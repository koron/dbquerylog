package parser

type EOFPacket struct {
	WarningCount uint16
	ServerStatus uint16
}

func NewEOFPacket(b []byte) (*EOFPacket, error) {
	var (
		pkt = &EOFPacket{}
		buf = &decbuf{buf: b[1:]}
	)
	pkt.WarningCount, _ = buf.ReadUint16()
	pkt.ServerStatus, _ = buf.ReadUint16()
	if buf.err != nil {
		return nil, buf.err
	}
	return pkt, nil
}
