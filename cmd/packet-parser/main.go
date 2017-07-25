package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/koron/mysql-packet-sniffer/parser"
	"github.com/koron/mysql-packet-sniffer/tcpasm"
)

var nextStrm = 0

const serverPort = 3306

var independentParsers = map[string]*parser.Parser{}

// start MySQL packet parser in goroutine with s.
func created(src, dst tcpasm.Endpoint, s io.ReadCloser) error {
	if src.Port != serverPort && dst.Port != serverPort {
		return fmt.Errorf("both port not for MySQL: src=%s dst=%s", src, dst)
	}

	n := nextStrm
	nextStrm++
	log.Printf("strm#%d: src=%s dst=%s", n, src, dst)

	// Create a parser.
	var pa *parser.Parser
	var addr string
	if src.Port == serverPort {
		pa = parser.NewFromServer(s)
		addr = dst.String()
	} else {
		pa = parser.NewFromClient(s)
		addr = src.String()
	}

	// Share parser state.
	if p0, ok := independentParsers[addr]; ok {
		pa.ShareContext(p0)
		delete(independentParsers, addr)
	} else {
		independentParsers[addr] = pa
	}

	go func() {
		defer s.Close()
		for {
			err := pa.Parse()
			if err == io.EOF {
				log.Printf("strm#%d: EOF", n)
				if len(pa.Body) > 0 {
					log.Printf("  %#x", pa.Body)
				}
				return
			}
			if err != nil {
				log.Printf("strm#%d: parse failed: %s", n, err)
				return
			}
			// show last parsed MySQL packet.
			log.Printf("strm#%d: %s", n, pa.String())
		}
	}()

	return nil
}

func main() {
	flag.Parse()
	asm := &tcpasm.Assembler{
		Warn:    log.New(os.Stderr, "WARN ", log.LstdFlags),
		Created: created,
	}
	err := asm.Assemble(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
}
