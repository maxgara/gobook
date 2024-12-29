package main

import (
	"bufio"
	"bytes"
	"testing"
)

func TestParse(t *testing.T) {
	b := bytes.NewBuffer([]byte("hello this is a test of the parsing tool. I need to see if word separation is done correctly\n\n\n\nok?"))
	p := parser{s: bufio.NewScanner(b)}
	p.parse()
}
