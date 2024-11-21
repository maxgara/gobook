package main

import (
	"fmt"
	"testing"
)

//	func ExampleParseNd_Parse() {
//		s := "abc abc abd abx bay pav"
//		x := NewParseNd(s)
//		words := x.Parse("ab")
//		fmt.Print(words)
//	}
func ExampleParseNd_Parse() {
	s := "abc abd abx bay pab"
	x := NewParseNd(s)
	words := x.Parse("abc (?<myname>..)")
	fmt.Println(words)
	fmt.Println(x)
	// Output: asdasasa
}
func TestPname(t *testing.T) {
	p := "abc (?<myname>..)"
	pn := pname(p)
	fmt.Println(pn)
}
