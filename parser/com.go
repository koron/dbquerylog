package parser

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
)

type COMPacket struct {
	FirstByte uint8
}

func NewCOMPacket(b []byte) (*COMPacket, error) {
	var (
		pkt = &COMPacket{}
	)
	pkt.FirstByte = b[0]
	// TODO:
	return &COMPacket{}, nil
}

type OKPacket struct {
	AffectedRows uint64
	InsertID     uint64
	Status       uint16
	WarningCount uint16
	Message      string
}

func NewOKPacket(b []byte) (*OKPacket, error) {
	var (
		err error
		pkt = &OKPacket{}
		r   = bufio.NewReader(bytes.NewReader(b[1:]))
	)
	if b[0] != 0x00 {
		return nil, fmt.Errorf("OK packet must start with 0x00: %02x", b[0])
	}
	pkt.AffectedRows, err = readLengthEncodedInteger(r)
	if err != nil {
		return nil, err
	}
	pkt.InsertID, err = readLengthEncodedInteger(r)
	err = binary.Read(r, binary.LittleEndian, &pkt.Status)
	if err != nil {
		return nil, err
	}
	err = binary.Read(r, binary.LittleEndian, &pkt.WarningCount)
	if err != nil {
		return nil, err
	}
	// FIXME: read string.
	return pkt, nil
}
