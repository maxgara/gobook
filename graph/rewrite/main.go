// Create SVG graphs of data.
// Usage: data is read from stdin, with data points separated by newlines. Output is to stdout and is comprised of
// HTML with embedded SVG graph(s).
// data points can be one coordinate or two, separated by a whitespace character.
// multiple data series can be displayed by separating the series with non-numerical text or an empty line.
// additional empty lines between plots are ignored.
// The following flags are supported between data series:
//
// -n 	New Chart		 By default curves are all plotted on the same chart. This starts a new chart
//
//	which is used for further data series.
//
// -p 	Parallel Plot 	Plot the next curve next to the previous ones on the same chart
// -css [<css properties>]
//
// Each flag must be on its own line with nothing other than whitespace. -n overrides -p if both are set.
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
	"os"
)

type svg struct {
	sset                   [][]float64 //data series
	label                  string
	Xmin, Xmax, Ymin, Ymax float64 //bounds for SVG viewbox
}

type grid struct {
	rows    []svg
	columns []svg
}

type parser struct {
	done bool
}

func (p *parser) parseLine(line string) {

}

func (g *grid) String() string

// print web page
func print()

func main() {
	// var p parser
	r := bufio.NewScanner(os.Stdin)
	for {
		ok := r.Scan()

		if !ok {
			break
		}
	}
}
