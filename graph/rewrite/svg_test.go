package main

import (
	"math/rand/v2"
	"os"
	"testing"
)

func TestSvgGen(t *testing.T) {
	b := newDocBuilder(os.Stdout)
	b.StartSVGBlock("test", "testx", "testy")
	for range 20 {
		b.startPoly("poly")
		x0 := rand.Float64()
		x1 := rand.Float64()
		y0 := rand.Float64()
		y1 := rand.Float64()
		b.vertex(x0, y0)
		b.vertex(x1, y1)
		b.endPoly()
	}
	b.EndSVGBlock()
	b.SendDoc()
}
func TestWebNode(t *testing.T) {
	// f, err := os.OpenFile("temp.html", os.O_WRONLY, os.ModePerm)
	// if err != nil {
	// 	return
	// }
	b := newDocBuilder(os.Stdout)
	b.StartContentNode(&Div{s: "hello world"}, Node{name: "div"})
	b.SendDoc()
}

func TestDivFunctions(t *testing.T) {
	b := newDocBuilder(os.Stdout)
	b.StartDiv()
	b.EndDiv()
	b.SendDoc()
}
