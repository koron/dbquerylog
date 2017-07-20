package parser

import (
	"bufio"
	"bytes"
	"encoding/binary"
)

type ServerHandshakePacket struct {
	ProtocolVersion  uint8
	ServerVersion    string
	ThreadID         uint32
	ScrambleBuffer   uint64
	Filter           uint8
	ServerCapability uint16
	ServerLanguage   uint8
	ServerStatus     uint16
}

func NewServerHandshakePacket(b []byte) (*ServerHandshakePacket, error) {
	var (
		err error
		pkt = &ServerHandshakePacket{}
		r   = bufio.NewReader(bytes.NewReader(b))
	)
	pkt.ProtocolVersion, err = r.ReadByte()
	if err != nil {
		return nil, err
	}
	pkt.ServerVersion, err = r.ReadString(0x00)
	if err != nil {
		return nil, err
	}
	err = binary.Read(r, binary.LittleEndian, &pkt.ThreadID)
	if err != nil {
		return nil, err
	}
	err = binary.Read(r, binary.LittleEndian, &pkt.ScrambleBuffer)
	if err != nil {
		return nil, err
	}
	pkt.Filter, err = r.ReadByte()
	if err != nil {
		return nil, err
	}
	err = binary.Read(r, binary.LittleEndian, &pkt.ServerCapability)
	if err != nil {
		return nil, err
	}
	pkt.ServerLanguage, err = r.ReadByte()
	if err != nil {
		return nil, err
	}
	err = binary.Read(r, binary.LittleEndian, &pkt.ServerStatus)
	if err != nil {
		return nil, err
	}
	// TODO:
	return pkt, nil
}

type ClientHandshakePacket struct {
}

func NewClientHandshakePacket(b []byte) (*ClientHandshakePacket, error) {
	// TODO:
	return &ClientHandshakePacket{}, nil
}
