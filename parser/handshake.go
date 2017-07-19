package parser

type ServerHandshakePacket struct {
}

func NewServerHandshakePacket(b []byte) (*ServerHandshakePacket, error) {
	// TODO:
	return &ServerHandshakePacket{}, nil
}

type ClientHandshakePacket struct {
}

func NewClientHandshakePacket(b []byte) (*ClientHandshakePacket, error) {
	// TODO:
	return &ClientHandshakePacket{}, nil
}
