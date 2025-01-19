package main

import (
	"fmt"
	"io"
)

const (
	REGULAR = iota
	TEXTNODE
	ROOT
)

type Node struct {
	name  string
	class string
	attrs [][2]string
	par   *Node
	cdn   []*Node
	text  string //only for text nodes
	t     int
}

// write node encoding to w
func (n *Node) Encode(w io.Writer) {
	var astr string //attribute string
	for _, at := range n.attrs {
		k := at[0]
		v := at[1]
		astr += fmt.Sprintf(" %v=%v", k, v)
	}

	start := `<` + n.name + astr + `>`
	end := `</` + n.name + `>`
	w.Write([]byte(start))
	if n.t == TEXTNODE {
		w.Write([]byte(n.text))
		w.Write([]byte(end))
		return
	}
	for _, c := range n.cdn {
		c.Encode(w)
	}
}

type docBuilder struct {
	loc  *Node
	root *Node
	w    io.Writer
}

func newDocBuilder(w io.Writer) *docBuilder {
	root := &Node{name: "root", t: ROOT}
	return &docBuilder{loc: root, w: w, root: root}
}

func (b *docBuilder) startNode(name string, class string, attrs [][2]string) {
	n := &Node{name: name, class: class}
	copy(n.attrs, attrs)
	b.loc.cdn = append(b.loc.cdn, n)
	n.par = b.loc
	b.loc = n
}

func (b *docBuilder) endNode() {
	b.loc = b.loc.par
}

// relevant Node Properties for functions below
type NodeArgs struct {
	func_id int
	name    string
	class   string
	attrs   [][2]string
}

/*func (d *docBuilder) writef(fstr string, args ...any) {
d.writef("<div class=labels>")
	d.writef(`<div id=label style="background-color: #%x">%v</div>`, pal[d.cidx], s)
d.writef("</div>")
d.writef(STARTDOC_FSTR, style)
d.writef(ENDDOC_FSTR)
d.writef(`<div class="text-block">%v</div>`, s)
d.writef(`<div class="title">%v</div>`, s)
d.writef(`<div class="pagetitle">%v</div>`, s)
d.writef(`<div class=grid><div class="grid-row">`)
d.writef("</div></div>")
	d.writef(`</div><div class="grid-row">`)
d.writef(`<div class="grid-elem">`)
d.writef(`</div>`)
d.writef(SVGSTART_FSTR, title, yaxis, view[0], view[1], view[2], view[3])
d.writef(SVGEND_FSTR, d.xl)
g.writef(POLYSTART_FSTR, stroke, width)
g.writef(POLYEND_FSTR)
g.writef("%v,%v ", x, y)*/

// initialization constants
var nodeinit = []NodeArgs{
	NodeArgs{NODE_DOC, "body", "", nil},
	NodeArgs{NODE_PAGETITLE, "div", "title", nil},
	NodeArgs{NODE_GRID, "div", "grid", nil},
	NodeArgs{NODE_GRIDROW, "div", "grid-row", nil},
	NodeArgs{NODE_GRIDITEM, "div", "grid-elem", nil},
	NodeArgs{NODE_SVG, "svg", "", [][2]string{[2]string{"preserveAspectRatio", "none"}, [2]string{"xmlns", "http://www.w3.org/2000/svg"}}},
	NodeArgs{NODE_POLY, "polyline", "", nil},
	NodeArgs{NODE_LABELAXISX, "div", "", nil},
	NodeArgs{NODE_LABELAXISY, "div", "", nil},
	NodeArgs{NODE_LEGEND, "div", "", nil},
	NodeArgs{NODE_LEGENDITEM, "div", "", nil}}

const (
	NODE_DOC = iota
	NODE_PAGETITLE
	NODE_GRID
	NODE_GRIDROW
	NODE_GRIDITEM
	NODE_SVG
	NODE_POLY
	NODE_LABELAXISX
	NODE_LABELAXISY
	NODE_LEGEND
	NODE_LEGENDITEM
)

func (b *docBuilder) startDoc() {
	return //unnecessary
}
func (b *docBuilder) endDoc() {
	b.root.Encode(b.w)
}
func (b *docBuilder) writePageTitle(s string) {
	b.startNode("div", "pagetitle")
}
func (b *docBuilder) writeLabels()
func (b *docBuilder) startSVG(title string, bounds [4]float64, xl string, yl string)
func (b *docBuilder) endSVG()
func (b *docBuilder) startGrid(colnum int)
func (b *docBuilder) endGrid()
func (b *docBuilder) startGridElem()
func (b *docBuilder) endGridElem()
func (b *docBuilder) startPoly(width float64, lab string)
func (b *docBuilder) endPoly()
func (b *docBuilder) vertex(x, y float64)
func (b *docBuilder) writeText(s string)
