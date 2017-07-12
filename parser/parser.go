package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

const maxPacketSize = 1<<24 - 1

type Parser struct {
	r *bufio.Reader

	header [4]byte
	pktLen int

	PktLens []int
	SeqNums []uint8
	Body    *bytes.Buffer
}

func New(r io.Reader) *Parser {
	return &Parser{
		r: bufio.NewReader(r),
	}
}

func (pa *Parser) initParse() {
	if pa.Body == nil {
		pa.Body = new(bytes.Buffer)
	}
	pa.Body.Reset()
	if pa.PktLens == nil {
		pa.PktLens = make([]int, 0, 10)
	}
	pa.PktLens = pa.PktLens[:0]
	if pa.SeqNums == nil {
		pa.SeqNums = make([]uint8, 0, 10)
	}
	pa.SeqNums = pa.SeqNums[:0]
}

func (pa *Parser) Parse() error {
	pa.initParse()
	for {
		err := readN(pa.r, pa.header[:])
		if err != nil {
			return err
		}
		pa.pktLen = packetLen(pa.header[:])
		pa.PktLens = append(pa.PktLens, pa.pktLen)
		pa.SeqNums = append(pa.SeqNums, pa.header[3])
		if pa.pktLen == 0 {
			return nil
		}
		_, err = io.CopyN(pa.Body, pa.r, int64(pa.pktLen))
		if err != nil {
			return err
		}
		if pa.pktLen != maxPacketSize {
			return nil
		}
	}
}

func (pa *Parser) String() string {
	b := pa.Body.Bytes()
	var fb byte
	if len(b) > 0 {
		fb = b[0]
	}
	return fmt.Sprintf("PktLens=%+v SeqNums=%+v First=%02x",
		pa.PktLens, pa.SeqNums, fb)
}
