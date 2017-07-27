package mysqlasm

import (
	"sync"

	"github.com/koron/mysql-packet-sniffer/parser"
	"github.com/koron/mysql-packet-sniffer/tcpasm"
)

type Conn interface {
	ID() string
	Received(pa *parser.Parser, fromServer bool)
	Closed()
}

type ConnFactory func(client, server tcpasm.Endpoint) Conn

type conn struct {
	sync.Mutex
	ref int
	c   Conn
	p0  *parser.Parser
}

func (c *conn) received(pa *parser.Parser, fromServer bool) {
	c.Lock()
	c.c.Received(pa, fromServer)
	c.Unlock()
}

func (c *conn) closed() {
	c.c.Closed()
}

func (c *conn) incRef() {
	c.Lock()
	c.ref++
	c.Unlock()
}

func (c *conn) decRef() bool {
	c.Lock()
	defer c.Unlock()
	if c.ref > 0 {
		c.ref--
		if c.ref == 0 {
			c.closed()
			return true
		}
	}
	return false
}
