package parser

import (
	"bytes"
	"compress/zlib"
	"io"
)

type decompressor struct {
	r io.Reader
	b *bytes.Buffer
	h [7]byte
}

func newDecompressor(r io.Reader) io.Reader {
	return &decompressor{
		r: r,
		b: new(bytes.Buffer),
	}
}

func (d *decompressor) Read(b []byte) (int, error) {
	if d.b.Len() == 0 {
		err := d.deflateNext()
		if err != nil {
			return 0, err
		}
	}
	return d.b.Read(b)
}

func (d *decompressor) deflateNext() error {
	d.b.Reset()
	err := readN(d.r, d.h[:])
	if err != nil {
		return err
	}
	var (
		clen = packetLen(d.h[0:3])
		_    = int(d.h[3])
		dlen = packetLen(d.h[4:7])
	)
	if clen == 0 {
		return nil
	}
	b := new(bytes.Buffer)
	_, err = io.CopyN(b, d.r, int64(clen))
	if err != nil {
		return err
	}
	// use raw bytes when decompressed length is less than compressed one.
	if dlen < clen {
		d.b.Grow(clen)
		_, err = io.CopyN(d.b, b, int64(clen))
		if err != nil {
			return err
		}
		return nil
	}
	d.b.Grow(dlen)
	fr, err := zlib.NewReader(b)
	if err != nil {
		return err
	}
	_, err = io.CopyN(d.b, fr, int64(dlen))
	if err != nil {
		return err
	}
	return nil
}
