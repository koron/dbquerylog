package main

import (
	"flag"
	"io"
	"log"
	"os"

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

func streamCreated(netFlow, tcpFlow gopacket.Flow, s *tcpreader.ReaderStream) {
	// start MySQL packet parser in goroutine with s.
	pa := parser.New(s)
	n := nextStrm
	nextStrm++
	log.Printf("strm#%d: netFlow=%s tcpFlow=%s", n, netFlow, tcpFlow)
	go func() {
		defer s.Close()
		for {
			err := pa.Parse()
			if err == io.EOF {
				log.Printf("strm#%d: EOF", n)
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
