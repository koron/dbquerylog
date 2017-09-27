package parser

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
)

const maxPacketSize = 1<<24 - 1

type dir int

const (
	fromServer dir = iota
	fromClient
)

type Parser struct {
	r   io.Reader
	dir dir
	ctx *Context

	compressing bool

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
		ctx: newContext(),
	}
}

// NewFromServer creates a parser to parse packet from client.
func NewFromClient(r io.Reader) *Parser {
	return &Parser{
		r:   bufio.NewReader(r),
		dir: fromClient,
		ctx: newContext(),
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
	pa.Body = nil

	if !pa.compressing && pa.ctx.Compressing {
		pa.switchDecompress()
	}
}

func (pa *Parser) switchDecompress() {
	log.Printf("switch to decompressing stream\n")
	pa.r = newDecompressor(pa.r)
	pa.compressing = true
}

func (pa *Parser) Parse() error {
	pa.initParse()
	for {
		err := readN(pa.r, pa.header[:])
		if err != nil {
			return err
		}
		// re-parse stream with decompressing.
		if !pa.compressing && pa.ctx.Compressing {
			// TODO: push back.
			switchDecompress()
			continue
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

func (pa *Parser) deflatePacket() error {
	h := make([]byte, 3)
	err := readN(pa.r, h)
	if err != nil {
		return err
	}
	deflateLen := packetLen(h)
	if deflateLen == 0 {
		return nil
	}
	fmt.Printf("pktLen=%d deflateLen=%d\n", pa.pktLen, deflateLen)
	bb := new(bytes.Buffer)
	_, err = io.CopyN(bb, pa.r, int64(pa.pktLen-len(h)))
	if err != nil {
		return err
	}
	/*
		r := flate.NewReader(bb)
		_, err = io.CopyN(pa.body, r, int64(deflateLen))
		if err != nil {
			return err
		}
	*/
	return nil
}

func (pa *Parser) parseServerPacket() error {
	if len(pa.Body) < 1 {
		return errors.New("less body as packet from server")
	}
	if pa.ctx.State == None {
		pkt, err := NewServerHandshakePacket(pa.Body)
		if err != nil {
			return err
		}
		pa.ctx.State = Handshake
		pa.Detail = pkt
		return nil
	}
	switch pa.Body[0] {
	case 0x00:
		if pa.ctx.State == Auth {
			// logged in successfully.
			pkt, err := NewOKPacket(pa.Body)
			if err != nil {
				return err
			}
			pa.Detail = pkt
			pa.ctx.State = Connected
			pa.ctx.Compressing = pa.ctx.WillCompress
			break
		}
		switch pa.ctx.LastCommand {

		case Prepare:
			pkt, err := NewPrepareResultPacket(pa.Body)
			if err != nil {
				return err
			}
			pa.Detail = pkt
			if pkt.ParameterCount > 0 {
				if pkt.FieldCount > 0 {
					pa.ctx.ResultState = PrepareParamsAndColumns
				} else {
					pa.ctx.ResultState = PrepareParams
				}
			} else {
				if pkt.FieldCount > 0 {
					pa.ctx.ResultState = PrepareColumns
				} else {
					pa.ctx.ResultState = 0
				}
			}
			pa.ctx.addStmt(Stmt{
				ID:         pkt.StatementID,
				NumParams:  pkt.ParameterCount,
				NumColumns: pkt.FieldCount,
			})
			return nil

		case Query, Execute, Reset:
			err := pa.parseServerResultPacket()
			if err != nil {
				return err
			}
			return nil

		default:
			// TODO: parse other server commands
			pa.Detail = nil
			return nil
		}

	case 0xfe:
		pkt, err := NewEOFPacket(pa.Body)
		if err != nil {
			return err
		}
		pa.Detail = pkt
		switch pa.ctx.ResultState {
		case PrepareParamsAndColumns:
			pa.ctx.ResultState = PrepareColumns
			return nil
		case PrepareParams, PrepareColumns:
			pa.ctx.ResultState = 0
			return nil
		}
		if !pa.ctx.IsClientDeprecateEOF() && pa.ctx.ResultState == Fields {
			pa.ctx.ResultState = Records
			return nil
		}
		pa.ctx.ResultState = 0

	case 0xff:
		pkt, err := NewErrorPacket(pa.Body)
		if err != nil {
			return err
		}
		pa.Detail = pkt
	default:
		// FIXME: any specific procedure?
		err := pa.parseServerResultPacket()
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (pa *Parser) parseServerResultPacket() error {
	switch pa.ctx.ResultState {
	case 0:
		buf := &decbuf{buf: pa.Body}
		p, err := buf.ReadUintV()
		if err != nil {
			return err
		}
		if p == nil {
			return nil
		}
		n := *p
		if n == 0 {
			pkt, err := NewResultNonePacket(buf.buf)
			if err != nil {
				return err
			}
			pa.Detail = pkt
			pa.ctx.ResultState = 0
			return nil
		}
		pa.Detail = &ResultFieldNumPacket{Num: uint64(n)}
		pa.ctx.ResultState = Fields
		pa.ctx.FieldNCurr = 0
		pa.ctx.FieldNMax = uint64(n)
		return nil

	case Fields, PrepareParamsAndColumns, PrepareParams, PrepareColumns:
		pkt, err := NewResultFieldPacket(pa.Body)
		if err != nil {
			return err
		}
		pa.Detail = pkt
		pa.ctx.FieldNCurr++
		if pa.ctx.IsClientDeprecateEOF() && pa.ctx.FieldNCurr >= pa.ctx.FieldNMax {
			pa.ctx.ResultState = Records
		}
		return nil

	case Records:
		var nfields int
		if pa.ctx.FieldNMax <= math.MaxInt32 {
			nfields = int(pa.ctx.FieldNMax)
		} else {
			nfields = math.MaxInt32
		}
		pkt, err := NewResultRecordPacket(pa.Body, nfields)
		if err != nil {
			return err
		}
		pa.Detail = pkt
		return nil

	default:
		return fmt.Errorf("unexpected query result mode: %d",
			pa.ctx.ResultState)
	}
}

func (pa *Parser) parseClientPacket() error {
	switch pa.ctx.State {
	case Handshake:
		pkt, err := NewClientHandshakePacket(pa.Body)
		if err != nil {
			return err
		}
		pa.ctx.ClientFlags = pkt.ClientFlags
		pa.ctx.WillCompress = pa.ctx.ClientFlags&ClientCompress != 0
		pa.ctx.State = Auth
		pa.Detail = pkt
	case AuthResend:
		pkt, err := NewClientAuthResendPacket(pa.Body)
		if err != nil {
			return err
		}
		pa.ctx.State = Auth
		pa.Detail = pkt
	case AwaitingReply:
		pkt, err := NewAwaitingReplyPacket(pa.Body)
		if err != nil {
			return err
		}
		pa.Detail = pkt
	default:
		pkt, err := NewCOMPacket(pa.Body, pa.ctx)
		if err != nil {
			return err
		}
		pa.Detail = pkt
		if c, ok := pkt.(Command); ok {
			pa.ctx.LastCommand = c.CommandType()
		} else {
			pa.ctx.LastCommand = -1
		}
	}
	return nil
}

func (pa *Parser) ShareContext(other *Parser) {
	if pa == other {
		return
	}
	pa.ctx = other.ctx
}

func (pa *Parser) Context() Context {
	return *pa.ctx
}

func (pa *Parser) String() string {
	if pa.Detail == nil {
		var fb byte
		if len(pa.Body) > 0 {
			fb = pa.Body[0]
		}
		return fmt.Sprintf("PktLens=%+v SeqNums=%+v First=%02x",
			pa.PktLens, pa.SeqNums, fb)
	}
	return fmt.Sprintf("Detail=%#v PktLens=%+v SeqNums=%+v",
		pa.Detail, pa.PktLens, pa.SeqNums)
}

func (pa *Parser) ContextData() interface{} {
	return pa.ctx.Data
}

func (pa *Parser) SetContextData(d interface{}) {
	pa.ctx.Data = d
}
