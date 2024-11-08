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
	const linestart = `<polyline stroke="grey" fill="none" stroke-width="0.7" points="`
	const lineend = `"> </polyline>`
	const svgstart = `<div style="width: 100%; height: 100%; display: flex; justify-content: center; align-items: center;"><svg width="100%" height="100%" viewBox="0 0 2000 2000" xmlns="http://www.w3.org/2000/svg" overflow="visible">`
	const svgend = `</svg></div>`
	var bounds = []float64{10000, 0, 10000, 0}
	var out string
	var outtemp string   //store string until it is ready to be added to out
	outtemp += linestart //add line start but do not add svg start until we have bounds
	//read input
	r := bufio.NewScanner(os.Stdin)
	var idx int
	for ok := r.Scan(); ok; ok = r.Scan() {
		// fmt.Printf(`scanned:"%s"`, r.Text())
		line := r.Text()
		t := linetype(line)
		//begin new line in same SVG
		if t == EMPTY {
			outtemp += lineend + linestart
			idx = 0
			continue
		}
		//begin new SVG and new line
		if t == NEWCHARTFLAG {
			//
			outtemp += lineend + svgend
			outtemp = printsvgstart(bounds) + outtemp // prepend svg start tag
			out += outtemp
			outtemp = linestart
			idx = 0
			continue
		}
		p := parsepoint(line, idx)
		updatebounds(bounds, p[0], p[1])
		outtemp += printpoint(p)
		idx++
	}
	outtemp += lineend + svgend
	outtemp = printsvgstart(bounds) + outtemp // prepend svg start tag
	out += outtemp
	fmt.Print(out)
}

func printpoint(p []float64) string {
	return fmt.Sprintf(" %d,%d", int(p[0]), int(p[1]))
}
func printsvgstart(b []float64) string {
	width := b[XMAX] - b[XMIN]
	height := b[YMAX] - b[YMIN]
	return fmt.Sprintf(`<div style="width: %s; height: %s; display: flex"><svg viewBox="%d %d %d %d" xmlns="http://www.w3.org/2000/svg">`, `100%`, `100%`, int(b[XMIN]), int(b[YMIN]), int(width), int(height))
}

const (
	DATA int = iota
	EMPTY
	NEWCHARTFLAG
)

func linetype(s string) int {
	if s == "" {
		return EMPTY
	}
	if s == "-n" {
		return NEWCHARTFLAG
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
