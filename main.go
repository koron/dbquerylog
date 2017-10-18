package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/koron/dbquerylog/mysqlasm"
	"github.com/koron/dbquerylog/parser"
	"github.com/koron/dbquerylog/tcpasm"
)

type conn struct {
	out io.Writer
	id  string

	report *Report

	preparing *statement
	prepared  map[uint32]*statement
}

type statement struct {
	id         uint32
	query      string
	fieldCount uint16
	paramCount uint16
}

var (
	warn = log.New(os.Stderr, "[WARN] ", 0)
	dbg  *log.Logger
)

func newConn(clientAddr, serverAddr tcpasm.Endpoint) mysqlasm.Conn {
	dbg.Println("")
	dbg.Printf("connected %s", clientAddr.String())
	return &conn{
		out: os.Stdout,
		id:  clientAddr.String(),
		report: &Report{
			ClientAddr: clientAddr,
			ServerAddr: serverAddr,
		},
		prepared: map[uint32]*statement{},
	}
}

func (c *conn) ID() string {
	return c.id
}

func (c *conn) Received(pa *parser.Parser, fromServer bool) {
	switch pkt := pa.Detail.(type) {

	case *parser.ClientHandshakePacket:
		c.report.Username = pkt.Username
		c.report.Database = pkt.Database

	case *parser.ServerHandshakePacket:
		// nothing to do.

	case *parser.QueryPacket:
		c.report.StartQuery(pkt.Query)

	case *parser.ExecuteQueryPacket:
		s, ok := c.getStatement(pkt.StatementID)
		if !ok {
			return
		}
		c.report.StartQuery(s.query, pkt.Parameters...)

	case *parser.ResultFieldNumPacket:
		if c.report.Querying() {
			c.report.ResponseSize += pa.PacketRawLen()
			c.report.ColumnNum = pkt.Num
		}

	case *parser.ResultFieldPacket:
		if c.report.Querying() {
			c.report.ResponseSize += pa.PacketRawLen()
		}

	case *parser.ResultRecordPacket:
		if c.report.Querying() {
			c.report.ResponseSize += pa.PacketRawLen()
			c.report.UpdatedRows++
		}

	case *parser.CloseQueryPacket:
		c.removeStatement(pkt.StatementID)

	case *parser.EOFPacket:
		if c.report.Querying() && pa.Context().ResultState == 0 {
			c.report.ResponseSize += pa.PacketRawLen()
			c.finishQuery()
			return
		}

	case *parser.ResultNonePacket:
		if c.report.Querying() {
			c.report.ResponseSize += pa.PacketRawLen()
			c.finishQuery()
			return
		}

	case *parser.PrepareQueryPacket:
		c.preparing = &statement{
			query: pkt.Query,
		}

	case *parser.PrepareResultPacket:
		if c.preparing == nil {
			return
		}
		c.preparing.id = pkt.StatementID
		c.preparing.fieldCount = pkt.FieldCount
		c.preparing.paramCount = pkt.ParameterCount
		c.addStatement(c.preparing)
		c.preparing = nil

	case *parser.OKPacket:
		// nothing to do yet.

	case *parser.ErrorPacket:
		if c.preparing != nil {
			c.preparing = nil
		}
		warn.Printf("ERROR: %s (%d)", pkt.Message, pkt.Number)

	case *parser.QuitPacket:
		dbg.Printf("QUIT: by command")

	default:
		if pkt == nil {
			dbg.Printf("IGNORED<nil>: first_byte=%02x", pa.Body[0])
			return
		}
		dbg.Printf("IGNORED: %#v", pkt)
	}
}

func (c *conn) Closed() {
	dbg.Printf("closed %s", c.id)
	dbg.Println("")
}

func (c *conn) finishQuery() {
	c.report.FinishQuery()
	defer c.report.Reset()
	if !includeSelect && strings.HasPrefix(strings.ToUpper(c.report.QueryString), "SELECT") {
		return
	}
	err := tsvWrite(c.out,
		c.report.StartTime.Format(time.RFC3339),
		strconv.FormatInt(c.report.StartTime.UnixNano(), 10),
		strconv.FormatInt(int64(c.report.ElapsedTime), 10),
		c.report.ClientAddr.Address,
		c.report.ClientAddr.PortString(),
		c.report.ServerAddr.Address,
		c.report.ServerAddr.PortString(),
		strconv.FormatUint(c.report.ResponseSize, 10),
		strconv.FormatUint(c.report.ColumnNum, 10),
		strconv.FormatUint(c.report.UpdatedRows, 10),
		c.report.Username,
		c.report.Database,
		c.report.QueryString,
		c.report.QueryParams,
	)
	if err != nil {
		warn.Printf("failed to output report: %s", err)
	}
}

func (c *conn) addStatement(s *statement) {
	t, ok := c.prepared[s.id]
	if ok {
		warn.Printf("duplicated statement %d: old=%q new=%q",
			s.id, t.query, s.query)
		return
	}
	c.prepared[s.id] = s
	dbg.Printf("PREPARE: %+v", s)
}

func (c *conn) getStatement(id uint32) (*statement, bool) {
	s, ok := c.prepared[id]
	if !ok {
		warn.Printf("statement not found: %d", id)
		return nil, false
	}
	return s, true
}

func (c *conn) removeStatement(id uint32) {
	s, ok := c.prepared[id]
	if !ok {
		warn.Printf("unknown statement %d", id)
		return
	}
	delete(c.prepared, id)
	dbg.Printf("DEALLOCATE: %+v", s)
}

var (
	debugFlag     bool
	includeSelect bool
	columnMaxlen  int
	decoder       string
)

func main() {
	flag.BoolVar(&debugFlag, "debug", false, "enable debug log")
	flag.BoolVar(&includeSelect, "select", false, "include SELECT statements")
	flag.IntVar(&columnMaxlen, "column_maxlen", 1024, "max length of columns")
	flag.StringVar(&decoder, "decoder", "Ethernet", "name of the decoder to use")
	flag.Parse()
	tsvValueMaxlen = columnMaxlen
	if debugFlag {
		dbg = log.New(os.Stderr, " [DBG] ", 0)
	} else {
		dbg = log.New(ioutil.Discard, "", 0)
	}
	dec, ok := gopacket.DecodersByLayerName[decoder]
	if !ok {
		log.Fatalf("no decoder: %s", decoder)
	}

	asm := mysqlasm.New(nil, newConn)
	asm.Warn = warn
	err := asm.Assemble(os.Stdin, dec)
	if err != nil {
		log.Fatal(err)
	}
}
