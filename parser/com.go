package parser

import (
	"fmt"
)

type COMPacket struct {
	Type CommandType
	Raw  []byte
}

func (p *COMPacket) CommandType() CommandType {
	return p.Type
}

func NewCOMPacket(b []byte) (interface{}, error) {
	if len(b) == 0 {
		return nil, fmt.Errorf("too short COM packet")
	}
	switch b[0] {

	case Quit:
		return &QuitPacket{}, nil

	case 0x03:
		pkt, err := NewQueryPacket(b)
		if err != nil {
			return nil, err
		}
		return pkt, err

	case 0x04:
		pkt, err := NewFieldListPacket(b)
		if err != nil {
			return nil, err
		}
		return pkt, err

	case 0x16:
		pkt, err := NewPrepareQueryPacket(b)
		if err != nil {
			return nil, err
		}
		return pkt, nil

	case 0x17:
		pkt, err := NewExecuteQueryPacket(b)
		if err != nil {
			return nil, err
		}
		return pkt, nil

	case Close:
		pkt, err := NewCloseQueryPacket(b)
		if err != nil {
			return nil, err
		}
		return pkt, nil

	default:
		// TODO: implement for other commands
		return &COMPacket{
			Type: CommandType(b[0]),
			Raw:  b[1:],
		}, nil
	}
}

type QuitPacket struct{}

func (pkt *QuitPacket) CommandType() CommandType {
	return Quit
}

type FieldListPacket struct {
	Table    string
	Wildcard string
}

func NewFieldListPacket(b []byte) (*FieldListPacket, error) {
	var (
		pkt = &FieldListPacket{}
		buf = &decbuf{buf: b[1:]}
	)
	pkt.Table, _ = buf.ReadString()
	pkt.Wildcard = string(buf.buf)
	if buf.err != nil {
		return nil, buf.err
	}
	return pkt, nil
}

func (pkt *FieldListPacket) CommandType() CommandType {
	return FieldList
}
