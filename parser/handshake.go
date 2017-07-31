package parser

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
	Charset        *UintV
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
	pkt.Charset, _ = buf.ReadUintV()
	buf.Discard(23)
	pkt.Username, _ = buf.ReadString()
	pkt.HashedPassword, _ = buf.ReadStringV()
	pkt.Database, _ = buf.ReadString()
	if buf.err != nil {
		return nil, buf.err
	}
	return pkt, nil
}
