package main

import (
	"bytes"
	"io"
	"strings"
	"sync"
)

func tsvWrite(w io.Writer, values ...string) error {
	for i, v := range values {
		if i != 0 {
			_, err := io.WriteString(w, "\t")
			if err != nil {
				return err
			}
		}
		_, err := io.WriteString(w, tsvEscape(v))
		if err != nil {
			return err
		}
	}
	_, err := io.WriteString(w, "\n")
	if err != nil {
		return err
	}
	return err
}

var tsvPool = &sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func tsvEscape(s string) string {
	if strings.IndexAny(s, "\t\n\r\\") == -1 {
		return s
	}
	b := tsvPool.Get().(*bytes.Buffer)
	for _, r := range s {
		switch r {
		case '\t':
			b.WriteString(`\t`)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '\\':
			b.WriteString(`\\`)
		default:
			b.WriteRune(r)
		}
	}
	t := b.String()
	b.Reset()
	tsvPool.Put(b)
	return t
}
