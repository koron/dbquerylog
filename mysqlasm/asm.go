package mysqlasm

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/koron/mysql-packet-sniffer/parser"
	"github.com/koron/mysql-packet-sniffer/tcpasm"
)

type Assembler struct {
	ctx   context.Context
	f     ConnFactory
	l     sync.Mutex
	conns map[string]*conn

	ServerPort uint16
	Warn       *log.Logger
}

func New(ctx context.Context, f ConnFactory) *Assembler {
	if ctx == nil {
		ctx = context.Background()
	}
	return &Assembler{
		ctx:        ctx,
		f:          f,
		ServerPort: 3306,
	}
}

func (a *Assembler) warnf(f string, args ...interface{}) {
	if a.Warn == nil {
		return
	}
	a.Warn.Printf(f, args...)
}

func (a *Assembler) created(src, dst tcpasm.Endpoint, r io.ReadCloser) error {
	if src.Port != a.ServerPort && dst.Port != a.ServerPort {
		return fmt.Errorf("both port not for MySQL: src=%s dst=%s", src, dst)
	}

	var (
		addr tcpasm.Endpoint
		pa   *parser.Parser
		from bool
	)
	if src.Port == a.ServerPort {
		addr = dst
		pa = parser.NewFromServer(r)
		from = true
	} else {
		addr = src
		pa = parser.NewFromClient(r)
		from = false
	}
	c := a.getConn(addr, pa, from)
	go a.parseLoop(r, pa, from, c)

	return nil
}

func (a *Assembler) getConn(addr tcpasm.Endpoint, pa *parser.Parser, fromServer bool) *conn {
	a.l.Lock()
	defer a.l.Unlock()
	if a.conns == nil {
		a.conns = make(map[string]*conn)
	}
	if c, ok := a.conns[addr.String()]; ok {
		if c.p0 != nil {
			pa.ShareContext(c.p0)
			c.p0 = nil
		}
		return c
	}
	c := &conn{
		c:  a.f(addr),
		p0: pa,
	}
	a.conns[addr.String()] = c
	return c
}

func (a *Assembler) parseLoop(r io.ReadCloser, pa *parser.Parser, fromServer bool, c *conn) {
	defer r.Close()
	for {
		select {
		case <-a.ctx.Done():
			return
		default:
		}
		err := pa.Parse()
		if err == io.EOF {
			a.warnf("stream closed for %s", c.c.ID())
			return
		}
		if err != nil {
			a.warnf("failed parser for %s: %s", c.c.ID(), err)
		}
		c.received(pa, fromServer)
	}
	// TODO: close conn
}

func (a *Assembler) Assemble(r io.Reader) error {
	asm := tcpasm.Assembler{
		Warn:    a.Warn,
		Created: a.created,
	}
	return asm.Assemble(r)
}