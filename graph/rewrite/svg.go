package main

import (
	"fmt"
	"io"
)

var coolpal = []uint32{0xf5f5b0, 0x154734, 0xcfa3a8, 0x6a7ba2, 0xd6d6d6} // cool-tone color palette

// SVG object methods facilitate drawing geometric objects in SVG format
type SVGPoly struct {
	fill   uint32
	stroke uint32
	width  float64
	points [][2]float64
}
type SVG struct {
	fill   uint32
	stroke uint32
	polys  []SVGPoly
	bounds [4]float64
	wp     SVGPoly //working poly (defined after polystart, before polyend)
}

func newSVG() *SVG {
	return &SVG{}
}
func setStrokeRGB(g *SVG, c uint32) {
	g.wp.stroke = c
}
func setFillRGP(g *SVG, c uint32) {
	g.wp.fill = c
}

func (g *SVG) polyStart() {
	cidx := len(g.polys)
	if g.stroke == 0 {
		g.stroke = coolpal[cidx] //default color scheme, may remove this later.
	}
	g.wp = SVGPoly{fill: g.fill, stroke: g.stroke, width: POLY_STROKE_WIDTH_DEFAULT}
}
func (g *SVG) polyEnd() {
	g.polys = append(g.polys, g.wp)
}
func (g *SVG) vertex(x, y float64) {
	g.wp.points = append(g.wp.points, [2]float64{x, y})
}
func (g *SVG) vertexMulti(pts [][2]float64) {
	g.wp.points = append(g.wp.points, pts...)
}

// generate svg markdown
func (g *SVG) render(w io.Writer) {
	//get bounds
	var xmin, ymin, xmax, ymax float64
	fp := g.polys[0].points[0] //first pt
	xmin = fp[0]
	xmax = fp[0]
	ymin = fp[0]
	ymax = fp[0]
	for _, poly := range g.polys {
		for _, p := range poly.points {
			x := p[0]
			y := p[1]
			xmax = max(xmax, x)
			xmin = min(xmin, x)
			ymax = max(ymax, y)
			ymin = min(ymin, y)
		}
	}
	bounds := []float64{xmin, ymin, xmax - xmin, ymax - ymin}
	s := fmt.Sprintf(svgstart, bounds[0], bounds[1], bounds[2], bounds[3])
	for _, poly := range g.polys {
		s += polymd(poly)
	}
	w.Write([]byte(s))
}

// generate polyline markdown
func polymd(ply SVGPoly) string {
	s := fmt.Sprintf(polystart, ply.stroke, ply.fill, ply.width)
	for _, p := range ply.points {
		s += fmt.Sprintf("%v,%v ", p[0], p[1])
	}
	return s + polyend
}

const (
	POLY_STROKE_WIDTH_DEFAULT = 0.5
	svgstart                  = `<svg viewBox="%v %v %v %v" preserveAspectRatio="none"
			style="width:94%%; height: 94%%; padding: 3%%; background: grey; border: coral solid"
			xmlns="http://www.w3.org/2000/svg">`
	svgend    = `</svg>` //FIX THIS
	polystart = `<polyline stroke="#%d" fill="%d" stroke-width="%v"
				points="`
	polyend = `">
			</polyline>`
)
