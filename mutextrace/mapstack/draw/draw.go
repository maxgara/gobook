package main

import (
	"fmt"
	"time"
)

func main() {
	// b := [1]byte{}
	// os.(b)
	fmt.Printf("testtest\n")
	time.Sleep(time.Second * 3)
	fmt.Printf("\b\b\033[A\b\b")
}
