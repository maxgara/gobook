// create svg graph from data points read from stdin. minimal version, no frills.
// reads data points as float64
package main

import (
	"fmt"
	"io"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

const styles = `

<head>
	<style>
		div {
			font-family: monospace;
			font-size: 16;
			font-weight: 700;
		}

		body {
			background-color: #FDFDF0;
		}
	</style>
</head>`

// read lines of data, when hitting a non-data line check for flag -n indicating a new chart, otherwise add line to current chart
func main() {
	//read all data in
	var data []byte
	data, _ = io.ReadAll(os.Stdin)
	svgs := parse(string(data))
	if len(svgs) == 0 {
		return
	}
	fmt.Printf("%s", styles)
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

// add point to current curve in svg
func (b *svg) Add(p point) {
	if len(b.Curves) == 0 {
		b.Newc()
	}
	cc := len(b.Curves) // curve count
	c := &(b.Curves[cc-1])
	c.Add(p)
}

// initialize curve in box b
func (b *svg) Newc() {
	cc := len(b.Curves)
	label := fmt.Sprintf("Plot %d", cc+1)
	c := curve{Col: b.colors.new(), Label: label}
	b.Curves = append(b.Curves, c)
}

// initialize new svg
func newsvg() svg {
	return svg{Curves: make([]curve, 0), Xmin: math.MaxFloat64, Ymin: math.MaxFloat64, Xmax: -math.MaxFloat64, Ymax: -math.MaxFloat64, colors: colorset{}}
}

// polyline curve data
type curve struct {
	P     []point //points on curve
	Col   color   //color of line
	Label string  //label for curve
}

// add point to polyline curve
func (c *curve) Add(p point) {
	c.P = append(c.P, p)
}

// polyline curve datapoint
type point struct {
	X float64
	Y float64
}

// read string input data points for 1 or more SVG elements, containing one or more curves each
func parse(data string) []svg {
	var boxes []svg
	var box = newsvg()
	var idx int // index of point in curve
	for _, l := range strings.Split(data, "\n") {
		t := linetype(l)
		switch t {
		//new curve plot in box
		case EMPTY:
			box.Newc()
		// new box
		case NEWCHARTFLAG:
			boxes = append(boxes, box)
			box = newsvg()
		//add labels
		case TEXT:
			//box label
			cc := len(box.Curves)
			if cc == 0 {
				box.Label = l
				continue
			}
			//curve label
			box.Curves[cc-1].Label = l
		case DATA:
			p, err := parsep(l)
			if err != nil {
				fmt.Printf("error on line \"%s\"\n", l)
				continue
			}
			//add idx as x coord if missing
			if math.IsNaN(p.X) {
				p.X = float64(idx)
			}

			box.Add(p)
			box.Xmin = min(p.X, box.Xmin)
			box.Xmax = max(p.X, box.Xmax)
			box.Ymin = min(p.Y, box.Ymin)
			box.Ymax = max(p.Y, box.Ymax)
		}
	}
	boxes = append(boxes, box)
	return boxes
}

func print(boxes []svg) {
	templ, err := template.ParseFiles("svg.tmpl")
	// templ, err := template.New("svg").Parse(tstr)
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
	s = strings.TrimSpace(s)
	if s == "" {
		return EMPTY
	}
	if s == "-n" {
		return NEWCHARTFLAG
	}
	if m, _ := regexp.MatchString(`[^\d\.\-\s]`, s); m {
		return TEXT
	}
	return DATA
}

// parse coordinates for point p, returns p and number of coords parsed, or -1 for error
func parsep(s string) (point, error) {
	words := strings.Fields(s)
	// wl := len(words)
	if len(words) != 1 && len(words) != 2 {
		err := fmt.Errorf("parsep: bad data line \"%v\"", s)
		// os.Exit(1)
		return point{}, err

	}
	//one coord case
	if len(words) == 1 {
		y, err := strconv.ParseFloat(words[0], 64)
		return point{math.NaN(), y}, err
	}
	// two coordinate case
	x, xerr := strconv.ParseFloat(words[0], 64)
	y, yerr := strconv.ParseFloat(words[1], 64)
	if xerr != nil {
		return point{x, y}, xerr
	}
	return point{x, y}, yerr
}
