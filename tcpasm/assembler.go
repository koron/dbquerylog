package tcpasm

import "github.com/google/gopacket/tcpassembly"

func New(created StreamCreated) *tcpassembly.Assembler {
	f := &StreamFactory{
		created: created,
	}
	p := tcpassembly.NewStreamPool(f)
	return tcpassembly.NewAssembler(p)
}
