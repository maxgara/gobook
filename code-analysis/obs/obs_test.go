package main

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	s := "abc abc abd abx bay pav"
	x := ParseNd{name: "text", s: s, p: make(map[string]ParseG)}
	words := x.Parse("words", "ab")
	fmt.Print(words)
}
