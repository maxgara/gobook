package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

const STROKE_WIDTH = "0.1"
const style_file = "wbstyle.css"

var pal = []string{"#f5f5b0", "#154734", "#cfa3a8", "#6a7ba2", "#d6d6d6"}

type Node struct {
	name  string      //HTML element
	class string      //CSS class
	attrs [][2]string //HTML attribute headers
	par   *Node       //parent
	cdn   []*Node     //children
	text  string      //only for text nodes
}

// write node encoding to w
func (n *Node) Encode(w io.Writer) {
	var astr string //attribute string
	if n.class != "" {
		astr += " class=" + n.class + ""
	}
	for _, at := range n.attrs {
		k := at[0]
		v := at[1]
		astr += fmt.Sprintf(` %v="%v"`, k, v)
	}

	start := `<` + n.name + astr + `>`
	end := `</` + n.name + `>`
	w.Write([]byte(start))
	for _, c := range n.cdn {
		c.Encode(w)
	}
	w.Write([]byte(n.text))
	w.Write([]byte(end))
}

type docBuilder struct {
	loc                *Node
	root               *Node
	w                  io.Writer
	svgPoints          [][2]float64
	svgPlotLabels      []string //names for each data series, displayed in legend
	svgPlotLabelColors []string
	gridCols           int //number of gridElement elements allowed in a grid row
	colorIdx           int //used to make each data series (polyline) a different color
}

func newDocBuilder(w io.Writer) *docBuilder {
	root := &Node{name: "root"}
	return &docBuilder{loc: root, w: w, root: root}
}

func (b *docBuilder) startNode(name string, class string, attrs [][2]string) {
	n := &Node{name: name, class: class, attrs: make([][2]string, len(attrs))}
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
	name  string
	class string
	attrs [][2]string
}

const (
	NODE_DOC = iota
	NODE_TITLE
	NODE_GRID
	NODE_GRIDROW
	NODE_GRIDELEMENT
	NODE_SVG
	NODE_POLY
	NODE_LABELAXISX
	NODE_LABELAXISY
	NODE_LEGEND
	NODE_LEGENDITEM
)

// initialization args
var nodeinit = []NodeArgs{
	NODE_DOC:         {"html", "", nil},
	NODE_TITLE:       {"div", "title", nil},
	NODE_GRID:        {"div", "grid", nil},
	NODE_GRIDROW:     {"div", "grid-row", nil},
	NODE_GRIDELEMENT: {"div", "grid-elem", nil},
	NODE_SVG:         {"svg", "", [][2]string{{"preserveAspectRatio", "none"}, {"xmlns", "http://www.w3.org/2000/svg"}}},
	NODE_POLY:        {"polyline", "", [][2]string{{"stroke-width", string(STROKE_WIDTH)}, {"fill", "none"}}},
	NODE_LABELAXISX:  {"div", "labelaxis-y", nil},
	NODE_LABELAXISY:  {"div", "labelaxis-x", nil},
	NODE_LEGEND:      {"div", "legend", nil},
	NODE_LEGENDITEM:  {"div", "legend-item", nil},
}

// 	NODE_DOC: NodeArgs{"body", "", nil},
// 	NODE_TITLE:       {"div", "title", nil},

// }

// start a node of type nt with default vals and return a pointer to it
func initNode(b *docBuilder, nt int) *Node {
	init := nodeinit[nt]
	b.startNode(init.name, init.class, init.attrs)
	return b.loc
}

func (b *docBuilder) startDoc() {
	return // unnecessary function, but currently in use
}
func (b *docBuilder) endDoc() {
	style, err := os.ReadFile(style_file)
	if err != nil {
		log.Fatalf("error reading style file %v\n", style)
		return
	}
	b.w.Write([]byte("<html><head><style>"))
	b.w.Write(style)
	b.w.Write([]byte("</style></head><body>"))
	b.root.Encode(b.w)
	b.w.Write([]byte("</body>"))
}
func (b *docBuilder) writePageTitle(s string) {
	n := initNode(b, NODE_TITLE)
	n.text = s
	b.endNode()
}
func (b *docBuilder) writeLabels() {
	return // unnecessary function, but currently in use
}

func (b *docBuilder) startSVG(title string, bounds [4]float64, xl string, yl string) {
	b.startNode("div", "svg-container", nil) //group all svg-related nodes (titles, labels, etc.)
	initNode(b, NODE_TITLE)
	b.writeText(title)
	b.endNode()

	initNode(b, NODE_LABELAXISY)
	b.writeText(yl)
	b.endNode()

	svg := initNode(b, NODE_SVG)
	s := fmt.Sprintf("%f %f %f %f", bounds[0], bounds[1], bounds[2], bounds[3])
	svg.attrs = append(svg.attrs, [2]string{"viewbox", s})
	b.endNode()

	initNode(b, NODE_LABELAXISX)
	b.writeText(xl)
	b.endNode()

	b.loc = svg //come back to svg after writing other elements in svg-container
}

func (b *docBuilder) endSVG() {
	b.endNode()              //leave svg
	initNode(b, NODE_LEGEND) //labels for each polyline element
	for i, v := range b.svgPlotLabels {
		l := initNode(b, NODE_LEGENDITEM)
		colstr := fmt.Sprintf("color: %v", b.svgPlotLabelColors[i])
		l.attrs = append(l.attrs, [2]string{"style", colstr})
		b.writeText(v)
		b.endNode()
	}
	b.writeText("Legend")
	b.endNode() //leave legend
	b.endNode() //leave svg-container
}

// grid handles formatting, caller uses startGridElement(). caller does not need to build their own grid rows.
func (b *docBuilder) startGrid(colnum int) {
	b.gridCols = colnum
	initNode(b, NODE_GRID)
	b.startNode("div", "grid-row", nil)
}
func (b *docBuilder) endGrid() {
	b.endNode()
}
func (b *docBuilder) startGridElem() {
	idx := len(b.loc.cdn) //index into current grid-row
	if idx != 0 && idx%b.gridCols == 0 {
		b.endNode()
		initNode(b, NODE_GRIDROW)
	}
	initNode(b, NODE_GRIDELEMENT)
}
func (b *docBuilder) endGridElem() {
	b.endNode()
}
func (b *docBuilder) startPoly(lab string) {
	b.svgPlotLabels = append(b.svgPlotLabels, lab)
	n := initNode(b, NODE_POLY)
	col := pal[b.colorIdx]
	n.attrs = append(n.attrs, [2]string{"stroke", col})
	b.colorIdx = (b.colorIdx + 1) % len(pal)
}
func (b *docBuilder) endPoly() {
	var s string
	n := b.loc
	for _, v := range b.svgPoints {
		s += fmt.Sprintf("%v,%v ", v[0], v[1])
	}
	//set dynamic points, set dynamic color
	b.svgPoints = nil
	n.attrs = append(b.loc.attrs, [2]string{"points", s})
	col := pal[b.colorIdx]
	n.attrs = append(n.attrs, [2]string{"color", col})
	b.colorIdx++
	//apply color to label for data series
	cols := &b.svgPlotLabelColors
	*cols = append(*cols, pal[b.colorIdx])
	b.endNode()
}
func (b *docBuilder) vertex(x, y float64) {
	b.svgPoints = append(b.svgPoints, [2]float64{x, y})
}
func (b *docBuilder) writeText(s string) {
	b.loc.text += s
}
