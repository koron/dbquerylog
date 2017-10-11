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

func TestLimit(t *testing.T) {
	ok := func(s string, max int, exp string) {
		tmp := tsvValueMaxlen
		tsvValueMaxlen = max
		defer func() { tsvValueMaxlen = tmp }()
		act := tsvLimit(s)
		if act != exp {
			t.Errorf("tsvLimit(%q) with %d failed: expect=%q actual=%q",
				s, max, exp, act)
			return
		}
	}
	ok("あいうえお", 8, "あい (...snipped)")
	ok("あいうえお", 9, "あいう (...snipped)")
	ok("あいうえお", 10, "あいう (...snipped)")
	ok("あいうえお", 11, "あいう (...snipped)")
	ok("あいうえお", 12, "あいうえ (...snipped)")
}
