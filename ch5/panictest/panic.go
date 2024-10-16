package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {
	// defer fmt.Printf("end print stack\n\n\n\n")
	// defer printStack()
	// defer fmt.Println("defer print")
	// defer fmt.Println("defer print 2")
	fmt.Printf("%d\n", bad())
	fmt.Printf("done :)\n")
}

func bad() (r int) {
	cheat := func() {
		recover()
		r = 5
	}
	defer cheat()
	panic("bad")
}
func fix() {
	x := recover()
	fmt.Println(x)
	fmt.Println("fixed~")
}
func printStack() {
	var buf [4096]byte
	runtime.Stack(buf[:], false)
	fmt.Fprintf(os.Stdout, "%s\n", buf)
}
