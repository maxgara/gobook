package main

//functions to build an SVG based HTML document

import (
	"fmt"
	"io"
	"os"
)

var pal = pal2
var coolpal = []uint32{0xf5f5b0, 0x154734, 0xcfa3a8, 0x6a7ba2, 0xd6d6d6} // cool-tone color palette
var pal0 = []uint32{0xfe938c, 0xe6b89c, 0xead2ac, 0x9cafb7, 0x4281a4}
var pal1 = []uint32{0x40f99b, 0x20a4f3, 0x941c2f, 0x03191e, 0xfcab10}
var pal2 = []uint32{0xf5f5b0, 0x154734, 0xcfa3a8, 0x6a7ba2, 0xd6d6d6, 0xfe938c, 0xe6b89c, 0xead2ac, 0x9cafb7, 0x4281a4, 0x40f99b, 0x20a4f3, 0x941c2f, 0x03191e, 0xfcab10}

const style_file = "wbstyle.css"

var style string

// document is made up of a title, and grids of text, svg, and svg-label elements
type docBuilder struct {
	gridcols int //grid column count
	grididx  int
	cidx     int //palette color index
	w        io.Writer
}

func newDocBuilder(w io.Writer) *docBuilder {
	b, err := os.ReadFile(style_file)
	if err != nil {
		panic(err)
	}
	style = string(b)
	return &docBuilder{w: w}
}
func (d *docBuilder) startDoc() {
	d.writef(STARTDOC_FSTR, style)
}
func (d *docBuilder) endDoc() {
	d.writef(ENDDOC_FSTR)
}

// write a document portion to w
func (d *docBuilder) writef(fstr string, args ...any) {
	s := fmt.Sprintf(fstr, args...)
	d.w.Write([]byte(s))
}

func (d *docBuilder) writeTitle(s string) {
	d.writef("<div id=title>%v</div>", s)
}
func (d *docBuilder) startGrid(cols int) {
	d.gridcols = cols
	d.grididx = 0
	d.writef("<div class=grid><div class=grid-row>")
}
func (d *docBuilder) endGrid() {
	d.writef("</div></div>")
}
func (d *docBuilder) startGridElem() {
	if d.grididx != 0 && d.grididx%d.gridcols == 0 {
		d.writef("</div><div class=grid-row>")
	}
	d.writef(`<div class="grid-elem">`)
	d.grididx++
}
func (d *docBuilder) endGridElem() {
	d.writef(`</div>`)
}

// start SVG element with given title and viewbox bounds
func (d *docBuilder) startSVG(title string, view [4]float64) {
	d.writef(SVGSTART_FSTR, title, view[0], view[1], view[2], view[3])
}
func (d *docBuilder) endSVG() {
	d.writef(SVGEND_FSTR)
}
func (g *docBuilder) startPoly(width float64) {
	stroke := pal[g.cidx]
	g.cidx = (g.cidx + 1) % len(pal)
	g.writef(POLYSTART_FSTR, stroke, width)
}
func (g *docBuilder) endPoly() {
	g.writef(POLYEND_FSTR)
}
func (g *docBuilder) vertex(x, y float64) {
	g.writef("%v,%v ", x, y)
}

const (
	POLY_STROKE_WIDTH_DEFAULT = 0.5
	SVGSTART_FSTR             = `<div class="svg-container"><div class="svg-title">%v</div><svg viewBox="%v %v %v %v" preserveAspectRatio="none"
			xmlns="http://www.w3.org/2000/svg">`
	SVGEND_FSTR    = `</svg></div>`
	POLYSTART_FSTR = `<polyline stroke="#%x" fill="none" stroke-width="%v"
				points="`
	POLYEND_FSTR = `">
			</polyline>`
	STARTDOC_FSTR = "<!DOCTYPE HTML><html><head><style>%v</style></head>"
	ENDDOC_FSTR   = "</body></html>"
)
