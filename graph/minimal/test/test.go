package main

import (
	"fmt"
	"math"
)

func main() {
	x := math.MaxFloat64
	fmt.Printf("%v\n%v\n", x, x*2)
	n := math.IsNaN(math.NaN())
	fmt.Printf("%v\n", n)
}
