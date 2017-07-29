package parser

import (
	"fmt"
)

type OKPacket struct {
	AffectedRows *UintV
	InsertID     *UintV
	Status       uint16
	WarningCount uint16
	Message      string
}

func NewOKPacket(b []byte) (*OKPacket, error) {
	var (
		pkt = &OKPacket{}
		buf = &decbuf{buf: b}
	)
	if b[0] != 0x00 {
		return nil, fmt.Errorf("OK packet must start with 0x00: %02x", b[0])
	}
	pkt.AffectedRows, _ = buf.ReadUintV()
	pkt.InsertID, _ = buf.ReadUintV()
	pkt.Status, _ = buf.ReadUint16()
	pkt.WarningCount, _ = buf.ReadUint16()
	// FIXME: read string.
	if buf.err != nil {
		return nil, buf.err
	}
	return pkt, nil
}
