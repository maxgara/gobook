package main

import (
	"fmt"
	"io"

	"maxgara-code.com/workspace/ch7/stringreader"
)

type dup stringreader.Sreader

func (d dup) Write() {
	return
}
func (d dup) Error() string {
	return "error?"
}
func (d dup) Read(p []byte) (n int, err error) {
	s := stringreader.Sreader(d)
	n, err = s.Read(p)
	return
}

func main() {
	sr := stringreader.Sreader{"test"}
	d := dup{"test"}
	test(&sr)
	test(&d)
}

func test(s io.Reader) {
	_, ok := s.(*stringreader.Sreader)
	if ok {
		fmt.Printf("%v is an io.Reader!\n", s)
	} else {
		fmt.Printf("%v is NOT an io.Reader! :(\n", s)
	}
	_, ok = s.(error)
	if ok {
		fmt.Printf("%v is an Error\n", s)
	} else {
		fmt.Printf("%v is not an error\n", s)
	}
}
