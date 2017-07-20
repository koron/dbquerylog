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
	Charset          uint8
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
	pkt.Charset, err = r.ReadByte()
	if err != nil {
		return nil, err
	}
	err = binary.Read(r, binary.LittleEndian, &pkt.ServerStatus)
	if err != nil {
		return nil, err
	}
	// FIXME: parse other fields.
	return pkt, nil
}

type ClientHandshakePacket struct {
	ClientFlags    uint32
	MaxPacketSize  uint32
	Charset        uint64
	Username       string
	HashedPassword string
	Database       string
}

func NewClientHandshakePacket(b []byte) (*ClientHandshakePacket, error) {
	var (
		err error
		pkt = &ClientHandshakePacket{}
		r   = bufio.NewReader(bytes.NewReader(b))
	)
	err = binary.Read(r, binary.LittleEndian, &pkt.ClientFlags)
	if err != nil {
		return nil, err
	}
	err = binary.Read(r, binary.LittleEndian, &pkt.MaxPacketSize)
	if err != nil {
		return nil, err
	}
	pkt.Charset, err = readLengthEncodedInteger(r)
	if err != nil {
		return nil, err
	}
	_, err = r.Discard(23)
	if err != nil {
		return nil, err
	}
	pkt.Username, err = r.ReadString(0x00)
	if err != nil {
		return nil, err
	}
	pkt.HashedPassword, err = readLengthEncodedString(r)
	if err != nil {
		return nil, err
	}
	pkt.Database, err = r.ReadString(0x00)
	if err != nil {
		return nil, err
	}
	return pkt, nil
}
