package parser

type EOFPacket struct {
	Raw []byte
}

func NewEOFPacket(b []byte) (*EOFPacket, error) {
	// TODO:
	return &EOFPacket{Raw: b}, nil
}
