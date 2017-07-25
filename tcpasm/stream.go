package tcpasm

import (
	"io"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
	"github.com/google/gopacket/tcpassembly/tcpreader"
)

type StreamCreated2 func(src, dst Endpoint, r io.ReadCloser) error

type StreamAssembler struct {
	WarnLog *log.Logger
	Created StreamCreated2
}

func (asm *StreamAssembler) warnf(f string, args ...interface{}) {
	if asm.WarnLog == nil {
		return
	}
	asm.WarnLog.Printf(f, args...)
}

func (asm *StreamAssembler) warn(args ...interface{}) {
	if asm.WarnLog == nil {
		return
	}
	asm.WarnLog.Print(args...)
}

func (asm *StreamAssembler) created(netFlow, tcpFlow gopacket.Flow, s *tcpreader.ReaderStream) {
	src, err := NewEndpoint(netFlow.Src(), tcpFlow.Src())
	if err != nil {
		asm.warnf("failed to build source: %s", err)
		return
	}
	dst, err := NewEndpoint(netFlow.Dst(), tcpFlow.Dst())
	if err != nil {
		asm.warnf("failed to build destination: %s", err)
		return
	}
	if asm.Created != nil {
		err :=	asm.Created(src, dst, s)
		if err != nil {
			asm.warnf("failed to create stream: %s", err)
		}
	}
}

func (asm *StreamAssembler) Assemble(r io.Reader) error {
	pr, err := pcapgo.NewReader(r)
	if err != nil {
		return err
	}
	src := gopacket.NewPacketSource(pr, layers.LayerTypeEthernet)
	factory := New(asm.created)
	for {
		p, err := src.NextPacket()
		if err == io.EOF {
			return nil
		} else if err != nil {
			asm.warn(err)
			continue
		}
		Assemble(factory, p)
	}
}

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
