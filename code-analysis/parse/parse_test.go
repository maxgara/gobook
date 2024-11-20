package main

import "fmt"

//	func ExampleParseNd_Parse() {
//		s := "abc abc abd abx bay pav"
//		x := NewParseNd(s)
//		words := x.Parse("ab")
//		fmt.Print(words)
//	}
func ExampleParseNd_Parse() {
	s := "abc abd abx bay pab"
	x := NewParseNd(s)
	words := x.Parse("ab.?")
	fmt.Print(words)
	//output: [node:abc
	// 	p:
	// 	anc:0x0
	//  node:abd
	// 	p:
	// 	anc:0x0
	//  node:abx
	// 	p:
	// 	anc:0x0
	//  node:ab
	// 	p:
	// 	anc:0x0
	// ]
}
func ExampleParseNd_NamedParse() {
	s := "abc abd abx bay pab"
	x := NewParseNd(s)
	words := x.NamedParse("abc (?<myname>..)")
	fmt.Print(words)
	//output: [node:abc
	// 	p:
	// 	anc:0x0
	//  node:abd
	// 	p:
	// 	anc:0x0
	//  node:abx
	// 	p:
	// 	anc:0x0
	//  node:ab
	// 	p:
	// 	anc:0x0
	// ]
}
