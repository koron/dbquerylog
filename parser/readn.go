package parser

import "io"

func readN(r io.Reader, b []byte) error {
	for len(b) > 0 {
		n, err := r.Read(b)
		if err != nil {
			return err
		}
		b = b[n:]
	}
	return nil
}

func packetLen(b []byte) int {
	if len(b) < 3 {
		return 0
	}
	return int(uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16)
}
