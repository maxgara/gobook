package main

import (
	"fmt"
	"os"
)

func main() {
	var s [1000]byte
	var s_idx int
	var sep byte = ' '

	for _, arg := range os.Args[1:] {
		for _, j := range arg {
			s[s_idx] = byte(j)
			s_idx++
		}
		s[s_idx] = sep
		s_idx++
	}
	fmt.Printf("%s\n", s)
}
