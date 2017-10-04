package tcpasm

import (
	"io"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"
)

type StreamCreated func(src, dst Endpoint, r io.ReadCloser) error

type Assembler struct {
	Warn    *log.Logger
	Decoder gopacket.Decoder
	Created StreamCreated
}

func (a *Assembler) warnf(f string, args ...interface{}) {
	if a.Warn == nil {
		return
	}
	a.Warn.Printf(f, args...)
}

func (a *Assembler) warn(args ...interface{}) {
	if a.Warn == nil {
		return
	}
	a.Warn.Print(args...)
}

func (a *Assembler) decoder() gopacket.Decoder {
	if a.Decoder == nil {
		return layers.LayerTypeEthernet
	}
	return a.Decoder
}

func (a *Assembler) New(netFlow, tcpFlow gopacket.Flow) tcpassembly.Stream {
	s := tcpreader.NewReaderStream()
	a.created(netFlow, tcpFlow, &s)
	return &s
}

func (a *Assembler) created(netFlow, tcpFlow gopacket.Flow, s *tcpreader.ReaderStream) {
	src, err := NewEndpoint(netFlow.Src(), tcpFlow.Src())
	if err != nil {
		s.Close()
		a.warnf("failed to build source: %s", err)
		return
	}
	dst, err := NewEndpoint(netFlow.Dst(), tcpFlow.Dst())
	if err != nil {
		s.Close()
		a.warnf("failed to build destination: %s", err)
		return
	}
	if a.Created != nil {
		err := a.Created(src, dst, s)
		if err != nil {
			s.Close()
			a.warnf("failed to create stream: %s", err)
		}
	}
}

func (a *Assembler) Assemble(r io.Reader) error {
	pr, err := pcapgo.NewReader(r)
	if err != nil {
		return err
	}
	src := gopacket.NewPacketSource(pr, a.decoder())
	asm := tcpassembly.NewAssembler(tcpassembly.NewStreamPool(a))
	for {
		p, err := src.NextPacket()
		if err == io.EOF {
			return nil
		} else if err != nil {
			a.warn(err)
			continue
		}
		assemble(asm, p)
	}
}

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
