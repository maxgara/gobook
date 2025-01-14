// Create SVG graphs of data.
// Usage: data is read from stdin, with data points separated by newlines. Output is to stdout and is comprised of
// HTML with embedded SVG graph(s).
// data points are one or more series, each in a separate column, separated by whitespace.
// data series may include a header row
// additional data series can be displayed by separating the series with non-numerical text or an empty line.
// additional empty lines between plots are ignored.

// The following flags are supported between data series:
// -pagetitle=<tstr>	set page title
// -title=<tstr>		set svg title
//
// -n 		New Chart	 By default curves are all plotted on the same chart. This starts a new chart
//					     which is used for further data series.
// -p 					parallel plot. place next svg next to previous one instead of below.
// -css=<style> 		add css style string(s) to next svg.
//
// Rough Format:
// [Title]
// [Graph Label]
// [Plot Label]
// <data>
//
//	...	(more data)
//	[-flag] or [labeltext] or [\n]
//
// [Graph Label]
// [Plot Label]
// <data> [data]
//
//	... (more data, new series)

package main

import (
	"flag"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

// type svg struct {
// 	sset                   [][]float64 //data series
// 	Xmin, Xmax, Ymin, Ymax float64     //bounds for SVG viewbox
// }

const POLYWIDTH = 0.1 //polyline stroke width

// data types
const (
	DATA = iota
	TEXT
	FLAGS
)

// flags
const (
	NEWCHART = 0b1 << iota
	PARALLEL
)

type parser struct {
	s         string
	lines     []string
	t         int //type of input processed
	dcols     int // # of data columns
	pagetitle string
	title     string //svg title
	text      string
	flags     int
	data      [][]float64
	err       error
	css       string
	readback  []byte
	secs      []sectionVars
}

func newParser(r io.Reader) parser {
	s, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return parser{s: string(s)}
}

func parseAllAndRender(r io.Reader) {
	p := newParser(r)
	for p.parse() {
		//continue
	}
	first := true
	var ingrid bool //do we need to close a grid?
	b := newDocBuilder(os.Stdout)
	for _, sec := range p.secs {
		if sec.gridFlag != 0 {
			ingrid = true
		}
		encodeSection(sec.text, sec.data, sec.d1Flag, sec.d2Flag, sec.dxFlag, sec.gridFlag, b, first, ingrid)
		first = false
	}
	endEncoding(b, ingrid)
}

// close html sections
func endEncoding(b *docBuilder, gfset bool) {
	if gfset {
		b.endGrid()
	}
	b.endDoc()
}

// parses a section of input into a sectionVars struct
// each call to parse reads until a break in the data series
// sets flags, titles, and text section properties for p
// returns false when parsing is complete, either due to error or end of input
func (p *parser) parse() bool {
	if p.lines == nil {
		p.lines = strings.Split(p.s, "\n")
	}
	lines := p.lines
	var lidx int
	var textdone, datadone bool
	var data [][]float64
	var flagstrs []string
	var done bool
	//loop over blocks
	for {
		if lidx >= len(lines) {
			done = true
			break
		}
		l := strings.Trim(lines[lidx], " \t\r")
		words := strings.Fields(l)
		//parse non-data first
		if !textdone {
			if len(l) == 0 {
				lidx++
				continue
			}
			if isdata(words) {
				textdone = true
				data = make([][]float64, len(words))
				continue //do not increment lidx so that l can be processed as data
			}
			if l[0] == '-' {
				flagstrs = append(flagstrs, words...)
				lidx++
				continue
			}
			p.text += l
			lidx++
			continue
		}
		// then parse data
		if !datadone {
			if len(l) == 0 {
				datadone = true
				break
			}
			for i, w := range words {
				d, err := strconv.ParseFloat(w, 64)
				if err != nil {
					datadone = true
					break
				}
				data[i] = append(data[i], d) //add ith data point to data series i
			}
			lidx++
			continue //next line
		}
	}
	//parse flags
	flags := flag.NewFlagSet("svgparseflags", flag.PanicOnError)
	d1Flag := flags.Bool("d1", true, "one coordinate series parsing")
	d2Flag := flags.Bool("d2", false, "two coordinate series parsing")
	dxFlag := flags.Bool("dx", false, "shared x value")
	gridFlag := flags.Int("grid", 0, "grid column count")

	flags.Parse(flagstrs)
	//handle result of parsing this block
	sec := sectionVars{text: p.text, data: data, d1Flag: *d1Flag, d2Flag: *d2Flag, dxFlag: *dxFlag, gridFlag: *gridFlag}
	p.secs = append(p.secs, sec)
	p.lines = p.lines[lidx:]
	if done {
		return false
	}
	return true
}

const (
	D1_FLAG = 0b1 << iota
	D2_FLAG
	DX_FLAG
)

type sectionVars struct {
	text     string
	data     [][]float64
	d1Flag   bool
	d2Flag   bool
	dxFlag   bool
	gridFlag int
}

// convert parsed data to document format
func encodeSection(text string, data [][]float64, d1Flag, d2Flag, dxFlag bool, gridFlag int, b *docBuilder, first bool, ingrid bool) {
	b.writeText(text)
	if len(data) == 0 {
		return
	}
	//assign pairs of data columns to contain x,y values of series based on flags
	var spairs [][2][]float64
	idxs := make([]float64, len(data[0])) //new index data series for convenience
	for i := range len(data[0]) {
		idxs[i] = float64(i)
	}
	switch {
	case d2Flag && !dxFlag:
		for i := 0; i < len(data); i += 2 {
			spair := [2][]float64{data[i], data[i+1]}
			spairs = append(spairs, spair)
		}
	case dxFlag:
		for i := 1; i < len(data); i++ {
			spair := [2][]float64{data[0], data[i]}
			spairs = append(spairs, spair)
		}
	case d1Flag:
		for _, s := range data {
			spair := [2][]float64{idxs, s}
			spairs = append(spairs, spair)
		}
	}
	//get bounds
	var xmin = math.MaxFloat64
	var ymin = math.MaxFloat64
	var xmax = -math.MaxFloat64
	var ymax = -math.MaxFloat64
	for _, s := range spairs {
		for _, x := range s[0] {
			xmin = min(xmin, x)
			xmax = max(xmax, x)
		}
		for _, y := range s[1] {
			ymin = min(ymin, y)
			ymax = max(ymax, y)
		}
	}
	//draw
	if first {
		b.startDoc()
	}
	if gridFlag != 0 {
		b.startGrid(gridFlag) // closed by endEncoding func
	}
	if ingrid {
		b.startGridElem()
	}
	bounds := [4]float64{xmin, ymin, xmax - xmin, ymax - ymin}
	b.startSVG("temporary svg title", bounds)
	for _, s := range spairs {
		b.startPoly(0.1)
		for i := range s[0] {
			x := s[0][i]
			y := s[1][i]
			b.vertex(x, y)
		}
		b.endPoly()
	}
	b.endSVG()
	if ingrid {
		b.endGridElem()
	}
}

func isdata(words []string) bool {
	datachars := []byte("-0123456789.E")
	for _, w := range words {
		for _, c := range []byte(w) {
			ok := false
			for _, dc := range datachars {
				if c == dc {
					ok = true
					break // next c
				}
			}
			if !ok {
				return false
			}
		}
	}

	return true
}

func main() {
	test := "-grid=2\n10 10\n20 10\n30 40\n\n0 0\n 1 2\n0 0\n-0.31 -0.5"
	parseAllAndRender(strings.NewReader(test))
}
