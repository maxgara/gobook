package main

import (
	"fmt"
	"os"
	"testing"
)

func TestWebBase(t *testing.T) {
	f, err := os.OpenFile("temp.html", os.O_RDWR, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	w := newDoc()
	g := newGrid(2)
	g.add("hello")
	g.add("world")
	g.add("what's")
	g.add("up")
	w.render(f)
}
