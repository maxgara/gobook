// create svg graph from data points read from stdin. minimal version, no frills.
// reads data points as float64
package main

import (
	"fmt"
	"io"
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
	var data []byte
	data, _ = io.ReadAll(os.Stdin)
	svgs := parse(string(data))
	if len(svgs) == 0 {
		return
	}
	print(svgs)
}

// svg data
type svg struct {
	Curves                 []curve //polyline curve data
	Xmin, Xmax, Ymin, Ymax float64 //bounds for SVG viewbox
	Label                  string
	Xrange                 float64
	Yrange                 float64
	colors                 colorset
}

func newsvg() svg {
	return svg{Curves: make([]curve, 0), Xmin: math.MaxFloat64, Ymin: math.MaxFloat64, colors: colorset{}}
}
func (b *svg) Add(p point) {
	c := &b.Curves[len(b.Curves)-1]
	c.Add(p)
}
func (b *svg) Endc() {
	c := curve{Fill: b.colors.new()}
	b.Curves = append(b.Curves, c)
}

// polyline
type curve struct {
	P    []point //points on curve
	Fill color   //color of line
}

func (c *curve) Add(p point) {
	c.P = append(c.P, p)
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
	cs := colorset{}
	c.Fill = cs.new()
	for _, l := range strings.Split(data, "\n") {
		t := linetype(l)
		switch t {
		// start a new curve plot
		case EMPTY:
			if c.P == nil {
				continue
			}
			box.Curves = append(box.Curves, c)
			c = curve{}
			c.Fill = cs.new()
		// start a new curve plot and a new chart
		case NEWCHARTFLAG:
			if c.P == nil {
				continue
			}
			box.Curves = append(box.Curves, c)
			boxes = append(boxes, box)
			box = svg{}
			c = curve{}
			cs = colorset{}
			c.Fill = cs.new()
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
	if c.P != nil {
		box.Curves = append(box.Curves, c)
	}
	if box.Curves != nil {
		boxes = append(boxes, box)
	}
	return boxes
}
func print(boxes []svg) {
	templ, err := template.New("svg").Parse(`<svg viewBox="{{.Xmin}} {{.Xmax}} {{.Xrange}} {{.Yrange}}" style="width: 100%; height: 100%; display: flex" xmlns="http://www.w3.org/2000/svg">
	{{range .Curves}}
	<polyline stroke="{{.Fill}}" fill="none" stroke-width="0.7" points="{{range .P}} {{.X}},{{.Y}}{{end}}">
	{{end}}</svg>`)
	if err != nil {
		fmt.Printf("err:%v\n", err)
	}
	for _, b := range boxes {
		//write svg to stdout
		b.Xrange = b.Xmax - b.Xmin
		b.Yrange = b.Ymax - b.Ymin
		err = templ.Execute(os.Stdout, b)
		// fmt.Printf("curves:%v, %v\n", b.Curves, b.Curves)
		if err != nil {
			fmt.Printf("err:%v\n", err)
		}
	}
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
