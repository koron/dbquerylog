package tcpasm

import (
	"io"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
)

func AssembleStream(s io.Reader, created StreamCreated) error {
	r, err := pcapgo.NewReader(s)
	if err != nil {
		return err
	}
	src := gopacket.NewPacketSource(r, layers.LayerTypeEthernet)
	asm := New(created)
	for {
		p, err := src.NextPacket()
		if err == io.EOF {
			return nil
		} else if err != nil {
			log.Print("WARN:", err)
			continue
		}
		Assemble(asm, p)
	}
}
