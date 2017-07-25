package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"
	"github.com/koron/mysql-packet-sniffer/parser"
	"github.com/koron/mysql-packet-sniffer/tcpasm"
)

func assemble(asm *tcpassembly.Assembler, p gopacket.Packet) {
	if p.NetworkLayer() == nil || p.TransportLayer() == nil || p.TransportLayer().LayerType() != layers.LayerTypeTCP {
		return
	}
	var (
		flow = p.NetworkLayer().NetworkFlow()
		tcp  = p.TransportLayer().(*layers.TCP)
		time = p.Metadata().Timestamp
	)
	asm.AssembleWithTimestamp(flow, tcp, time)
}

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
	}
	dstPort, err := port(tcpFlow.Dst())
	if err != nil {
		log.Printf("strm#%d: failed to parse destination port: %s", n, err)
	}
	if srcPort != serverPort && dstPort != serverPort {
		log.Printf("strm#%d: both ports (%s) are not for MySQL", n, tcpFlow)
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
	r, err := pcapgo.NewReader(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	src := gopacket.NewPacketSource(r, layers.LayerTypeEthernet)
	asm := tcpasm.New(streamCreated)
	for {
		p, err := src.NextPacket()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Print("WARN:", err)
			continue
		}
		assemble(asm, p)
	}
}
