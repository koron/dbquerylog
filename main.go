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
}

var warn = log.New(os.Stderr, "", log.LstdFlags)

func newConn(clientAddr, serverAddr tcpasm.Endpoint) mysqlasm.Conn {
	return &conn{
		out: os.Stdout,
		id:  clientAddr.String(),
		report: &Report{
			ClientAddr: clientAddr,
			ServerAddr: serverAddr,
		},
	}
}

func (c *conn) ID() string {
	return c.id
}

func (c *conn) Received(pa *parser.Parser, fromServer bool) {
	switch pkt := pa.Detail.(type) {

	case *parser.ClientHandshakePacket:
		c.report.Username = pkt.Username

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

	default:
		warn.Printf("IGNORED: %#v", pkt)
	}
}

func (c *conn) Closed() {
	warn.Printf("closed %s", c.id)
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
		warn.Print(err)
	}
	c.report.Reset()
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
