package main

import "fmt"

func main() {
	var x int
	var y int
	var z int
	z = x / 2
	z = z * 2
	y = z - 3
	fmt.Println(y)
}

func f1(x int) int {
	return 3 * x
}
