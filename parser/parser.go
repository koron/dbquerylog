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
	body   *bytes.Buffer

	PktLens []int
	SeqNums []uint8
	Body    []byte
}

func New(r io.Reader) *Parser {
	return &Parser{
		r: bufio.NewReader(r),
	}
}

func (pa *Parser) initParse() {
	if pa.body == nil {
		pa.body = new(bytes.Buffer)
	}
	pa.body.Reset()
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
			break
		}
		_, err = io.CopyN(pa.body, pa.r, int64(pa.pktLen))
		if err != nil {
			return err
		}
		if pa.pktLen != maxPacketSize {
			break
		}
	}
	pa.Body = pa.body.Bytes()
	return nil
}

func (pa *Parser) String() string {
	var fb byte
	if len(pa.Body) > 0 {
		fb = pa.Body[0]
	}
	return fmt.Sprintf("PktLens=%+v SeqNums=%+v First=%02x",
		pa.PktLens, pa.SeqNums, fb)
}
