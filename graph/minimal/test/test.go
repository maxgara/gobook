package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Printf("%f, %f\n", math.MaxFloat64, -math.MaxFloat64)
	fmt.Printf("%v, %v\n ", 5.0 < math.MaxFloat64, 5.0 > -math.MaxFloat64)
}
