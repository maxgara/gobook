package main

import (
	"fmt"
	"math"
)

type Point struct{ X, Y float64 }

func main() {
	p := Point{3.0, 4.0}
	fmt.Println(p.Length())
}

func (p Point) Length() float64 {
	return math.Sqrt(p.X*p.X + p.Y*p.Y)
}
