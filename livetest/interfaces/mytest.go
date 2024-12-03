package main

import (
	"fmt"
	"io"
	"os"
)

type mytype struct {
	s string
}

func (mt *mytype) Val() string {
	return mt.s
}

func (mt *mytype) GetReader() io.ReadCloser {
	return os.Stdin
}

type Valer interface {
	Val() string
	GetReader() io.Reader
}

func main() {
	var x = mytype{"hello"}
	var abs Valer
	abs = &x
	fmt.Println(abs)
}
