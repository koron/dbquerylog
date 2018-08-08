package parser

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
)

const maxPacketSize = 1<<24 - 1

// ReuseBufferMaxSize represent maximum size of buffer to reuse.  When size of
// buffer is bigger than this, buffer will be destory and create agein to free
// a big bunch of memory.
var ReuseBufferMaxSize = 8 * 1024 * 1024

type dir int

const (
	fromServer dir = iota
	fromClient
)

func (v dir) String() string {
	switch v {
	case fromServer:
		return "server"
	case fromClient:
		return "client"
	default:
		return "(dir:unknown)"
	}
}

type Parser struct {
	currR io.Reader
	raw   *CountReader

	dir dir
	ctx *Context

	decompressing bool

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
	cr := &CountReader{R: bufio.NewReader(r)}
	return &Parser{
		currR: cr,
		raw:   cr,
		dir:   fromServer,
		ctx:   newContext(),
	}
}

// NewFromServer creates a parser to parse packet from client.
func NewFromClient(r io.Reader) *Parser {
	cr := &CountReader{R: bufio.NewReader(r)}
	return &Parser{
		currR: cr,
		raw:   cr,
		dir:   fromClient,
		ctx:   newContext(),
	}
}

func (pa *Parser) resetBuffer() {
	if pa.body == nil {
		pa.body = new(bytes.Buffer)
		return
	}
	if ReuseBufferMaxSize > 0 && pa.body.Cap() > ReuseBufferMaxSize {
		pa.body = new(bytes.Buffer)
		return
	}
	pa.body.Reset()
}

func (pa *Parser) initParse() {
	pa.resetBuffer()
	if pa.PktLens == nil {
		pa.PktLens = make([]int, 0, 10)
	}
	pa.PktLens = pa.PktLens[:0]
	if pa.SeqNums == nil {
		pa.SeqNums = make([]uint8, 0, 10)
	}
	pa.SeqNums = pa.SeqNums[:0]
	pa.Body = nil

	if pa.shouldStartDecompress() {
		pa.switchDecompress(pa.raw)
	}

	// reset read bytes count.
	pa.raw.N = 0
}

func (pa *Parser) shouldStartDecompress() bool {
	return !pa.decompressing && pa.ctx.Compressing
}

func (pa *Parser) switchDecompress(r io.Reader) {
	pa.currR = newDecompressor(r)
	pa.decompressing = true
}

func (pa *Parser) toReader(b []byte) io.Reader {
	b2 := make([]byte, len(b))
	copy(b2, b)
	return bytes.NewBuffer(b)
}

func (pa *Parser) Parse() error {
	pa.initParse()
	for {
		err := readN(pa.currR, pa.header[:])
		if err != nil {
			return err
		}
		// re-parse stream with decompressing.
		if pa.shouldStartDecompress() {
			pa.switchDecompress(io.MultiReader(pa.toReader(pa.header[:]), pa.raw))
			continue
		}
		pa.pktLen = packetLen(pa.header[:])
		pa.PktLens = append(pa.PktLens, pa.pktLen)
		pa.SeqNums = append(pa.SeqNums, pa.header[3])
		if pa.pktLen == 0 {
			break
		}
		_, err = io.CopyN(pa.body, pa.currR, int64(pa.pktLen))
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

func (pa *Parser) PacketRawLen() uint64 {
	return pa.raw.N
}
