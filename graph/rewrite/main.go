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
	"io"
	"os"
	"strconv"
	"strings"
)

type svg struct {
	sset                   [][]float64 //data series
	Xmin, Xmax, Ymin, Ymax float64     //bounds for SVG viewbox
}

type grid struct {
	rows [][]svg
}

const (
	NEWCHART = 0b1 << iota
	PARALLEL
)

type parser struct {
	r         io.Reader //input to read from
	lidx      int       //line index
	isdata    bool
	isflag    bool
	dcols     int // # of data columns
	pagetitle string
	title     string //svg title
	text      string
	flags     int
	data      [][]float64
	err       error
	css       string
}

// parses a section of input into native golang types (slices, strings, etc.)
// returns false when parsing is complete, either due to error or end of input
func (p *parser) parse() bool {
	s := bufio.NewScanner(p.r)
	for {
		ok := s.Scan()
		if !ok {
			return false
		}
		l := s.Bytes()
		if len(l) == 0 {
			continue
		}
		words := lwords(l)
		if len(words) == 0 {
			continue
		}
		//handle data
		_, err := strconv.ParseFloat(string(words[0]), 64)
		if err != nil {
			p.isdata = false
		}
		if p.isdata {
			p.dcols = len(words)
			p.data, err = parsedstream(s, p.dcols)
			if err != nil {
				p.err = err
				return false
			}
		}
		if words[0][0] != '-' {
			p.text = strings.Trim(string(l), "\t ")
		}
		//handle flags
		for _, f := range words {
			fs := string(f)
			switch {
			case fs == "-n":
				p.flags |= NEWCHART
			case fs == "-p":
				p.flags |= PARALLEL
			case strings.HasPrefix(fs, "-css="):
				p.css = strings.TrimPrefix(fs, "-css")
			case strings.HasPrefix(fs, "-pagetitle="):
				p.pagetitle = strings.TrimPrefix(fs, "-pagetitle=")
			case strings.HasPrefix(fs, "-title="):
				p.title = strings.TrimPrefix(fs, "-title=")
			}
		}
	}
}

// parse data stream into slices
func parsedstream(s *bufio.Scanner, dcols int) ([][]float64, error) {
	first := true
	data := make([][]float64, dcols)
	for {
		ok := true
		//scan thru data, first line is pre-scanned by parse()
		if !first {
			if ok = s.Scan(); !ok {
				break
			}
		}
		l := s.Text()
		row := strings.Fields(l)
		for i, xstr := range row {
			x, err := strconv.ParseFloat(xstr, 64)
			if err != nil {
				return data, err
			}
			data[i] = append(data[i], x)
		}
	}
	return data, nil
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
	p := parser{r: os.Stdin}
	for {
		ok := p.parse()
		//handle parse results
		if !ok {
			break
		}
	}
}
