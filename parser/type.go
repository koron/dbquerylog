package parser

import "math"

type FieldType uint16

const (
	FieldTypeDecimal FieldType = iota
	FieldTypeTiny
	FieldTypeShort
	FieldTypeLong
	FieldTypeFloat
	FieldTypeDouble
	FieldTypeNULL
	FieldTypeTimestamp
	FieldTypeLongLong
	FieldTypeInt24
	FieldTypeDate
	FieldTypeTime
	FieldTypeDateTime
	FieldTypeYear
	FieldTypeNewDate
	FieldTypeVarChar
	FieldTypeBit
)
const (
	FieldTypeJSON FieldType = iota + 0xf5
	FieldTypeNewDecimal
	FieldTypeEnum
	FieldTypeSet
	FieldTypeTinyBLOB
	FieldTypeMediumBLOB
	FieldTypeLongBLOB
	FieldTypeBLOB
	FieldTypeVarString
	FieldTypeString
	FieldTypeGeometry
)

type fieldValue string

func (rv fieldValue) GoString() string {
	return string(rv)
}

const (
	fvNULL = fieldValue("NULL")
	fvNA   = fieldValue("N/A")
)

func (ft FieldType) readValue(b *decbuf) (interface{}, error) {
	switch ft {
	case FieldTypeNULL:
		return fvNULL, nil

	case FieldTypeTiny:
		v, err := b.ReadUint8()
		if err != nil {
			return nil, err
		}
		return v, nil

	case FieldTypeShort, FieldTypeYear:
		v, err := b.ReadUint16()
		if err != nil {
			return nil, err
		}
		return v, nil

	case FieldTypeLong, FieldTypeInt24:
		v, err := b.ReadUint32()
		if err != nil {
			return nil, err
		}
		return v, nil

	case FieldTypeLongLong:
		v, err := b.ReadUint64()
		if err != nil {
			return nil, err
		}
		return v, nil

	case FieldTypeFloat:
		v, err := b.ReadUint32()
		if err != nil {
			return nil, err
		}
		return math.Float32frombits(v), nil

	case FieldTypeDouble:
		v, err := b.ReadUint64()
		if err != nil {
			return nil, err
		}
		return math.Float64frombits(v), nil

	case FieldTypeDecimal, FieldTypeNewDecimal, FieldTypeVarChar,
		FieldTypeBit, FieldTypeEnum, FieldTypeSet, FieldTypeTinyBLOB,
		FieldTypeMediumBLOB, FieldTypeLongBLOB, FieldTypeBLOB,
		FieldTypeVarString, FieldTypeString, FieldTypeGeometry, FieldTypeJSON:
		v, err := b.ReadStringV()
		if err != nil {
			return nil, err
		}
		if v == nil {
			return fvNULL, nil
		}
		return *v, nil

	case FieldTypeDate, FieldTypeNewDate, FieldTypeTime, FieldTypeTimestamp,
		FieldTypeDateTime:
		v, err := b.ReadUintV()
		if err != nil {
			return nil, err
		}
		if v == nil {
			return fvNULL, nil
		}
		// TODO: parse time
		return *v, nil

	default:
		return fvNA, nil
	}
}
