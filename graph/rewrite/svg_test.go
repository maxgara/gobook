package main

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
)

func TestGen(t *testing.T) {
	//this is art lol
	f, err := os.OpenFile("temp.html", os.O_RDWR|os.O_TRUNC, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	d := newDocBuilder(f)
	d.startDoc()
	d.writeTitle("Test Doc Title")
	d.startGrid(3)
	for range 20 {
		d.startGridElem()
		d.startSVG("Test SVG Title", [4]float64{0, 0, 1, 1})
		for range 20 {
			d.startPoly(0.05)
			x0 := rand.Float64()
			x1 := rand.Float64()
			y0 := rand.Float64()
			y1 := rand.Float64()
			d.vertex(x0, y0)
			d.vertex(x1, y1)
			d.endPoly()
		}
		d.endSVG()
		d.endGridElem()
	}
	d.endGrid()
	d.endDoc()
}
