package draw

import (
	"fmt"
)

// func main() {
// 	// b := [1]byte{}
// 	// os.(b)
// 	fmt.Printf("testtest\n")
// 	time.Sleep(time.Second * 3)
// 	fmt.Printf("\b\b\033[A\b\b")
// }

type Line []byte
type Canvas struct {
	lines []Line
}

// place line of length l at line y, pos x, and oriented in direction d. d has allowed values ('u','d','l','r').
// if pointy is set to true, the line will terminate with an "arrowhead" ie. ---->
func putLine(x int, y int, d byte, l int, can *Canvas, pointy bool) {
	var c byte         //character to write
	var p byte         //pointy end of line (when point is set true)
	var xoff, yoff int //direction to move writier per char
	switch d {
	case 'u':
		c = '|'
		yoff = -1
		p = '^'
	case 'd':
		c = '|'
		yoff = 1
		p = '?'
	case 'l':
		c = '-'
		xoff = -1
		p = '<'
	case 'r':
		c = '-'
		xoff = 1
		p = '>'
	}
	var effectivel = l
	if pointy {
		effectivel++
	}
	if x+(effectivel*xoff) < 0 || y+(effectivel*yoff) < 0 || x+(effectivel*xoff) > len(can.lines[0]) || y+(effectivel*yoff) > len(can.lines) {
		can.PrintBytes()
		panic(fmt.Sprintf("address below zero. Tried to draw line from (%d,%d) dir %c; l %d in canvas:\n", x, y, d, l))
	}
	for i := 0; i < l; i++ {
		x += xoff
		y += yoff
		can.lines[y][x] = c
	}
	if !pointy {
		return
	}
	x += xoff
	y += yoff
	can.lines[y][x] = p
}

// write byte slice to canvas, starting at line y, positiion x
func putString(x int, y int, s []byte, can *Canvas) {
	for _, c := range s {
		can.lines[y][x] = c
	}
}

// print canvas to console
func (can Canvas) Print() {
	for _, l := range can.lines {
		for _, c := range l {
			fmt.Printf("%c", c)
		}
		fmt.Println()
	}
}

// convenience function to see all byte values in canvas
func (can Canvas) PrintBytes() {
	for _, l := range can.lines {
		for _, c := range l {
			fmt.Printf("%d", c)
		}
		fmt.Println()
	}
}

// initialize blank canvas of width w and height h
func makeCanvas(w int, h int) *Canvas {
	blank := make([]byte, w)
	for i := range blank {
		blank[i] = ' '
	}
	can := Canvas{make([]Line, h)}
	for i := range can.lines {
		can.lines[i] = make(Line, w)
		copy(can.lines[i], blank) //note: must use reference (i), not value. Also, must use copy not := or every line points to the same slice :(.)
	}
	return &can
}
