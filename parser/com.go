package parser

import (
	"fmt"
)

type COMPacket struct {
	Raw []byte
}

type QueryPacket struct {
	Query string
}

func NewCOMPacket(b []byte) (interface{}, error) {
	if len(b) == 0 {
		return nil, fmt.Errorf("too short COM packet")
	}
	switch b[0] {
	case 0x03:
		return &QueryPacket{Query: string(b[1:])}, nil
	default:
		// TODO: implement for other commands
		return &COMPacket{Raw: b}, nil
	}
}
