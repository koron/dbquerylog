package main

import (
	"io"
	"log"
	"os"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"
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

func streamCreated(netFlow, tcpFlow gopacket.Flow, s *tcpreader.ReaderStream) {
	// TODO: start MySQL packet parser in goroutine with s.
	log.Printf("netFlow=%s tcpFlow=%s")
}

func main() {
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
