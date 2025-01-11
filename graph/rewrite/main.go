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
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// type svg struct {
// 	sset                   [][]float64 //data series
// 	Xmin, Xmax, Ymin, Ymax float64     //bounds for SVG viewbox
// }

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
}

func newParser(r io.Reader) parser {
	s, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return parser{s: string(s)}
}

// parses a section of input into native golang types (slices, strings, etc.)
// each call to parse reads until a break in the data series
// sets flags, titles, and text section properties for p
// returns false when parsing is complete, either due to error or end of input

func (p *parser) parse() bool {
	lines := strings.Split(p.s, "\n")
	var lidx int
	var textdone, datadone bool
	var data [][]float64
	var flagstrs []string
	//loop over blocks
	for {
		l := strings.Trim(lines[lidx], " \t\r")
		words := strings.Fields(l)
		//loop over non-data
		if !textdone {
			if isdata(l) {
				textdone = true
				data = make([][]float64, len(words))
				continue
			}
			if len(l) == 0 {
				lidx++
				continue
			}
			if l[0] == '-' {
				flagstrs = append(flagstrs, words...)
				lidx++
			}
			//handle text paragraphs
			p.text += l
			lidx++
			continue
		}
		if !datadone {
			// handle data
			for i, w := range words {
				d, err := strconv.ParseFloat(w, 64)
				if err != nil {
					datadone = true
					break
				}
				data[i] = append(data[i], d)
			}
			continue
		}
		flags := flag.NewFlagSet("svgparseflags", flag.PanicOnError)
		var fbits uint32
		d1Flag := flags.Bool("d1", true, "one coordinate series parsing")
		d2Flag := flags.Bool("d2", false, "two coordinate series parsing")
		dxFlag := flags.Bool("dx", false, "shared x value")
		gridFlag := flags.Int("grid", 0, "grid column count")

		flags.Parse(flagstrs)
		//handle result of parsing this block

		encode(p.text, p.data, *d1Flag, *d2Flag, *dxFlag, *gridFlag, newDocBuilder(os.Stdout))
	}
}

const (
	D1_FLAG = 0b1 << iota
	D2_FLAG
	DX_FLAG
)

// convert parsed data to document format
func encode(text string, data [][]float64, d1Flag, d2Flag, dxFlag bool, gridFlag int, b *docBuilder) {
	b.writeText(text)
	//assign pairs of data columns to contain x,y values of series based on flags
	var spairs [][2][]float64
	idxs := make([]float64, len(data)) //new data series for convenience
	for i := range len(data) {
		idxs[i] = float64(i)
	}
	switch {
	case d2Flag && !dxFlag:
		for i := 0; i < len(data); i += 2 {
			spair := [2][]float64{data[i], data[i+1]}
			spairs = append(spairs, spair)
		}
	case d1Flag:
		for _, s := range data {
			spair := [2][]float64{idxs, s}
			spairs = append(spairs, spair)
		}
	case dxFlag:
		for i := 1; i < len(data); i++ {
			spair := [2][]float64{data[0], data[i]}
			spairs = append(spairs, spair)
		}
	}

	//get bounds
	// var xmin = math.MaxFloat64
	// var ymin = math.MaxFloat64
	// var xmax = -math.MaxFloat64
	// var ymax = -math.MaxFloat64
	// b.startSVG("placeholder", 0, 0, 0, 0)

}

func isdata(s string) bool {
	datachars := []byte("0123456789.E")
	for _, c := range []byte(s) {
		ok := true
		for _, dc := range datachars {
			if c == dc {
				ok = true
				break
			}
		}
		if !ok {
			return false
		}
	}
	return true
}
func (p *parser) old_parse() bool {
	*p = parser{s: p.s, readback: p.readback, flags: p.flags, title: p.title, pagetitle: p.pagetitle, css: p.css} //reset p
	for {
		var l []byte
		if p.readback != nil {
			l = p.readback
			p.readback = nil
		} else {
			if ok := p.s.Scan(); !ok {
				p.err = io.EOF
				return false
			}
			l = p.s.Bytes()
		}
		if len(l) == 0 {
			continue
		}
		words := lwords(l)
		if len(words) == 0 {
			continue
		}
		//handle data
		isdata := true //default
		_, err := strconv.ParseFloat(string(words[0]), 64)
		if err != nil {
			isdata = false
		}
		if isdata {
			p.t = DATA
			p.dcols = len(words)
			p.data, p.err = parsedstream(p.s, p.dcols)
			if p.err != nil {
				return false
			}
			p.readback = p.s.Bytes() // capture last line read by parsedstream, for re-processing
			return true
		}
		//handle text
		if words[0][0] != '-' {
			p.t = TEXT
			p.text = strings.Trim(string(l), "\t ")
			return true
		}
		//handle flags
		p.t = FLAGS
		ls := string(l)
		switch {
		case ls == "-n":
			p.flags |= NEWCHART
		case ls == "-p":
			p.flags |= PARALLEL
		case strings.HasPrefix(ls, "-css="):
			p.css = strings.TrimPrefix(ls, "-css=")
		case strings.HasPrefix(ls, "-pagetitle="):
			p.pagetitle = strings.TrimPrefix(ls, "-pagetitle=")
		case strings.HasPrefix(ls, "-title="):
			p.title = strings.TrimPrefix(ls, "-title=")
		}
		return true
	}
}

// parse data stream into slices
func parsedstream(s *bufio.Scanner, dcols int) ([][]float64, error) {
	first := true
	data := make([][]float64, dcols)
	ok := true
	for {
		//scan thru data, first line is pre-scanned by parse()
		if !first {
			if ok = s.Scan(); !ok {
				return data, io.EOF
			}
		}
		first = false
		l := s.Text()
		row := strings.Fields(l)
		for i, xstr := range row {
			x, err := strconv.ParseFloat(xstr, 64)
			if err != nil && i == 0 {
				return data, nil
			}
			if err != nil && i != 0 {
				return data, fmt.Errorf("parsedstream: non-data in data series")
			}
			data[i] = append(data[i], x)
		}
	}
}

// func (g *grid) String() string

// print web page
// func print()

// break line into words
func lwords(l []byte) [][]byte {
	var words [][]byte
	after := -1
	for i, b := range l {
		if b != ' ' && b != '\t' {
			continue
		}
		if i == after+1 {
			after = i //drop whitespace at beginning of word
		}
		words = append(words, l[after+1:i])
		after = i
	}
	if after+1 < len(l) {
		words = append(words, l[after+1:])
	}
	return words
}

func main() {
	// s := bufio.NewScanner(os.Stdin)
	sflag := flag.Bool("multiseries", false, "multi-series mode")
	flag.BoolFunc("s", "set multi-series mode", func(s string) error {
		flag.Set("multiseries", "1")
		return nil
	})
	flag.Parse()
	fmt.Printf("flag: %v\n", *sflag)
	// p := parser{s: s}
	// for {
	// 	ok := p.parse()
	// 	if p.t == DATA {
	// 	}
	// 	//handle parse results
	// 	if !ok {
	// 		break
	// 	}
	// }
}
