package mysqlasm

import (
	"sync"

	"github.com/koron/mysql-packet-sniffer/parser"
	"github.com/koron/mysql-packet-sniffer/tcpasm"
)

type Conn interface {
	ID() string
	Received(pa *parser.Parser, fromServer bool)
}

type ConnFactory func(tcpasm.Endpoint) Conn

type conn struct {
	sync.Mutex
	c  Conn
	p0 *parser.Parser
}

func (c *conn) received(pa *parser.Parser, fromServer bool) {
	c.Lock()
	c.c.Received(pa, fromServer)
	c.Unlock()
}
