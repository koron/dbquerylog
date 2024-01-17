package parser

import "fmt"

type QueryPacket struct {
	ParamCount    *UintV
	ParamSetCount *UintV
	Query         string
}

// NewQueryPacket parses a COM_QUERY packet.
// See https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_com_query.html also
func NewQueryPacket(b []byte, ctx *Context) (*QueryPacket, error) {
	// skip first byte, caller must check it.
	var (
		pkt QueryPacket
		buf = &decbuf{buf: b[1:]}
	)
	// Parse query attributes (enabled by ClientFlag.ClientQueryAttributes)
	if ctx.QueryAttributes {
		pkt.ParamCount, _ = buf.ReadUintV()
		pkt.ParamSetCount, _ = buf.ReadUintV() // currently always 1
		if num := pkt.ParamCount.Uint64(); num != 0 {
			// TODO: read parameters. see ExecuteQueryPacket.readParams also
		}
	}
	pkt.Query, _ = buf.ReadStringAll()
	if buf.err != nil {
		return nil, buf.err
	}
	return &pkt, nil
}

func (p *QueryPacket) CommandType() CommandType {
	return Query
}

type PrepareQueryPacket struct {
	Query string
}

func NewPrepareQueryPacket(b []byte) (*PrepareQueryPacket, error) {
	// skip first byte, caller must check it.
	return &PrepareQueryPacket{Query: string(b[1:])}, nil
}

func (p *PrepareQueryPacket) CommandType() CommandType {
	return Prepare
}

type PrepareResultPacket struct {
	StatementID    uint32
	FieldCount     uint16
	ParameterCount uint16
	WarningCount   uint16
}

func NewPrepareResultPacket(b []byte) (*PrepareResultPacket, error) {
	// skip first byte, caller must check it.
	var (
		pkt = &PrepareResultPacket{}
		buf = &decbuf{buf: b[1:]}
	)
	pkt.StatementID, _ = buf.ReadUint32()
	pkt.FieldCount, _ = buf.ReadUint16()
	pkt.ParameterCount, _ = buf.ReadUint16()
	n1, _ := buf.ReadUint8()
	pkt.WarningCount, _ = buf.ReadUint16()
	if buf.err != nil {
		return nil, buf.err
	}
	if n1 != 0 {
		return nil, fmt.Errorf(
			"unexpected first spacer for prepare_result packet: %d", n1)
	}
	return pkt, nil
}

type ExecuteQueryPacket struct {
	StatementID uint32
	CursorType  uint8
	Types       []FieldType
	Parameters  []interface{}
}

func NewExecuteQueryPacket(b []byte, ctx *Context) (*ExecuteQueryPacket, error) {
	// skip first byte, caller must check it.
	var (
		pkt = &ExecuteQueryPacket{}
		buf = &decbuf{buf: b[1:]}
	)
	pkt.StatementID, _ = buf.ReadUint32()
	pkt.CursorType, _ = buf.ReadUint8()
	n1, _ := buf.ReadUint32()
	err := pkt.readParams(buf, ctx)
	if err != nil {
		return nil, err
	}
	if buf.err != nil {
		return nil, buf.err
	}
	if n1 != 1 {
		return nil, fmt.Errorf(
			"unexpected first spacer for exec_query packet: %d", n1)
	}
	return pkt, nil
}

func (p *ExecuteQueryPacket) readParams(buf *decbuf, ctx *Context) error {
	if buf.err != nil {
		return buf.err
	}
	st, ok := ctx.PreparedStmts[p.StatementID]
	if !ok {
		return fmt.Errorf("statement not found: %d", p.StatementID)
	}
	// read and parse bitmap.
	var bm *bitmap
	if st.NumParams > 0 {
		b := make([]byte, (st.NumParams+7)/8)
		_, err := buf.Read(b)
		if err != nil {
			return err
		}
		bm = &bitmap{b: b, m: st.NumParams}
	}
	sp, _ := buf.ReadUint8()
	if sp != 1 {
		return fmt.Errorf(
			"unexpected 2nd parser for exec_query packet: %d", sp)
	}
	if st.NumParams == 0 {
		return nil
	}
	types := make([]FieldType, st.NumParams)
	// read and parse parameter types
	for i := uint16(0); i < st.NumParams; i++ {
		t, err := buf.ReadUint16()
		if err != nil {
			return err
		}
		types[i] = FieldType(t)
	}
	p.Types = types
	// read and parse parameter values
	values := make([]interface{}, 0, st.NumParams)
	for i := uint16(0); i < st.NumParams; i++ {
		if bm.get(i) {
			values = append(values, fvNULL)
			continue
		}
		var t FieldType
		t, types = types[0], types[1:]
		v, err := t.readValue(buf)
		if err != nil {
			values = append(values, fieldValue(fmt.Sprintf("<ERR:%s>", err)))
			break
		}
		values = append(values, v)
	}
	p.Parameters = values
	return nil
}

func (p *ExecuteQueryPacket) CommandType() CommandType {
	return Execute
}

type CloseQueryPacket struct {
	StatementID uint32
}

func NewCloseQueryPacket(b []byte) (*CloseQueryPacket, error) {
	var (
		pkt = &CloseQueryPacket{}
		buf = &decbuf{buf: b[1:]}
	)
	pkt.StatementID, _ = buf.ReadUint32()
	if buf.err != nil {
		return nil, buf.err
	}
	return pkt, nil
}
