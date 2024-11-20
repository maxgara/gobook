package main

import "fmt"

func ExampleParseNd_Parse_temp() {
	s := "abc abd abx bay pab"
	x := NewParseNd(s)
	words := x.Parse("ab.?")
	fmt.Print(words)
	//Output: [node:abc
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
