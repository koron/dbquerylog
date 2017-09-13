package parser

import (
	"bytes"
	"compress/flate"
	"io"
	"log"
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
	log.Printf("deflateNext: HERE1")
	d.b.Reset()
	err := readN(d.r, d.h[:])
	log.Printf("deflateNext: HERE2")
	if err != nil {
		log.Printf("deflateNext: ERROR: %s", err)
		return err
	}
	var (
		clen = packetLen(d.h[0:3])
		pnum = int(d.h[3])
		dlen = packetLen(d.h[4:7])
	)
	log.Printf("deflateNext: clen=%d pnum=%d dlen=%d", clen, pnum, dlen)
	if clen == 0 {
		return nil
	}
	b := new(bytes.Buffer)
	_, err = io.CopyN(b, d.r, int64(clen))
	if err != nil {
		return err
	}
	d.b.Grow(dlen)
	fr := flate.NewReader(b)
	_, err = io.CopyN(d.b, fr, int64(dlen))
	if err != nil {
		return err
	}
	return nil
}
