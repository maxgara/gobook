package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

// html doc
type doc struct {
	style string
	g     *grid
}

// simple grid to contain markdown sub-elements using <div>s
type grid struct {
	md   [][]string //markdown
	cols int        //setting for columns per row
}

func (d *doc) startGrid(cols int) *grid {
	d.g = &grid{cols: cols}
	return d.g
}

// create new document loading css styles from file st
func newDoc(st string) *doc {
	style, err := os.ReadFile(st)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	return &doc{style: string(style)}
}

// add markdown content cell to grid
func (g *grid) add(el string) {
	ridx := len(g.md)
	//new row case
	if ridx == 0 || len(g.md[ridx-1]) >= g.cols {
		r := []string{el}
		g.md = append(g.md, r)
		return
	}
	//continue row case
	g.md[ridx-1] = append(g.md[ridx-1], el)
}

// write grid to w
func (g *grid) render(w io.Writer) {
	for _, r := range g.md {
		fmt.Fprint(w, rowstart)
		for _, cell := range r {
			fmt.Fprintf(w, `<div class="cell">%s</div>`, cell)
		}
		fmt.Fprint(w, rowend)
	}
}
func (g *grid) String() string {
	buff := bytes.Buffer{}
	g.render(&buff)
	return buff.String()
}
func (d *doc) render(w io.Writer) {
	fmt.Fprintf(w, docbase, d.style, d.g)
}

const (
	rowstart = `<div class="row">`
	rowend   = `</div>`
	docbase  = `<!DOCTYPE html><html><head><style>%s</style></head><body>%s</body></html>`
)
