// create svg graph from data points read from stdin. minimal version, no frills.
// reads data points as float64
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const width = 600
const height = 500
const sep = " "

// const header = `<svg xmlns="http://www.w3.org/2000/svg" overflow="visible">`

// read lines of data, when hitting a non-data line check for flag -n indicating a new chart, otherwise add line to current chart
func main() {
	printSVGs("")
}
func printSVGs(instr string) string {
	// const linestart = `<polyline stroke="grey" fill="none" stroke-width="0.7" points="`
	const lineend = `"> </polyline>`
	// const svgstart = `<div style="width: 100%; height: 100%; display: flex; justify-content: center; align-items: center;"><svg width="100%" height="100%" viewBox="0 0 2000 2000" xmlns="http://www.w3.org/2000/svg" overflow="visible">`
	const svgend = `</svg></div>`
	var cols = colors{}
	var bounds = []float64{10000, 0, 10000, 0}
	var out string
	var outtemp string                         //store string until it is ready to be added to out
	outtemp += printlinestart(cols.newcolor()) //add line start but do not add svg start until we have bounds
	//read input
	var r *bufio.Scanner
	if instr != "" {
		sr := strings.NewReader(instr)
		r = bufio.NewScanner(sr)
	} else {
		r = bufio.NewScanner(os.Stdin)
	}
	var idx int // idx of current point in line plot
	for ok := r.Scan(); ok; ok = r.Scan() {
		// fmt.Printf(`scanned:"%s"`, r.Text())
		line := r.Text()
		t := linetype(line)
		//begin new line plot in same SVG
		if t == EMPTY {
			outtemp += lineend + printlinestart(cols.newcolor())
			idx = 0
			continue
		}
		//begin new SVG and new line plot
		if t == NEWCHARTFLAG {
			cols = colors{} //reset colors
			outtemp += lineend + printBoundsHTML(bounds) + svgend
			outtemp = printsvgstart(bounds) + outtemp // prepend svg start tag
			out += outtemp
			outtemp = printlinestart(cols.newcolor())
			idx = 0
			continue
		}
		p := parsepoint(line, idx)
		updatebounds(bounds, p[0], p[1])
		outtemp += printpoint(p)
		idx++
	}
	outtemp += lineend + svgend + printBoundsHTML(bounds)
	outtemp = printsvgstart(bounds) + outtemp // prepend svg start tag
	out += outtemp
	fmt.Print(out)
	return out
}

func printpoint(p []float64) string {
	return fmt.Sprintf(" %d,%d", int(p[0]), int(p[1]))
}
func printsvgstart(b []float64) string {
	width := b[XMAX] - b[XMIN]
	height := b[YMAX] - b[YMIN]
	return fmt.Sprintf(`<div style="width: %s; height: %s; display: flex"><svg viewBox="%d %d %d %d" style="width: %s; height: %s; display: flex" xmlns="http://www.w3.org/2000/svg">`, `100%`, `100%`, int(b[XMIN]), int(b[YMIN]), int(width), int(height), "100%", "100%")
}

func printlinestart(c color) string {
	return fmt.Sprintf(`<polyline stroke="%v" fill="none" stroke-width="0.7" points="`, c)
}
func printBoundsHTML(b []float64) string {
	return fmt.Sprintf("<div> XMIN=%v, XMAX=%v, YMIN=%v, YMAX=%v</div>", b[XMIN], b[XMAX], b[YMIN], b[YMAX])
}
func printBoundsSVG(bounds []float64) string {
	const off = 50
	b := []int{int(bounds[0]), int(bounds[1]), int(bounds[2]), int(bounds[3])}
	const style = `<style>".small {
      font: italic 13px sans-serif;
    }"</style>`
	return fmt.Sprintf(`%s<text x="%d" y="%d" class="small">%d,%d</text>
	 <text x="%d" y="%d" class="small">%d</text>
	  <text x="%d" y="%d" class="small">%d</text>
	  `, "",
		b[XMIN]+off, b[YMIN]+off, b[XMIN], b[YMIN], b[XMAX]-off, b[YMIN]+off, b[XMAX], b[XMIN]+off, b[YMAX]-off, b[YMAX])
}

type colors struct {
	allcolors []color
	used      map[color]bool
}

type color uint32

func (c color) String() string {
	return fmt.Sprintf("#%06x", uint32(c))
}

const red color = 0xff0000
const green color = 0x00ff00
const blue color = 0x0000ff

func blend(c1, c2 color) color {
	const rfilter = 0xff0000
	const gfilter = 0x00ff00
	const bfilter = 0x0000ff
	r := (((c1 & rfilter) + (c2 & rfilter)) / 2) & rfilter
	g := (((c1 & gfilter) + (c2 & gfilter)) / 2) & gfilter
	b := (((c1 & bfilter) + (c2 & bfilter)) / 2) & bfilter
	return r | g | b
}
func (c *colors) newcolor() color {
	//initialization case
	l := len(c.allcolors)
	if l == 0 {
		c.allcolors = []color{red, green, blue}
		c.used = make(map[color]bool)
		return 0x000000 //start with black
	}
	//out of colors case, add new colors
	if l == len(c.used) {
		var cnew []color
		// for each pair of colors a_n, a_n+1 insert new blended color in the middle
		for i, _ := range c.allcolors {
			nc := blend(c.allcolors[i], c.allcolors[(i+1)%l]) //new color
			cnew = append(cnew, c.allcolors[i])               //original color
			cnew = append(cnew, nc)
		}
		c.allcolors = cnew
		return c.newcolor()
	}
	for _, col := range c.allcolors {
		if !c.used[col] {
			c.used[col] = true
			return col
		}
	}
	fmt.Println("hit end of newcolor() - that's not supposed to happen")
	return red
}

// types of per-line input
const (
	DATA int = iota
	EMPTY
	NEWCHARTFLAG
	LABEL
)

// decide what type of input a given line read is. default is DATA for coordinates/datapoints
func linetype(s string) int {
	if s == "" {
		return EMPTY
	}
	if s == "-n" {
		return NEWCHARTFLAG
	}
	if strings.ContainsAny(s, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return LABEL
	}
	return DATA
}

// convert string into data point
func parsepoint(s string, idx int) []float64 {
	//break s into field strings
	p := strings.Fields(s)
	if len(p) == 0 {
		panic("failed to parse coordinates")
	}
	//convert each field string to a float
	var pfloat []float64
	for _, pstr := range p {
		f, _ := strconv.ParseFloat(pstr, 64)
		pfloat = append(pfloat, f)
	}
	//add idx as the x coord field if missing.
	if len(pfloat) == 1 {
		return []float64{float64(idx), pfloat[0]}
	}
	return pfloat
}

// b = {xmin, xmax, ymin, ymax}
const (
	XMIN int = iota
	XMAX
	YMIN
	YMAX
)

func updatebounds(b []float64, x float64, y float64) {

	if x < b[XMIN] {
		b[XMIN] = x
	}
	if x > b[XMAX] {
		b[XMAX] = x
	}
	if y < b[YMIN] {
		b[YMIN] = y
	}
	if y > b[YMAX] {
		b[YMAX] = y
	}
}
func scalesvg(bounds []float64) {

}
