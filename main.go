package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/koron/mysql-packet-sniffer/mysqlasm"
	"github.com/koron/mysql-packet-sniffer/parser"
	"github.com/koron/mysql-packet-sniffer/tcpasm"
)

type conn struct {
	out  *log.Logger
	addr tcpasm.Endpoint

	queryString string
	queryStart  time.Time
}

func newConn(addr tcpasm.Endpoint) mysqlasm.Conn {
	prefix := fmt.Sprintf("%s ", addr.String())
	return &conn{
		out:  log.New(os.Stdout, prefix, 0),
		addr: addr,
	}
}

func (c *conn) ID() string {
	return c.addr.String()
}

func (c *conn) Received(pa *parser.Parser, fromServer bool) {
	switch pkt := pa.Detail.(type) {
	case *parser.QueryPacket:
		c.queryString = pkt.Query
		c.queryStart = time.Now()
	case *parser.EOFPacket:
		if c.queryString != "" {
			c.finishQuery(pa)
			return
		}
		// TODO:
	case *parser.ResultNonePacket:
		if c.queryString != "" {
			c.finishQuery(pa)
			return
		}
		// TODO:
	}
}

func (c *conn) finishQuery(pa *parser.Parser) {
	d := time.Since(c.queryStart)
	c.out.Printf("query %q finished in %s", c.queryString, d)
	c.queryString = ""
}

func main() {
	flag.Parse()
	asm := mysqlasm.New(nil, newConn)
	asm.Warn = log.New(os.Stderr, "", log.LstdFlags)
	err := asm.Assemble(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
}
