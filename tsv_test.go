package main

import "testing"

func TestTsvEscape(t *testing.T) {
	ok := func(s, exp string) {
		act := tsvEscape(s)
		if act != exp {
			t.Errorf("tsvEscape(%q) failed: expect=%q actual=%q", s, exp, act)
			return
		}
	}
	ok("foo", "foo")
	ok("ABC", "ABC")
	ok("foo\tbar", `foo\tbar`)
	ok("abc\nxyz", `abc\nxyz`)
	ok("ABC\rXYZ", `ABC\rXYZ`)
	ok("012\\789", `012\\789`)
}
