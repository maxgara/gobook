package main

import (
	"fmt"
	"os"
	"testing"
)

func TestGen(t *testing.T) {
	f, err := os.OpenFile("temp.html", os.O_RDWR, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	g := newSVG()
	g.polyStart()
	g.vertex(1, 1)
	g.vertex(2, 2)
	g.vertex(3, 1)
	g.polyEnd()
	g.render(f)
}
