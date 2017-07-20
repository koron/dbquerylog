package parser

import (
	"fmt"
	"io"
)

type reader interface {
	io.Reader
	io.ByteReader
}

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

func readInt(r io.ByteReader, nbytes int) (uint64, error) {
	n := uint64(0)
	for i := 0; i < nbytes; i++ {
		b, err := r.ReadByte()
		if err != nil {
			return 0, err
		}
		n += uint64(b) << (uint(i) * 8)
	}
	return n, nil
}

func readLengthEncodedInteger(r reader) (n uint64, err error) {
	b, err := r.ReadByte()
	if err != nil {
		return 0, err
	}
	if b <= 0xFB {
		return uint64(b), nil
	}
	switch b {
	case 0xFC:
		return readInt(r, 2)
	case 0xFD:
		return readInt(r, 3)
	case 0xFE:
		return readInt(r, 8)
	default:
		return 0, fmt.Errorf(
			"invalid first byte for length-encoded integer: %02x", b)
	}
}

func readLengthEncodedString(r reader) (string, error) {
	n, err := readLengthEncodedInteger(r)
	if err != nil {
		return "", err
	}
	b := make([]byte, n)
	err = readN(r, b)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
