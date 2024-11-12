// create svg graph from data points read from stdin. minimal version, no frills.
// reads data points as float64
package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"text/template"
)

// const header = `<svg xmlns="http://www.w3.org/2000/svg" overflow="visible">`
//rewrite to store data before formatting (simpler), and not use so much string concatination. Use templates instead

// read lines of data, when hitting a non-data line check for flag -n indicating a new chart, otherwise add line to current chart
func main() {
	//read all data in
	var data string
	data = io.ReadAll(os.Stdin)
	svgs := parse(data)
	if len(svgs) == 0 {
		return
	}
	// fmt.Println(svgs)
	print(svgs)
}

// svg data
type svg struct {
	Curves                 []curve //polyline curve data
	Xmin, Xmax, Ymin, Ymax float64 //bounds for SVG viewbox
	Label                  string
	Xrange                 float64
	Yrange                 float64
}

func newsvg() svg {
	return svg{Xmin: math.MaxFloat64, Ymin: math.MaxFloat64}
}

// polyline
type curve struct {
	P []point //points on curve
}

// datapoint
type point struct {
	X float64
	Y float64
}

// read string input input datapoints for 1 or more curves and parse into 1 or more SVG elements
func parse(data string) []svg {
	var boxes []svg
	var box = newsvg()
	var c curve
	for _, l := range strings.Split(data, "\n") {
		t := linetype(l)
		switch t {
		// start a new curve plot
		case EMPTY:
			box.Curves = append(box.Curves, c)
			c = curve{}
		// start a new curve plot and a new chart
		case NEWCHARTFLAG:
			box.Curves = append(box.Curves, c)
			boxes = append(boxes, box)
			box = svg{}
			c = curve{}
		//add label to plot if it appears before data
		case TEXT:
			if box.Label != "" || c.P != nil {
				continue
			}
			box.Label = l
		case DATA:
			p, n := parsep(l)
			if n == -1 {
				fmt.Printf("error on line %v\n", l)
			}
			//add idx as x coord if missing
			if n == 1 {
				p.X = float64(len(c.P))
			}
			c.P = append(c.P, p)
			box.Xmin = min(p.X, box.Xmin)
			box.Xmax = max(p.X, box.Xmax)
			box.Ymin = min(p.Y, box.Ymin)
			box.Ymax = max(p.Y, box.Ymax)
		}
	}
	box.Curves = append(box.Curves, c)
	boxes = append(boxes, box)
	return boxes
}
func print(boxes []svg) {
	templ, err := template.New("svg").Parse(`<svg viewBox="{{.Xmin}} {{.Xmax}} {{.Xrange}} {{.Yrange}}" style="width: 100%; height: 100%; display: flex" xmlns="http://www.w3.org/2000/svg">
	{{range .Curves}}
	<polyline stroke="%v" fill="none" stroke-width="0.7" points="{{range .P}} {{.X}},{{.Y}}{{end}}">
	{{end}}</svg>`)
	if err != nil {
		fmt.Printf("err:%v\n", err)
	}
	for _, b := range boxes {
		//write svg to stdout
		err = templ.Execute(os.Stdout, b)
		// fmt.Printf("curves:%v, %v\n", b.Curves, b.Curves)
		if err != nil {
			fmt.Printf("err:%v\n", err)
		}
	}
}

type colorset struct {
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
func (c *colorset) newcolor() color {
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
	TEXT
)

// decide what type of input a given line read is. default is DATA for coordinates/datapoints
func linetype(s string) int {
	if s == "" {
		return EMPTY
	}
	if strings.TrimSpace(s) == "-n" {
		return NEWCHARTFLAG
	}
	if strings.ContainsAny(s, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return TEXT
	}
	return DATA
}

// parse coordinates for point p, returns p and number of coords parsed, or -1 for error
func parsep(s string) (point, int) {
	words := strings.Fields(s)
	// wl := len(words)
	if len(words) != 1 && len(words) != 2 {
		fmt.Fprintf(os.Stderr, `ERROR: bad data line :"%v"\n`, s)
		os.Exit(1)
		return point{}, -1
	}
	//one coord case
	if len(words) == 1 {
		y, err := strconv.ParseFloat(words[0], 64)
		if err != nil {
			fmt.Fprintf(os.Stdout, "ERROR:%v\n", err)
			return point{}, -1
		}
		return point{0, y}, 1
	}
	// two coordinate case
	x, xerr := strconv.ParseFloat(words[0], 64)
	y, yerr := strconv.ParseFloat(words[1], 64)
	if xerr != nil || yerr != nil {
		err := xerr.Error() + yerr.Error()
		fmt.Printf("ERROR:%v\n", err)
		return point{0, y}, -1
	}
	return point{x, y}, 2
}

