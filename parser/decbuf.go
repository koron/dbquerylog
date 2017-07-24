package parser

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
)

type decbuf struct {
	buf []byte
	err error
}

func (b *decbuf) ReadUint8() (uint8, error) {
	if b.err != nil {
		return 0, b.err
	}
	if len(b.buf) < 1 {
		b.err = io.EOF
		return 0, b.err
	}
	r := b.buf[0]
	b.buf = b.buf[1:]
	return r, nil
}

func (b *decbuf) ReadUint16() (uint16, error) {
	if b.err != nil {
		return 0, b.err
	}
	if len(b.buf) < 2 {
		b.err = io.EOF
		return 0, b.err
	}
	r := binary.LittleEndian.Uint16(b.buf)
	b.buf = b.buf[2:]
	return r, nil
}

func (b *decbuf) ReadUint32() (uint32, error) {
	if b.err != nil {
		return 0, b.err
	}
	if len(b.buf) < 4 {
		b.err = io.EOF
		return 0, b.err
	}
	r := binary.LittleEndian.Uint32(b.buf)
	b.buf = b.buf[4:]
	return r, nil
}

func (b *decbuf) ReadUint64() (uint64, error) {
	if b.err != nil {
		return 0, b.err
	}
	if len(b.buf) < 8 {
		b.err = io.EOF
		return 0, b.err
	}
	r := binary.LittleEndian.Uint64(b.buf)
	b.buf = b.buf[8:]
	return r, nil
}

func (b *decbuf) readNUint(n int) (uint64, error) {
	if len(b.buf) < n {
		b.err = io.EOF
		return 0, b.err
	}
	r := uint64(0)
	for i := 0; i < n; i++ {
		r += uint64(b.buf[i]) << (uint(i) * 8)
	}
	b.buf = b.buf[n:]
	return r, nil
}

func (b *decbuf) ReadUintV() (uint64, error) {
	if b.err != nil {
		return 0, b.err
	}
	if len(b.buf) < 1 {
		b.err = io.EOF
		return 0, b.err
	}
	f := b.buf[0]
	b.buf = b.buf[1:]
	if f < 0xFB {
		return uint64(f), nil
	}
	switch f {
	case 0xFC:
		return b.readNUint(2)
	case 0xFD:
		return b.readNUint(3)
	case 0xFE:
		return b.readNUint(8)
	default:
		b.err = fmt.Errorf("invalid byte for length-encoded integer: %02x", f)
		return 0, b.err
	}
}

func (b *decbuf) ReadString() (string, error) {
	if b.err != nil {
		return "", b.err
	}
	n := bytes.IndexByte(b.buf, 0)
	if n < 0 {
		b.err = errors.New("string didn't found \\0")
		return "", b.err
	}
	s := string(b.buf[:n])
	b.buf = b.buf[n+1:]
	return s, nil
}

func (b *decbuf) ReadStringV() (string, error) {
	if b.err != nil {
		return "", b.err
	}
	n, err := b.ReadUintV()
	if err != nil {
		return "", err
	}
	if n > math.MaxInt32 {
		b.err = fmt.Errorf("too long string: %d", n)
		return "", b.err
	}
	if len(b.buf) < int(n) {
		b.err = io.EOF
		return "", b.err
	}
	s := string(b.buf[:n])
	b.buf = b.buf[n:]
	return s, nil
}

func (b *decbuf) Discard(n int) error {
	if b.err != nil {
		return b.err
	}
	if len(b.buf) < n {
		b.err = io.EOF
		return b.err
	}
	b.buf = b.buf[n:]
	return nil
}
