package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/google/gopacket"
	"github.com/google/gopacket/tcpassembly/tcpreader"
	"github.com/koron/mysql-packet-sniffer/parser"
	"github.com/koron/mysql-packet-sniffer/tcpasm"
)

var nextStrm = 0

const serverPort = 3306

func port(e gopacket.Endpoint) (uint16, error) {
	n, err := strconv.ParseUint(e.String(), 10, 16)
	if err != nil {
		return 0, err
	}
	return uint16(n), nil
}

var independentParsers = map[string]*parser.Parser{}

func streamCreated(netFlow, tcpFlow gopacket.Flow, s *tcpreader.ReaderStream) {
	// start MySQL packet parser in goroutine with s.
	n := nextStrm
	nextStrm++
	log.Printf("strm#%d: netFlow=%s tcpFlow=%s", n, netFlow, tcpFlow)

	// Check ports of source and destination.
	srcPort, err := port(tcpFlow.Src())
	if err != nil {
		log.Printf("strm#%d: failed to parse source port: %s", n, err)
		return
	}
	dstPort, err := port(tcpFlow.Dst())
	if err != nil {
		log.Printf("strm#%d: failed to parse destination port: %s", n, err)
		return
	}
	if srcPort != serverPort && dstPort != serverPort {
		log.Printf("strm#%d: both ports (%s) are not for MySQL", n, tcpFlow)
		return
	}

	// Create a parser.
	var pa *parser.Parser
	var addr string
	if srcPort == serverPort {
		pa = parser.NewFromServer(s)
		addr = fmt.Sprintf("%s:%s", netFlow.Dst(), tcpFlow.Dst())
	} else {
		pa = parser.NewFromClient(s)
		addr = fmt.Sprintf("%s:%s", netFlow.Src(), tcpFlow.Src())
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
}

func main() {
	flag.Parse()
	err := tcpasm.AssembleStream(os.Stdin, streamCreated)
	if err != nil {
		log.Fatal(err)
	}
}
