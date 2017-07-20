package parser

import "fmt"

type ErrorPacket struct {
	Number  uint16
	State   string
	Message string
}

func NewErrorPacket(b []byte) (*ErrorPacket, error) {
	pkt := &ErrorPacket{}
	if b[0] != 0xff {
		return nil, fmt.Errorf("error packet must start with 0xff: %02x", b[0])
	}
	if len(b) < 9 {
		return nil, fmt.Errorf("too short data for error packet: %d", len(b))
	}
	pkt.Number = uint16(b[1]) | uint16(b[2])<<8
	if b[3] == 0x23 {
		pkt.State = string(b[4:9])
		pkt.Message = string(b[9:])
	} else {
		pkt.Message = string(b[3:])
	}
	return pkt, nil
}
