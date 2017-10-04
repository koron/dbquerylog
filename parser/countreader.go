package parser

import "io"

type CountReader struct {
	R io.Reader
	N uint64
}

func (cr *CountReader) Read(p []byte) (int, error) {
	n, err := cr.R.Read(p)
	if n > 0 {
		cr.N += uint64(n)
	}
	return n, err
}
