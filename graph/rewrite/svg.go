package main

import (
	"fmt"
	"io"
)

// document is made up of a title, and grids of text, svg, and svg-label elements
type docBuilder struct {
	// gridcols   int //grid column count
	// grididx    int
	// cidx       int //palette color index
	// labelstack []string
	// xl         string
	// yl         string
	w   io.Writer
	loc *Node
}

// element in doc
type Node struct {
	name string
	attrs [][2]string //key value pairs
	chl []*Node //children
	t   int //node type
	c   NodeContent
	par *Node   //parent
}

type NodeContent interface {
	String()
}
type SVGBlock struct {
	title  string
	xl     string //axis label x
	yl     string //axis label y
	viewBox [4]float64
	series [][2][]float64
	labels []string
}
func (svg SVGBlock) String() string {
	viewBox := [4]float64
	preserveAspectRatio="none"
}
type Grid struct {
	rlen int
}
type TextBlock struct {
	s string
}
type Title struct {
	s string
}
type PageTitle struct {
	s string
}
func NodeString (n *Node) string{
	var astr string
	for i := range n.attrs{
		k := n.attrs[i][0]
		v := n.attrs[i][1]
		astr += fmt.Sprintf(` "%v"="%v"`,k,v)
	}
	var cstr string
	for _, ch := range n.chl{
		cstr += n.chl.NodeContent.String()
	}
	return fmt.Sprintf("<%v %v> %v </%1%v>", n.name, astr, cstr)
}

// svgFstr:=`<svg viewBox="%d %d %d %d" preserveAspectRatio="none"
                            xmlns="http://www.w3.org/2000/svg">
                            %v
                        </svg>`