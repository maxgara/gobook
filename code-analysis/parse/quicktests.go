package main

import "fmt"

func ExampleParseNd_Parse() {
	s := "abc abd abx bay pab"
	x := NewParseNd(s)
	words := x.Parse("ab.?")
	fmt.Print(words)
}
