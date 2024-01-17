package parser

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
)

type decbuf struct {
	buf []byte
	err error

	errorChecked bool
}

var ErrNotEnoughBuffer = errors.New("not enough buffer")

func (b *decbuf) Read(p []byte) (int, error) {
	n, m := len(p), len(b.buf)
	if n > m {
		n = m
	}
	copy(p[:n], b.buf[:n])
	b.buf = b.buf[n:]
	if n < len(p) {
		b.err = io.EOF
		return n, b.err
	}
	return n, nil
}

func (b *decbuf) ReadUint8() (uint8, error) {
	if b.err != nil {
		return 0, b.err
	}
	if len(b.buf) < 1 {
		b.err = ErrNotEnoughBuffer
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
		b.err = ErrNotEnoughBuffer
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
		b.err = ErrNotEnoughBuffer
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
		b.err = ErrNotEnoughBuffer
		return 0, b.err
	}
	r := binary.LittleEndian.Uint64(b.buf)
	b.buf = b.buf[8:]
	return r, nil
}

func (b *decbuf) readNUint(n int) (uint64, error) {
	if len(b.buf) < n {
		b.err = ErrNotEnoughBuffer
		return 0, b.err
	}
	r := uint64(0)
	for i := 0; i < n; i++ {
		r += uint64(b.buf[i]) << (uint(i) * 8)
	}
	b.buf = b.buf[n:]
	return r, nil
}

// ReadUintV parses an int<lenenc> value.
func (b *decbuf) ReadUintV() (*UintV, error) {
	if b.err != nil {
		return nil, b.err
	}
	if len(b.buf) < 1 {
		b.err = ErrNotEnoughBuffer
		return nil, b.err
	}
	f := b.buf[0]
	b.buf = b.buf[1:]
	if f < 0xFB {
		n := UintV(f)
		return &n, nil
	}
	switch f {
	case 0xFB:
		return nil, nil
	case 0xFC:
		m, err := b.readNUint(2)
		n := UintV(m)
		return &n, err
	case 0xFD:
		m, err := b.readNUint(3)
		n := UintV(m)
		return &n, err
	case 0xFE:
		m, err := b.readNUint(8)
		n := UintV(m)
		return &n, err
	default:
		b.err = fmt.Errorf("invalid byte for length-encoded integer: %02x (len:%d)", f, len(b.buf))
		return nil, b.err
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

func (b *decbuf) ReadStringV() (*StringV, error) {
	if b.err != nil {
		return nil, b.err
	}
	p, err := b.ReadUintV()
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, nil
	}
	n := *p
	if n > math.MaxInt32 {
		b.err = fmt.Errorf("too long string: %d", n)
		return nil, b.err
	}
	if len(b.buf) < int(n) {
		b.err = ErrNotEnoughBuffer
		return nil, b.err
	}
	s := StringV(b.buf[:n])
	b.buf = b.buf[n:]
	return &s, nil
}

// ReadStringAll reads remained whole buffer as a string.
func (b *decbuf) ReadStringAll() (string, error) {
	if b.err != nil {
		return "", b.err
	}
	s := string(b.buf)
	b.buf = b.buf[len(b.buf):]
	return s, nil
}

func (b *decbuf) Discard(n int) error {
	if b.err != nil {
		return b.err
	}
	if len(b.buf) < n {
		b.err = ErrNotEnoughBuffer
		return b.err
	}
	b.buf = b.buf[n:]
	return nil
}

func (b *decbuf) DiscardV() (*UintV, error) {
	if b.err != nil {
		return nil, b.err
	}
	p, err := b.ReadUintV()
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, nil
	}
	n := *p
	if n > math.MaxInt32 {
		b.err = fmt.Errorf("too long blob: %d", n)
		return nil, b.err
	}
	return p, b.Discard(int(n))
}

func (b *decbuf) ErrorCheck(label string) {
	if b.errorChecked || b.err == nil {
		return
	}
	b.errorChecked = true
	log.Printf("error occurred on %s: %s", label, b.err)
}
