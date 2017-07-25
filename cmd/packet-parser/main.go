package main

import (
	"flag"
	"log"
	"os"

	"github.com/koron/mysql-packet-sniffer/mysqlasm"
	"github.com/koron/mysql-packet-sniffer/parser"
	"github.com/koron/mysql-packet-sniffer/tcpasm"
)

type conn struct {
	addr tcpasm.Endpoint
}

func (c *conn) ID() string {
	return c.addr.String()
}

func (c *conn) Received(pa *parser.Parser, fromServer bool) {
	dir := "client"
	if fromServer {
		dir = "server"
	}
	log.Printf("%s(%s): %s", c.ID(), dir, pa.String())
}

func newConn(addr tcpasm.Endpoint) mysqlasm.Conn {
	log.Printf("connected with %s", addr)
	return &conn{
		addr: addr,
	}
}

func main() {
	flag.Parse()
	asm := mysqlasm.New(nil, newConn)
	asm.Warn = log.New(os.Stderr, "WARN ", log.LstdFlags)
	err := asm.Assemble(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
}
