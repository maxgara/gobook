package main

import (
	"fmt"
	"io"
)

type ezreader struct {
	s *string
}

func NewReader(s string) *ezreader {
	var r = ezreader{&s}
	return &r
}
func (r ezreader) Read(p []byte) (int, error) {
	n := copy(p, *r.s)
	if n < len(*r.s) {
		*r.s = (*r.s)[n:]
		return n, nil
	}
	return n, io.EOF
}

func testReader() {
	var r = NewReader("test hello test")
	var s = make([]byte, 4)
	r.Read(s)
	fmt.Printf("%s\n", s)
	fmt.Printf("rval:%s\n", r.s)
	r.Read(s)
	fmt.Printf("%s\n", s)
}
