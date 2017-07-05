package tcpasm

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"
)

type StreamCreated func(netFlow, tcpFlow gopacket.Flow, r *tcpreader.ReaderStream)

type StreamFactory struct {
	created StreamCreated
}

func (f *StreamFactory) New(netFlow, tcpFlow gopacket.Flow) tcpassembly.Stream {
	s := tcpreader.NewReaderStream()
	if f.created != nil {
		f.created(netFlow, tcpFlow, &s)
	}
	return &s
}
