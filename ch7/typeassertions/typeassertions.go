package main

import (
	"fmt"
	"io"
	"os"

	"maxgara-code.com/workspace/ch7/stringreader"
)

type dup stringreader.Sreader

func (d dup) Write() {
  fmt.Printf("this is a string %v\n", d)
  
  os.ReadFile("this is a file name")
	return
}
func (d dup) Error() string {
