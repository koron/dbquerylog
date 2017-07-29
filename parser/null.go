package parser

import (
	"fmt"
	"strconv"
)

type StringV string

func (v *StringV) String() string {
	if v == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%q", *v)
}

func (v *StringV) GoString() string {
	return v.String()
}

type UintV uint64

func (v *UintV) String() string {
	if v == nil {
		return "<nil>"
	}
	return "0x" + strconv.FormatUint(uint64(*v), 16)
}

func (v *UintV) GoString() string {
	return v.String()
}
