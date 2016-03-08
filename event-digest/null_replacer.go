package main

import (
	"io"
	"unicode/utf8"
)

// NullReplacer wraps an io.Reader and replaces instances of \x00 with a rune
// of your choice!
type NullReplacer struct {
	Replacement rune
	Stream      io.Reader
}

// NewNullReplacer returns a valid NullReplacer, ensuring that the replacement
// given is only one byte when encoded in UTF-8
func NewNullReplacer(stream io.Reader, replacement rune) NullReplacer {
	if utf8.RuneLen(replacement) != 1 {
		panic("multibyte rune cannot replace single byte")
	}
	return NullReplacer{
		Replacement: replacement,
		Stream:      stream,
	}
}

func (r NullReplacer) Read(p []byte) (n int, err error) {
	n, err = r.Stream.Read(p)
	for i := 0; i < n; i++ {
		if p[i] == 0x00 {
			p[i] = byte(r.Replacement)
		}
	}
	return n, err
}
