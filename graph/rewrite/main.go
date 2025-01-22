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
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

const POLYWIDTH = 0.1 //polyline stroke width

type parser struct {
	s     string
	lines []string
	sec   section
}

func newParser(r io.Reader) parser {
	s, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(s), "\n")
	return parser{s: string(s), lines: lines}
}

// parse all input from r and print document to stdout
func parseAllAndRender(r io.Reader) {
	p := newParser(r)
	var secs []section
	for p.parse() {
		secs = append(secs, p.sec)
	}
	first := true
	var ingrid bool //do we need to close a grid?
	b := newDocBuilder(os.Stdout)
	for _, sec := range secs {
		if sec.flags.grid != 0 {
			ingrid = true
		}
		encodeSection(sec.text, sec.data, sec.headers, sec.flags, b, first, ingrid)
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

// read specified line
func (p *parser) readline(lidx int) (l string, words []string, ok bool) {
	if lidx >= len(p.lines) || lidx < 0 {
		return "", nil, false
	}
	l = strings.Trim(p.lines[lidx], " \t\r")
	words = strings.Fields(l)
	ok = true
	return
}

// parses a section of input into a section struct
// each call to parse reads until a break in the data series(s)
// sets flags, titles, and text section properties for p.sec
// returns false only after input string has been completely parsed.
// if parse() returns false, no data is present in p.sec and therefore
// the last successful call to parse will be after data is already done being parsed and len(p.lines) == 0.
func (p *parser) parse() (ok bool) {
	if len(p.lines) == 0 {
		return false
	}
	var lidx int
	var data [][]float64
	var text []string
	var flagstrs []string
	var done bool // end of input
	//parse non-data first
TLOOP:
	for !done {
		l, words, ok := p.readline(lidx)
		switch {
		case !ok:
			done = true
		case len(l) == 0:
			lidx++ //blank lines ignored in text mode
		case isdata(words):
			data = make([][]float64, len(words))
			break TLOOP // do not increment lidx, allow processing of l as data
		case l[0] == '-':
			flagstrs = append(flagstrs, words...)
			lidx++
		default:
			text = append(text, l)
			lidx++
		}
	}
	hlidx := lidx - 1      //header line idx
	htidx := len(text) - 1 //header text idx
	// then parse data
DLOOP:
	for !done {
		l, words, ok := p.readline(lidx)
		switch {
		case !ok:
			done = true
		case len(l) == 0:
			break DLOOP //blank lines mean end of series in data mode
		case !isdata(words):
			break DLOOP // do not increment lidx, allow processing of l as text
		default:
			for i, w := range words {
				d, err := strconv.ParseFloat(w, 64)
				if err != nil {
					break DLOOP // same as !isdata case
				}
				data[i] = append(data[i], d)
			}
			lidx++
		}
	}
	//parse headers if applicable
	var headers []string
	l, h, ok := p.readline(hlidx)
	if ok && len(l) != 0 && l[0] != '-' && len(data) != 0 {
		headers = h
		//drop headers from text
		haft := text[htidx:]
		if len(haft) > 1 {
			copy(haft, haft[1:])
		}
		text = text[:len(text)-1]
	}
	//parse collected flags
	flags := flag.NewFlagSet("svgparseflags", flag.PanicOnError)
	sf := sectionFlags{}
	flags.BoolVar(&sf.d1, "d1", true, "one coordinate series parsing")
	flags.BoolVar(&sf.d2, "d2", false, "two coordinate series parsing")
	flags.BoolVar(&sf.dx, "dx", false, "shared x value")
	flags.IntVar(&sf.grid, "grid", 0, "grid column count")
	flags.StringVar(&sf.title, "title", "", "section title")
	flags.StringVar(&sf.pageTitle, "pagetitle", "Graph Output", "page title")
	flags.StringVar(&sf.xl, "xlabel", "x-axis", "x axis label")
	flags.StringVar(&sf.yl, "ylabel", "y-axis", "y axis label")

	flags.Parse(flagstrs)

	//return parsed section
	p.sec = section{text: strings.Join(text, "\n"), data: data, headers: headers, flags: sf}
	if done {
		p.lines = nil
	} else {
		p.lines = p.lines[lidx:]
	}
	return true //ok
}

type sectionFlags struct {
	d1        bool
	d2        bool
	dx        bool
	grid      int
	title     string
	xl        string //x axis label
	yl        string //y axis label
	pageTitle string
}
type section struct {
	text    string
	data    [][]float64
	headers []string
	flags   sectionFlags
}

// convert parsed data to document format
func encodeSection(text string, data [][]float64, h []string, f sectionFlags, b *docBuilder, first bool, ingrid bool) {
	if first {
		b.startDoc()
		b.writePageTitle(f.pageTitle)
	}
	if f.grid != 0 {
		b.startGrid(f.grid) // closed by endEncoding func
		ingrid = true
	}
	if ingrid {
		b.startGridElem()
	}
	// if f.title != "" {
	// 	b.writeTitle(f.title)
	// }
	b.writeText(text)
	if len(data) != 0 {
		encodeSVG(data, h, f, b)
	}
	if ingrid {
		b.endGridElem()
	}
}

// build the SVG document, including svg-title and plot legend
func encodeSVG(data [][]float64, h []string, f sectionFlags, b *docBuilder) {

	//assign pairs of data columns to contain x,y values of series based on flags
	var spairs [][2][]float64
	if len(data) == 0 {

	}
	idxs := make([]float64, len(data[0])) //new index data series for convenience
	for i := range len(data[0]) {
		idxs[i] = float64(i)
	}
	switch {
	case f.d2 && !f.dx:
		for i := 0; i < len(data); i += 2 {
			spair := [2][]float64{data[i], data[i+1]}
			spairs = append(spairs, spair)
		}
	case f.dx:
		for i := 1; i < len(data); i++ {
			spair := [2][]float64{data[0], data[i]}
			spairs = append(spairs, spair)
		}
	case f.d1:
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
	bounds := [4]float64{xmin, ymin, xmax - xmin, ymax - ymin}
	//fill in missing headers
	for len(h) < len(spairs) {
		nh := fmt.Sprintf("series # %v", len(h)+1)
		h = append(h, nh)
	}
	//draw
	b.startSVG(f.title, bounds, f.xl, f.yl)
	for n, s := range spairs {
		b.startPoly(h[n]) //header #n mapped to spair #n
		for i := range s[0] {
			x := s[0][i]
			y := s[1][i]
			b.vertex(x, y)
		}
		b.endPoly()
	}
	b.endSVG()
	b.writeLabels()
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

	test := "-grid=2\njane bob\n10 10\n20 10\n30 40\n\nsammy ethyl\n0 0\n 1 2\n0 0\n-0.31 -0.5"
	parseAllAndRender(strings.NewReader(test))
	// pa	rseAllAndRender(os.Stdin)
}
