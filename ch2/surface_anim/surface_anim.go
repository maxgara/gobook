package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
)

const (
	width, height = 600, 500
	cells         = 80
	xyrange       = 30.0
	xyscale       = width / 2 / xyrange
	zscale        = height * 0.4
	angle         = math.Pi / 6 // angle of x,y axes
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle)
var t int

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))

}

// return canvas point x, y at corner of cell i,j. last return val indicates status (ok/err).
func corner(i, j int) (float64, float64, string, int) {
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)
	z := f2(x, y)
	if math.IsNaN(z) {
		return 0, 0, "", 1
	}
	//calculate color
	rcomp := uint8(z * 255 / xyrange * 3)
	bcomp := uint8(255 - rcomp)
	str := fmt.Sprintf("%02x00%02x", rcomp, bcomp) //SVG stroke attr

	//project x,y,z onto 2-D canvas surface
	sx := width/2 + (x-y)*cos30*xyscale
	sy := width/2 + (x+y)*sin30*xyscale + -z*zscale
	return sx, sy, str, 0

}

func f(x, y float64) float64 {
	r := math.Hypot(x, y) // distance from (0,0)
	return math.Sin(r) / r
}
func f2(x, y float64) float64 {
	r := math.Hypot(x, y) // distance from (0,0)
	return math.Sin(r+float64(t)/10) / r
}
func f1(x, y float64) float64 {
	return x * y / 1000
}

/*
generate an SVG image and write it to out. if ref==true then the svg element tags at end+beginning will be omitted,
and only polygon tags will be returned
*/
func svggen(out io.Writer) {
	fmt.Fprintf(out, `
		<svg xmlns='http://www.w3.org/2000/svg' `+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)
	var maxx float64
	var maxy float64
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay, astr, aerr := corner(i+1, j)
			bx, by, _, berr := corner(i, j)
			cx, cy, _, cerr := corner(i, j+1)
			dx, dy, _, derr := corner(i+1, j+1)
			//skip error polygons and continue
			if aerr|berr|cerr|derr > 0 {
				continue
			}

			fmt.Fprintf(out, "<polygon points='%.5g,%.5g %.5g,%.5g %.5g,%.5g %.5g,%.5g' stroke='#%s'/>\n",
				ax, ay, bx, by, cx, cy, dx, dy, astr)
			var xs = []float64{ax, bx, cx, dx}
			var ys = []float64{ay, by, cy, dy}
			for k := range xs {
				if xs[k] > maxx {
					maxx = xs[k]
				}
				if ys[k] > maxy {
					maxy = ys[k]
				}
			}
		}

	}
	fmt.Fprintf(out, "</svg>")
	_, _ = maxx, maxy
	// fmt.Printf("max x:%f, maxy:%f", maxx, maxy)
}

func handler(w http.ResponseWriter, r *http.Request) {
	t++
	fmt.Printf("t=%v\n", t)

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Refresh", "0.1")
	svggen(w)
}
