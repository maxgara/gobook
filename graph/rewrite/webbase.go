package main

import (
	"fmt"
	"io"
	"os"
)

var style []byte //load from file (or embed?)

// html doc
type doc struct {
	md string //markdown
}

// simple grid to contain markdown elements
type grid struct {
	md   [][]string //markdown
	cols int        //setting for columns per row
}

func newGrid(cols int) *grid {
	return &grid{cols: cols}
}
func newDoc() *doc {
	var err error
	style, err = os.ReadFile("wbstyle.css") //load styles
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return &doc{}
}
func (g *grid) add(el string) {
	rowidx := len(g.md)
	//new row case
	if rowidx == 0 || len(g.md[rowidx-1]) >= g.cols {
		r := []string{el}
		g.md = append(g.md, r)
		return
	}
	//continue row case
	g.md[rowidx-1] = append(g.md[rowidx-1], el)
}
func (g *grid) render(w io.Writer) {
	for _, r := range g.md {
		fmt.Fprint(w, rowstart)
		for _, cell := range r {
			fmt.Fprintf(w, `<div class="cell">%v</div>\n`, cell)
		}
		fmt.Fprint(w, rowend)
	}
}
func (d *doc) render(w io.Writer) {
	fmt.Fprintf(w, docbase, style)
}

const (
	rowstart = `<div class="row">`
	rowend   = `</div>`
	docbase  = `<html><head>%v</head><body></body></html>`
)
