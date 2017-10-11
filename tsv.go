package main

import (
	"io"
	"strconv"
	"unicode/utf8"
)

func tsvWrite(w io.Writer, values ...string) error {
	for i, v := range values {
		if i != 0 {
			_, err := io.WriteString(w, "\t")
			if err != nil {
				return err
			}
		}
		_, err := io.WriteString(w, tsvLimit(tsvEscape(v)))
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

var tsvValueMaxlen int = 0

func tsvLimit(s string) string {
	if tsvValueMaxlen <= 0 || len(s) < tsvValueMaxlen {
		return s
	}
	n := tsvValueMaxlen
	for i, r := range s {
		if i+utf8.RuneLen(r) > tsvValueMaxlen {
			n = i
			break
		}
	}
	return s[0:n] + " (...snipped)"
}

func tsvEscape(s string) string {
	s = strconv.Quote(s)
	return s[1 : len(s)-1]
}
