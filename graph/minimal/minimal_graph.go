package main

import (
	"embed"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"text/template"
)

const styles = `

<head>
	<style>
		body {
			font-family: monospace;
			font-size: 16;
			font-weight: 700;
			text-transform: lowercase;
		}
	</style>
</head>
`

var pageTitle string

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
func (b *svg) Push(p point) {
	if len(b.Curves) == 0 {
		b.NewCurve()
	}
	cc := len(b.Curves) // curve count
	c := &(b.Curves[cc-1])
	c.P = append(c.P, p)
}

// initialize curve in box b
func (b *svg) NewCurve() {
	cc := len(b.Curves)
	label := fmt.Sprintf("Plot %d", cc+1) //default naming
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

// polyline curve datapoint
type point struct {
	X float64
	Y float64
}

// read string input data points for 1 or more SVG elements, containing one or more curves each
func parse(data string) []svg {
	var boxes []svg
	var box = newsvg()
	var pidx int // index of point in curve
	var cidx int
	var bidx int
	var lidx int
	var onData bool //true=currently reading data series; false=between data series
	var newChart bool
	var sideBySide bool
	var labels []string
	parseData := func(l string) {
		t := linetype(l)
		switch t {
		case DATA:
			p, err := parsep(l)
			if err != nil {
				panic(fmt.Sprintf("error on line \"%s\"\n", l))
			}
			//add idx as x coord if missing
			if math.IsNaN(p.X) {
				p.X = float64(pidx)
			}
			box.Push(p)
			box.Xmin = min(p.X, box.Xmin)
			box.Xmax = max(p.X, box.Xmax)
			box.Ymin = min(p.Y, box.Ymin)
			box.Ymax = max(p.Y, box.Ymax)
			pidx++
		default:
			onData = false
			lidx--
		}
	}
	//call at beginning of new data series
	applyFlags := func() {
		if newChart {
			boxes = append(boxes, box)
			box = newsvg()
			cidx = 0
			bidx++
		}
		if !sideBySide {
			pidx = 0
		}
		box.NewCurve() // order important here
		cidx++
		for _, lab := range labels {
			if lab == "" {
				continue
			}
			switch {
			//page title case
			case bidx == 0 && pageTitle == "":
				pageTitle = lab
			//box label case
			case cidx == 1 && box.Label == "":
				box.Label = lab
				//curve label case
			default:
				cc := len(box.Curves)
				c := &box.Curves[cc-1]
				c.Label = lab
			}
		}
		newChart = false
		sideBySide = false
		labels = []string{}
	}
	parseNonData := func(l string) {
		t := linetype(l)
		// new box
		switch t {
		case NEWCHARTFLAG:
			//ignore newchartflag if current chart has only 1 curve containing 0 points.
			cc := len(box.Curves)
			if cc == 0 || len(box.Curves[cc-1].P) == 0 {
				return
			}
			newChart = true
		case TEXT:
			//assign label to element
			labels = append(labels, l)
		case PARALLELFLAG:
			sideBySide = true
			newChart = false //overrides previous flag if set
		case DATA:
			onData = true
			applyFlags()
			parseData(l)
		}
	}

	//main loop
	lines := strings.Split(data, "\n")
	lmax := len(lines)
	for lidx < lmax {
		l := lines[lidx]
		if onData {
			parseData(l)
		} else {
			parseNonData(l)
		}
		lidx++
	}
	boxes = append(boxes, box)
	return boxes
}

//go:embed svg.tmpl
var templates embed.FS

func print(boxes []svg) {
	templ, err := template.ParseFS(templates, "svg.tmpl")
	// templ, err := template.New("svg").Parse(tstr)
	if err != nil {
		fmt.Printf("err:%v\n", err)
	}
	fmt.Printf(`<h1 style="background-color: cornflowerblue;text-align: center;">%v</h1>`, pageTitle)
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
	NEWCHARTFLAG
	PARALLELFLAG
	TEXT
)

// decide what type of input a given line read is. default is TEXT
func linetype(s string) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return TEXT
	}
	s0 := s[0]
	switch {
	case s == "-n":
		return NEWCHARTFLAG
	case s == "-p":
		return PARALLELFLAG
	case (s0 >= '0' && s0 <= '9') || s0 == '-':
		return DATA
	default:
		return TEXT
	}
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
