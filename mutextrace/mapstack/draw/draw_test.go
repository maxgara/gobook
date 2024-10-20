package draw

import (
	"fmt"
	"testing"
)

func Test0(t *testing.T) {

	can := makeCanvas(50, 20)
	can.Print()
}
func Test1(t *testing.T) {
	can := makeCanvas(20, 20)
	fmt.Printf("line count:%d\nchar count:%d\n", len(can.lines), len(can.lines[0]))
	can.Print()
}
func TestLines(t *testing.T) {
	w := 20
	h := 20
	can := makeCanvas(w, h)
	//draw lines
	l := 5
	dirs := []byte{'l', 'r', 'u', 'd'}
	for _, dir := range dirs {
		putLine(w/2, h/2, dir, l, can, true)
	}
	can.Print()
	can.PrintBytes()
}
