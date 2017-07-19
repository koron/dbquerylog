package parser

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
)

const maxPacketSize = 1<<24 - 1

type dir int

const (
	fromServer dir = iota
	fromClient
)

type state int

const (
	none state = iota
	handshake
	auth
	connected
)

type Parser struct {
	r   *bufio.Reader
	dir dir
	st  state

	header [4]byte
	pktLen int
	body   *bytes.Buffer

	PktLens []int
	SeqNums []uint8
	Body    []byte

	Detail interface{}
}

// NewFromServer creates a parser to parse packet from server.
func NewFromServer(r io.Reader) *Parser {
	return &Parser{
		r:   bufio.NewReader(r),
		dir: fromServer,
	}
}

// NewFromServer creates a parser to parse packet from client.
func NewFromClient(r io.Reader) *Parser {
	return &Parser{
		r:   bufio.NewReader(r),
		dir: fromClient,
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
	pa.Detail = nil
	switch pa.dir {
	case fromServer:
		return pa.parseServerPacket()
	case fromClient:
		return pa.parseClientPacket()
	default:
		return fmt.Errorf("unknown direction: %s", pa.dir)
	}
}

func (pa *Parser) parseServerPacket() error {
	if len(pa.Body) < 1 {
		return errors.New("less body as packet from server")
	}
	if pa.st == none {
		pkt, err := NewServerHandshakePacket(pa.Body)
		if err != nil {
			return err
		}
		pa.st = handshake
		pa.Detail = pkt
		return nil
	}
	switch pa.Body[0] {
	case 0x00:
		if pa.st == auth {
			// TODO: logged in successfully.
			pa.st = connected
			break
		}
		return pa.parseServerResultPacket()
	case 0xfe:
		pkt, err := NewEOFPacket(pa.Body)
		if err != nil {
			return err
		}
		pa.Detail = pkt
	case 0xff:
		pkt, err := NewErrorPacket(pa.Body)
		if err != nil {
			return err
		}
		pa.Detail = pkt
	default:
		// FIXME: any specific procedure?
		return pa.parseServerResultPacket()
	}
	return nil
}

func (pa *Parser) parseServerResultPacket() error {
	// TODO: parse processing results if any commands are running.
	pa.Detail = &ServerResultPacket{}
	return nil
}

func (pa *Parser) parseClientPacket() error {
	// TODO:
	return nil
}

func (pa *Parser) String() string {
	if pa.Detail == nil {
		var fb byte
		if len(pa.Body) > 0 {
			fb = pa.Body[0]
		}
		return fmt.Sprintf("[%d] PktLens=%+v SeqNums=%+v First=%02x",
			pa.dir, pa.PktLens, pa.SeqNums, fb)
	}
	return fmt.Sprintf("[%d] PktLens=%+v SeqNums=%+v Detail=%+v",
		pa.dir, pa.PktLens, pa.SeqNums, pa.Detail)
}
