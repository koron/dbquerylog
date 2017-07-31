package main

import (
	"flag"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/koron/mysql-packet-sniffer/mysqlasm"
	"github.com/koron/mysql-packet-sniffer/parser"
	"github.com/koron/mysql-packet-sniffer/tcpasm"
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
	dbg  = log.New(os.Stderr, " [DBG] ", 0)
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

	case *parser.ServerHandshakePacket:
		// nothing to do.

	case *parser.QueryPacket:
		c.report.StartQuery(pkt.Query)

	case *parser.ResultFieldNumPacket:
		if c.report.Querying() {
			c.report.ResponseSize += uint64(len(pa.Body))
			c.report.ColumnNum = pkt.Num
		}

	case *parser.ResultFieldPacket:
		if c.report.Querying() {
			c.report.ResponseSize += uint64(len(pa.Body))
		}

	case *parser.ResultRecordPacket:
		if c.report.Querying() {
			c.report.ResponseSize += uint64(len(pa.Body))
			c.report.UpdatedRows++
		}

	case *parser.EOFPacket:
		if c.report.Querying() {
			c.report.ResponseSize += uint64(len(pa.Body))
			c.finishQuery()
			return
		}

	case *parser.ResultNonePacket:
		if c.report.Querying() {
			c.report.ResponseSize += uint64(len(pa.Body))
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
	err := tsvWrite(c.out,
		c.report.StartTime.String(),
		c.report.ClientAddr.String(),
		c.report.ServerAddr.String(),
		c.report.Username,
		strconv.FormatUint(c.report.ResponseSize, 10),
		strconv.FormatUint(c.report.ColumnNum, 10),
		strconv.FormatUint(c.report.UpdatedRows, 10),
		strconv.FormatInt(int64(c.report.ElapsedTime), 10),
		c.report.QueryString,
		c.report.QueryParams,
	)
	if err != nil {
		warn.Printf("failed to output report: %s", err)
	}
	c.report.Reset()
}

func (c *conn) addStatement(s *statement) {
	t, ok := c.prepared[s.id]
	if ok {
		warn.Printf("duplicated statement %d: old=%q new=%q",
			s.id, t.query, s.query)
		return
	}
	c.prepared[s.id] = s
	dbg.Printf("prepared: %+v", s)
}

func main() {
	flag.Parse()
	asm := mysqlasm.New(nil, newConn)
	asm.Warn = warn
	err := asm.Assemble(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
}
