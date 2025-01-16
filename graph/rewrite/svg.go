package main

import (
	"fmt"
)

// global consts + vars

// colors for svg generation
var pal = []uint32{0xf5f5b0, 0x154734, 0xcfa3a8, 0x6a7ba2, 0xd6d6d6, 0xfe938c, 0xe6b89c, 0xead2ac, 0x9cafb7, 0x4281a4, 0x40f99b, 0x20a4f3, 0x941c2f, 0x03191e, 0xfcab10}
var style string

const SVG_STROKEWIDTH = 0.1 //default stroke width
const style_file = "wbstyle.css"

// implementation of a general <div> as Content.
type Div struct {
	s string
	n Node
}

// div uses default node rendering with no modifications
func (d *Div) String() string {
	return NodeStringC(&d.n, d.s+fmt.Sprint(d.n.chl))
}
func (d *Div) Node() *Node {
	return &d.n
}
func (b *docBuilder) StartDiv() {
	b.StartContentNode(&Div{}, Node{name: "div"})
}
func (b *docBuilder) EndDiv() {
	b.EndNode()
}

// SVGBlock implements an SVG based graph involving multiple markup elements
type SVGBlock struct {
	title   string         //graph title
	xl      string         //axis label x
	yl      string         //axis label y
	viewBox [4]float64     //xmin, ymin, width, height
	series  [][2][]float64 //polyline, x|y, datapoint idx
	colors  []uint32       //RGB
	pidx    int            //color palette idx
	labels  []string       //series data labels
	n       Node
}

func (svg *SVGBlock) Node() *Node {
	return &svg.n
}

// default Node definitions
var (
	defaultSVGNode  = Node{name: "div", class: "svg-container", attrs: [][2]string{{"preserveAspectRatio", "none"}, {"xmlns", "http://www.w3.org/2000/svg"}}}
	defaultGridNode = Node{name: "div", class: "grid"}
	defaultLegend   = Node{name: "div", class: "labels"}
)

// enter SVG Bloc
func (b *docBuilder) StartSVGBlock(title, xl, yl string) {
	newsvg := SVGBlock{n: defaultSVGNode}
	newsvg.n.c = &newsvg
	b.StartNode(&newsvg.n)
}
func (b *docBuilder) EndSVGBlock() {
	b.loc = b.loc.par
}

// start new polyline with given stroke width and label
func (b *docBuilder) startPoly(label string) {
	c := b.loc.c.(*SVGBlock)
	color := pal[c.pidx]
	c.labels = append(c.labels, label)
	c.colors = append(c.colors, color)
	c.pidx++
	c.pidx = c.pidx % len(pal)
	c.series = append(c.series, [2][]float64{})
	//write opening tag
	// b.writef(`<polyline stroke="#%x" fill="none" stroke-width="%v" points="`, color, width)
}
func (b *docBuilder) vertex(x, y float64) {
	c := b.loc.c.(*SVGBlock)
	i := len(c.series) - 1
	c.series[i][0] = append(c.series[i][0], x)
	c.series[i][1] = append(c.series[i][1], x)
}

// write end tag, leave element
func (b *docBuilder) endPoly() {
	// b.writef("></polyline>")
	b.EndNode()
}
func (svg SVGBlock) String() string {
	boundstr := fmt.Sprintf("%v %v %v %v", svg.viewBox[0], svg.viewBox[1], svg.viewBox[2], svg.viewBox[3])
	svg.n.attrs = append(svg.n.attrs, [2]string{"viewBox", boundstr})
	var b buffer //content string buffer
	for sidx, ser := range svg.series {
		b.writef(`<polyline stroke="#%x" fill="none" stroke-width="%v" points="`, svg.colors[sidx], SVG_STROKEWIDTH)
		for i := range ser[0] {
			x := ser[0][i]
			y := ser[1][i]
			b.writef("%v,%v ", x, y)
		}
		b.writef(`"></polyline>`)
	}
	return NodeStringC(&svg.n, b.String())
}

type Grid struct {
	rlen int
}
type TextBlock struct {
	s string
	n Node
}

func (tb TextBlock) String() string {
	return NodeString(&tb.n)
}

type Title struct {
	s string
}
type PageTitle struct {
	s string
}

// NodeString but with specified inner HTML string
func NodeStringC(n *Node, cs string) string {
	var as string
	for i := range n.attrs {
		k := n.attrs[i][0]
		v := n.attrs[i][1]
		as += fmt.Sprintf(` "%v"="%v"`, k, v)
	}
	return fmt.Sprintf("<%v%v>%v</%[1]v>", n.name, as, cs)
}

// default html node printing, no text content
func NodeString(n *Node) string {
	var as string
	for i := range n.attrs {
		k := n.attrs[i][0]
		v := n.attrs[i][1]
		as += fmt.Sprintf(` "%v"="%v"`, k, v)
	}
	var cs string
	for _, cd := range n.chl {
		cs += cd.String()
	}
	return fmt.Sprintf("<%v%v>%v</%[1]v>", n.name, as, cs)
}

// // svgFstr:=`<svg viewBox="%d %d %d %d" preserveAspectRatio="none"
//                             xmlns="http://www.w3.org/2000/svg">
//                             %v
//                         </svg>`
