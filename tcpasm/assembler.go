package tcpasm

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/tcpassembly"
)

func New(created StreamCreated) *tcpassembly.Assembler {
	f := &StreamFactory{
		created: created,
	}
	p := tcpassembly.NewStreamPool(f)
	return tcpassembly.NewAssembler(p)
}

func Assemble(asm *tcpassembly.Assembler, p gopacket.Packet) {
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
