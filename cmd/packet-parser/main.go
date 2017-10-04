package main

import (
	"flag"
	"log"
	"os"

	"github.com/koron/dbquerylog/mysqlasm"
	"github.com/koron/dbquerylog/parser"
	"github.com/koron/dbquerylog/tcpasm"
)

var l = log.New(os.Stdout, "", log.LstdFlags)

type conn struct {
	addr tcpasm.Endpoint
}

func newConn(clientAddr, serverAddr tcpasm.Endpoint) mysqlasm.Conn {
	l.Printf("connected with %s", clientAddr)
	return &conn{
		addr: clientAddr,
	}
}

func (c *conn) ID() string {
	return c.addr.String()
}

func (c *conn) Received(pa *parser.Parser, fromServer bool) {
	dir := "client"
	if fromServer {
		dir = "server"
	}
	l.Printf("%s(%s): %s", c.ID(), dir, pa.String())
}

func (c *conn) Closed() {
	l.Printf("closed %s", c.ID())
}

func main() {
	flag.Parse()
	asm := mysqlasm.New(nil, newConn)
	asm.Warn = log.New(os.Stderr, "WARN ", log.LstdFlags)
	err := asm.Assemble(os.Stdin, nil)
	if err != nil {
		log.Fatal(err)
	}
}
