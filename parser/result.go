package parser

import (
	"fmt"
)

type ResultPacket struct {
}

type ResultNonePacket struct {
	ResultPacket
	AffectedRows *UintV
	InsertID     *UintV
	ServerStatus uint16
	WarningCount uint16
	Message      *StringV
}

func NewResultNonePacket(b []byte) (*ResultNonePacket, error) {
	var (
		pkt = &ResultNonePacket{}
		buf = &decbuf{buf: b}
	)
	pkt.AffectedRows, _ = buf.ReadUintV()
	pkt.InsertID, _ = buf.ReadUintV()
	pkt.ServerStatus, _ = buf.ReadUint16()
	pkt.WarningCount, _ = buf.ReadUint16()
	if len(buf.buf) > 0 {
		pkt.Message, _ = buf.ReadStringV()
	}
	if buf.err != nil {
		return nil, buf.err
	}
	return pkt, nil
}

type ResultFieldNumPacket struct {
	Num uint64
}

type ResultFieldPacket struct {
	ResultPacket
	Database     *StringV
	Table        *StringV
	TableOrigin  *StringV
	Column       *StringV
	ColumnOrigin *StringV
	Charset      uint16
	Length       uint32
	Type         uint8
	Flag         uint16
	DotN         uint8
	Default      *StringV
}

func NewResultFieldPacket(b []byte) (*ResultFieldPacket, error) {
	var (
		pkt = &ResultFieldPacket{}
		buf = &decbuf{buf: b}
	)
	s, _ := buf.ReadStringV()
	if s == nil || *s != "def" {
		return nil, fmt.Errorf(
			"unexpected header for result field packet: %+v", s)
	}
	pkt.Database, _ = buf.ReadStringV()
	pkt.Table, _ = buf.ReadStringV()
	pkt.TableOrigin, _ = buf.ReadStringV()
	pkt.Column, _ = buf.ReadStringV()
	pkt.ColumnOrigin, _ = buf.ReadStringV()
	n1, _ := buf.ReadUint8()
	pkt.Charset, _ = buf.ReadUint16()
	pkt.Length, _ = buf.ReadUint32()
	pkt.Type, _ = buf.ReadUint8()
	pkt.Flag, _ = buf.ReadUint16()
	pkt.DotN, _ = buf.ReadUint8()
	n2, _ := buf.ReadUint16()
	if len(buf.buf) > 0 {
		pkt.Default, _ = buf.ReadStringV()
	}
	if buf.err != nil {
		return nil, buf.err
	}
	if n1 != 12 {
		return nil, fmt.Errorf(
			"unexpected first spacer for result field packet: %d", n1)
	}
	if n2 != 0 {
		return nil, fmt.Errorf(
			"unexpected second spacer for result field packet: %d", n2)
	}
	return pkt, nil
}

type ResultRecordPacket struct {
	ResultPacket
	Columns []*StringV
}

func NewResultRecordPacket(b []byte, nfields int) (*ResultRecordPacket, error) {
	var (
		pkt = &ResultRecordPacket{
			Columns: make([]*StringV, nfields),
		}
		buf = &decbuf{buf: b}
	)
	var err error
	for i := 0; i < nfields; i++ {
		pkt.Columns[i], err = buf.ReadStringV()
		if err != nil {
			return nil, err
		}
	}
	return pkt, nil
}
