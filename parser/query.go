package parser

import "fmt"

type QueryPacket struct {
	Query string
}

func NewQueryPacket(b []byte) (*QueryPacket, error) {
	// skip first byte, caller must check it.
	return &QueryPacket{Query: string(b[1:])}, nil
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
}

func NewExecuteQueryPacket(b []byte) (*ExecuteQueryPacket, error) {
	// skip first byte, caller must check it.
	var (
		pkt = &ExecuteQueryPacket{}
		buf = &decbuf{buf: b[1:]}
	)
	pkt.StatementID, _ = buf.ReadUint32()
	pkt.CursorType, _ = buf.ReadUint8()
	n1, _ := buf.ReadUint32()
	if buf.err != nil {
		return nil, buf.err
	}
	if n1 != 1 {
		return nil, fmt.Errorf(
			"unexpected first spacer for exec_query packet: %d", n1)
	}
	// TODO: parse other fields
	return pkt, nil
}

func (p *ExecuteQueryPacket) CommandType() CommandType {
	return Execute
}
