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
	gridcols   int //grid column count
	grididx    int
	cidx       int //palette color index
	labelstack []string
	xl         string
	yl         string
	w          io.Writer
}

const (
	RMASK = 0xff0000 >> (8 * iota)
	GMASK
	BMASK
)

// darken rgb color by given percent
func darken(col uint32, per uint32) uint32 {
	adj := func(k uint32) uint32 {
		return k - (per*k)/100
	}
	var newcol uint32
	for i := range 3 {
		bsh := 8 * i
		bm := uint32(0xff0000) >> bsh
		cmp := col & bm >> (16 - bsh)
		c := adj(cmp) << (16 - bsh)
		newcol |= c
	}
	fmt.Printf("adjusted %x to %x (%v%%)\n", col, newcol, per)
	return newcol
}
func newDocBuilder(w io.Writer) *docBuilder {
	b, err := os.ReadFile(style_file)
	if err != nil {
		panic(err)
	}
	style = string(b)
	return &docBuilder{w: w}
}

// write format string with args to w; svg.go internal use
func (d *docBuilder) writef(fstr string, args ...any) {
	s := fmt.Sprintf(fstr, args...)
	d.w.Write([]byte(s))
}

// write a block of labels below svg
func (d *docBuilder) writeLabels() {
	d.writef("<div class=labels>")
	//reset cidx to match color at the bottom of the label stack
	d.cidx = d.cidx - len(d.labelstack)
	for d.cidx < 0 {
		d.cidx += len(pal)
	}
	//write labels
	for _, s := range d.labelstack {
		d.writef(`<div id=label style="background-color: #%x">%v</div>`, pal[d.cidx], s)
		d.cidx++
	}
	d.writef("</div>")
	d.labelstack = nil
}
func (d *docBuilder) startDoc() {
	d.writef(STARTDOC_FSTR, style)
}
func (d *docBuilder) endDoc() {
	d.writef(ENDDOC_FSTR)
}
func (d *docBuilder) writeText(s string) {
	d.writef(`<div class="text-block">%v</div>`, s)
}
func (d *docBuilder) writeTitle(s string) {
	d.writef(`<div class="title">%v</div>`, s)
}
func (d *docBuilder) writePageTitle(s string) {
	d.writef(`<div class="pagetitle">%v</div>`, s)
}
func (d *docBuilder) startGrid(cols int) {
	d.gridcols = cols
	d.grididx = 0
	d.writef(`<div class=grid><div class="grid-row">`)
}
func (d *docBuilder) endGrid() {
	d.writef("</div></div>")
}
func (d *docBuilder) startGridElem() {
	if d.grididx != 0 && d.grididx%d.gridcols == 0 {
		d.writef(`</div><div class="grid-row">`)
	}
	d.writef(`<div class="grid-elem">`)
	d.grididx++
}
func (d *docBuilder) endGridElem() {
	d.writef(`</div>`)
}

// start SVG element with given title and viewbox bounds
func (d *docBuilder) startSVG(title string, view [4]float64, xaxis string, yaxis string) {
	d.yl = yaxis
	d.xl = xaxis
	d.writef(SVGSTART_FSTR, title, yaxis, view[0], view[1], view[2], view[3])
}
func (d *docBuilder) endSVG() {
	d.writef(SVGEND_FSTR, d.xl)
}

// start a polyline with given stroke width and label, which can be written later
func (g *docBuilder) startPoly(width float64, label string) {
	g.labelstack = append(g.labelstack, label)
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
	SVGSTART_FSTR             = `<div class="svg-container"><div class="svg-title">%v</div><div class="yaxislabel">%v</div><div class="svg-axisalignment"><svg viewBox="%v %v %v %v" preserveAspectRatio="none"
			xmlns="http://www.w3.org/2000/svg">`
	SVGEND_FSTR    = `</svg><div class="xaxislabel">%v</div></div></div>`
	POLYSTART_FSTR = `<polyline stroke="#%x" fill="none" stroke-width="%v"
				points="`
	POLYEND_FSTR = `">
			</polyline>`
	STARTDOC_FSTR = "<!DOCTYPE HTML><html><head><style>%v</style></head>"
	ENDDOC_FSTR   = "</body></html>"
)
