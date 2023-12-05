package parser

import (
	"fmt"
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
		pkt = &ServerHandshakePacket{}
		buf = &decbuf{buf: b}
	)
	pkt.ProtocolVersion, _ = buf.ReadUint8()
	pkt.ServerVersion, _ = buf.ReadString()
	pkt.ThreadID, _ = buf.ReadUint32()
	pkt.ScrambleBuffer, _ = buf.ReadUint64()
	pkt.Filter, _ = buf.ReadUint8()
	pkt.ServerCapability, _ = buf.ReadUint16()
	pkt.Charset, _ = buf.ReadUint8()
	pkt.ServerStatus, _ = buf.ReadUint16()
	if buf.err != nil {
		return nil, buf.err
	}
	// FIXME: parse other fields.
	return pkt, nil
}

type ClientHandshakePacket struct {
	ClientFlags    ClientFlags
	MaxPacketSize  uint32
	Charset        uint8
	Username       string
	HashedPassword *StringV
	Database       string
}

func NewClientHandshakePacket(b []byte) (*ClientHandshakePacket, error) {
	var (
		pkt = &ClientHandshakePacket{}
		buf = &decbuf{buf: b}
	)
	cflags, _ := buf.ReadUint32()
	pkt.ClientFlags = ClientFlags(cflags)
	pkt.MaxPacketSize, _ = buf.ReadUint32()
	pkt.Charset, _ = buf.ReadUint8()
	buf.Discard(23)
	// XXX: To support SSL, stop parsing here.
	// See also https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_connection_phase_packets_protocol_ssl_request.html
	pkt.Username, _ = buf.ReadString()
	pkt.HashedPassword, _ = buf.ReadStringV()
	if pkt.ClientFlags&ClientConnectWithDB != 0 {
		pkt.Database, _ = buf.ReadString()
	}
	if buf.err != nil {
		return nil, fmt.Errorf("failed on parsing ClientHandshakePacket: %w", buf.err)
	}
	return pkt, nil
}
