package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	width, height = 800, 800 //window dims
)

var file *os.File
var fileVerts []vert
var fileFaces []face
var filebounds struct {
	xmin, xmax, ymin, ymax float64
}

func main() {
	vec1 := [3]float64{1, 0, 0}
	vec2 := [3]float64{0, 1, 0}
	fmt.Printf("%v X %v = %v\n", vec1, vec2, cross(vec1, vec2))
	// b := bary(v, w, [2]float64{0.15, 0.1})

	// fmt.Printf("b=<%f %f>\n\n\n", b[0], b[1])
	// fmt.Println(b)
	// <-time.After(1 * time.Hour)
	// file, _ = os.Create("cpuprofile")
	// pprof.StartCPUProfile(file)
	// defer pprof.StopCPUProfile()

	// file, _ = os.Create("memprofile")
	fileVerts, fileFaces = loadobjfile("african_head.obj")
	// fmt.Println(fileVerts)
	// fmt.Println("faces")
	// fmt.Println(fileFaces)
	// fmt.Println("bounds")
	// fmt.Println(filebounds)
	setupAndDraw()
}

type vert [3]float64
type face = [3]int

// load vertex from string
func loadVertex(s string, verts *[]vert) {
	fs := strings.Fields(s)
	var vt vert
	for i, coordstr := range fs[1:] {
		coord, err := strconv.ParseFloat(coordstr, 64)
		if err != nil {
			log.Fatal(err)
		}
		vt[i] = coord
	}
	//adjust bounds
	b := &filebounds
	b.xmax = max(b.xmax, vt[0])
	b.xmin = min(b.xmin, vt[0])
	b.ymax = max(b.ymax, vt[1])
	b.ymin = min(b.ymin, vt[1])
	*verts = append(*verts, vt)
}
func loadface(s string, faces *[]face) {
	fs := strings.Fields(s)
	var f face
	for i, field := range fs[1:] {
		idxstr := strings.Split(field, "/")[0]
		fidx, err := strconv.ParseInt(idxstr, 10, 32)
		if err != nil {
			log.Fatal(err)
		}
		f[i] = int(fidx)
	}
	*faces = append(*faces, f)
}
func loadobjfile(filename string) (verts []vert, faces []face) {
	f, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	s := strings.Split(string(f), "\n")

	for _, v := range s {
		switch {
		case strings.HasPrefix(v, "v "):
			loadVertex(v, &verts)
		case strings.HasPrefix(v, "f "):
			loadface(v, &faces)
		}
	}
	// fmt.Println(verts)
	return verts, faces
}

// convert vertices to pixels
func vtop(v vert) (p [3]int) {
	for i := range 3 {
		vadj := ((v[i] + 1) / 2) //adjust so v >= 0
		scale := min(width, height)
		p[i] = int(float64(scale) * vadj)
	}
	return p
}
func vtop2(v vert) (p [2]int) {
	for i := range 2 {
		vadj := ((v[i] + 1) / 2) //adjust so v >= 0
		scale := min(width, height)
		p[i] = int(float64(scale) * vadj)
	}
	return p
}

// rotate around +y axis by t(heta) radians
func rotateVert(v vert, t float64) vert {
	st := math.Sin(t)
	ct := math.Cos(t)
	x := v[0]
	y := v[1]
	z := v[2]
	x1 := x*ct + z*st
	y1 := y
	z1 := z*ct - x*st
	return vert{x1, y1, z1}
}
func update() {
	for i := range fileVerts {
		v := fileVerts[i]
		v = rotateVert(v, 0.01)
		fileVerts[i] = v
	}
}

// bgra
func draw(pixels []byte) {
	//simpler line drawing func for convenience
	// line := func(x0, y0, x1, y1 int) {
	// 	DrawLine(x0, y0, x1, y1, pixels)
	// }

	// v0 := [2]int{0, 0}
	// v1 := [2]int{width / 2, height / 3}
	// v2 := [2]int{0, height}
	// fillTriangle(v0, v1, v2, pixels)

	vline := func(vertex0, vertex1 vert) {
		p0 := vtop(vertex0)
		p1 := vtop(vertex1)
		DrawLine(p0[0], p0[1], p1[0], p1[1], pixels)
	}
	fill3 := func(vt0, vt1, vt2 vert) {
		p0, p1, p2 := vtop2(vt0), vtop2(vt1), vtop2(vt2)
		fillTriangle(p0, p1, p2, pixels)
	}

	// u := vert{0, 0, 0}
	// v := vert{0.2, 0, 0}
	// w := vert{0.1, 0.1, 0}
	// U, V, W := vtop2(u), vtop2(v), vtop2(w)
	// fillTriangle(U, V, W, pixels)

	for _, f := range fileFaces {
		vidx0, vidx1, vidx2 := f[0], f[1], f[2]
		v0, v1, v2 := fileVerts[vidx0-1], fileVerts[vidx1-1], fileVerts[vidx2-1]
		// vline(v0, v1)
		// vline(v1, v2)
		// vline(v2, v0)
		fill3(v0, v1, v2)
	}

	for _, f := range fileFaces {
		vidx0, vidx1, vidx2 := f[0], f[1], f[2]
		v0, v1, v2 := fileVerts[vidx0-1], fileVerts[vidx1-1], fileVerts[vidx2-1]
		vline(v0, v1)
		vline(v1, v2)
		vline(v2, v0)
		// fill3(v0, v1, v2)
	}
}
func DrawLine(x0, y0, x1, y1 int, pixels []byte) {
	// var color uint32 = 0x0000ff00
	var color uint32 = 0xff000000
	//x_i=x0 + i*(x1-x0)/N
	if x1 < x0 {
		x0, x1 = x1, x0
		y0, y1 = y1, y0
	}
	var yflip bool
	if y1 < y0 {
		yflip = true
	}
	Dx := x1 - x0
	Dy := y1 - y0
	if yflip {
		Dy = y0 - y1
	}

	N := max(Dx, Dy)

	x, y := x0, y0
	var nxerr, nyerr int
	for i := 0; i < N; i++ {
		nxerr += Dx
		nyerr += Dy
		if nxerr > N {
			nxerr -= N
			x++
		}
		if nyerr > N && yflip {
			nyerr -= N
			y--
		}
		if nyerr > N && !yflip {
			nyerr -= N
			y++
		}
		putpixel(x, y, color, pixels)
	}
}
func putpixel(x, y int, color uint32, pixels []byte) {
	if x >= width || y >= height || x < 0 || y < 0 {
		return
	}
	idx := 4 * (x + width*y)
	a := byte(color & 0xff000000 >> 24)
	r := byte(color & 0x00ff0000 >> 16)
	g := byte(color & 0x0000ff00 >> 8)
	b := byte(color & 0x000000ff)
	pixels[idx] = a
	pixels[idx+1] = r
	pixels[idx+2] = g
	pixels[idx+3] = b
}

func setupAndDraw() {

	var window *sdl.Window
	var renderer *sdl.Renderer
	var winTitle string = "TinyRenderer"
	window, err := sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(width), int32(height), sdl.WINDOW_SHOWN|sdl.WINDOW_ALWAYS_ON_TOP)
	if err != nil {
		log.Fatal(err)
	}
	surf, err := window.GetSurface()
	if err != nil {
		log.Fatal(err)
	}
	//blank surface to blit before redrawing
	blanksurf, err := sdl.CreateRGBSurface(5, 800, 800, 32, 0, 0, 0, 0)
	if err != nil {
		log.Fatal(err)
	}

	defer window.Destroy()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		os.Exit(1)
	}
	//benchmarking
	var dur time.Duration
	start := time.Now()
	var done bool
	var loopsPerSec float64
	go func() {
		<-time.After(time.Second * 120)

		dur = time.Since(start)
		done = true
	}()

	var loops uint64
	// squares := make([][]sdl.Point, 0, cells*cells)
	//draw loop
	for {
		loops++
		rect := sdl.Rect{0, 0, height, width}
		blanksurf.Blit(&rect, surf, &rect)
		surf.Lock()
		pix := surf.Pixels()
		// drawLines(arr, pix)
		update()
		draw(pix)
		surf.Unlock()
		window.UpdateSurface()
		renderer.Present()

		if done {
			loopsPerSec = float64(loops) / float64(dur.Seconds())
			surf.Lock()
			// fmt.Println(surf.Pixels())
			runtime.GC() // get up-to-date statistics
			pprof.WriteHeapProfile(file)
			fmt.Printf("dur=%v; loops=%v; lps=%v\n", dur.Seconds(), loops, loopsPerSec)
			return
		}
		// for event := sdl.PollEvent(); event == nil; {
		// 	event = sdl.PollEvent()
		// }
		if event := sdl.PollEvent(); event != nil {
			if event, ok := event.(*sdl.KeyboardEvent); ok {
				fmt.Printf("event.Keysym.Scancode: %v %[1]c\n", event.Keysym.Scancode)
				switch event.Keysym.Scancode {
				case 20: //q
					dur = time.Since(start)
					done = true
				case 82: //up
				case 81: //down
				case 80: //left
				case 79: //right
				}

			}
			if _, ok := event.(*sdl.QuitEvent); ok {
				return
			}
		}
		sdl.Delay(25)

	}
}

// fmt.Println(arr)
// }
// 	setp := func(p sdl.Point) {
// 		x := p.X
// 		y := p.Y
// 		if x < 0 || x >= width || y < 0 || y >= width {
// 			return
// 		}
// 		idx := x*4 + y*width*4
// 		pix[idx] = 255
// 	}
// var N = 100
// drawLine := func(x0, y0, x1, y1 int32) {
// 	for i := range N {
// 		perc := (100 * i) / (N - 1)
// 		// fmt.Println("percent=%v\n", perc)
// 		x := (x0*(100-int32(perc)) + x1*int32(perc)) / 100
// 		y := (y0*(100-int32(perc)) + y1*int32(perc)) / 100
// 		p := sdl.Point{x, y}
// 		setp(p)
// 	}
// }
// drawLineSDLPoint := func(p0, p1 sdl.Point) {
// 	drawLine(p0.X, p0.Y, p1.X, p1.Y)
// }
// for i := 0; i < len(arr); i += 4 {
// 	a, b, c, d := arr[i], arr[i+1], arr[i+2], arr[i+3]
// 	drawLineSDLPoint(a, b)
// 	drawLineSDLPoint(b, c)
// 	drawLineSDLPoint(c, d)
// 	drawLineSDLPoint(d, a)
// }
// // for i := 0; i < width; i++ {
// // 	drawPoint(sdl.Point{X: int32(i), Y: height / 2})

// // }
// drawLine(0, height/2, width-1, height/2)
// }
type Box struct {
	xmin int
	ymin int
	xmax int
	ymax int
}

func bound(points [][2]int) Box {
	xmax, ymax := -1000, -1000
	xmin, ymin := 1000, 1000
	for _, p := range points {
		xmin, ymin = min(xmin, p[0]), min(ymin, p[1])
		xmax, ymax = max(xmax, p[0]), max(ymax, p[1])
	}
	return Box{xmin: xmin, ymin: ymin, xmax: xmax, ymax: ymax}
}

// fill in a triangle between pixels p{0,1,2}
func fillTriangle(p0, p1, p2 [2]int, pix []byte) {
	u := vdiff(p1, p0)
	v := vdiff(p2, p0)
	box := bound([][2]int{p0, p1, p2})
	for i := box.xmin; i < box.xmax; i++ {
		for j := box.ymin; j < box.ymax; j++ {
			// fmt.Printf("looping: i=%v j=%v box=%v\n", i, j, box)
			//offset X to be a vector w/r/t p0
			X := [2]float64{float64(i - p0[0]), float64(j - p0[1])}
			b := bary(u, v, X)
			// fmt.Printf("\tbary= %v %v", b[0], b[1])
			// mag := math.Hypot(b[0], b[1])
			mag := b[0] + b[1]
			//check if b is in triangle made up of vectors u,v
			if b[0] < 0 || b[1] < 0 || mag > 1 {
				continue
			}
			// fmt.Printf("pass check, in triangle")
			putpixel(i, j, 0x0000ff00, pix)
		}
	}
	// fmt.Println("done")
}

// vector cross product
func cross(a, b [3]float64) [3]float64 {
	a1, a2, a3 := a[0], a[1], a[2]
	b1, b2, b3 := b[0], b[1], b[2]
	x := a2*b3 - a3*b2
	y := a3*b1 - a1*b3
	z := a1*b2 - a2*b1
	return [3]float64{x, y, z}
}

// vector subtract u - w
func vdiff(u, w [2]int) [2]float64 {
	x := float64(u[0] - w[0])
	y := float64(u[1] - w[1])
	return [2]float64{x, y}
}

// return representation of X in terms of basis vectors v,w
func bary(v, w, X [2]float64) [2]float64 {
	a := v[0]
	b := w[0]
	c := v[1]
	d := w[1]
	x := X[0]
	y := X[1]
	//make sure determinant is ! = 0
	if a*d-b*c == 0 {
		return [2]float64{-1, -1}
	}
	det := 1 / (a*d - b*c)
	//calculate coefficients for vectors v,w
	cv := d*det*x - b*det*y
	cw := -c*det*x + a*det*y
	return [2]float64{cv, cw}
}

/*
generate an SVG image and write it to out. if ref==true then the svg element tags at end+beginning will be omitted,
and only polygon tags will be returned
*/
// func svggen(out io.Writer) {
// 	fmt.Fprintf(out, `
// 		<svg xmlns='http://www.w3.org/2000/svg' `+
// 		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
// 		"width='%d' height='%d'>", width, height)
// 	var maxx float64
// 	var maxy float64
// 	for i := 0; i < cells; i++ {
// 		for j := 0; j < cells; j++ {
// 			ax, ay, astr, aerr := corner(i+1, j)
// 			bx, by, _, berr := corner(i, j)
// 			cx, cy, _, cerr := corner(i, j+1)
// 			dx, dy, _, derr := corner(i+1, j+1)
// 			//skip error polygons and continue
// 			if aerr|berr|cerr|derr > 0 {
// 				continue
// 			}

// 			fmt.Fprintf(out, "<polygon points='%.5g,%.5g %.5g,%.5g %.5g,%.5g %.5g,%.5g' stroke='#%s'/>\n",
// 				ax, ay, bx, by, cx, cy, dx, dy, astr)
// 			var xs = []float64{ax, bx, cx, dx}
// 			var ys = []float64{ay, by, cy, dy}
// 			for k := range xs {
// 				if xs[k] > maxx {
// 					maxx = xs[k]
// 				}
// 				if ys[k] > maxy {
// 					maxy = ys[k]
// 				}
// 			}
// 		}

// 	}
// 	fmt.Fprintf(out, "</svg>")
// 	_, _ = maxx, maxy
// 	// fmt.Printf("max x:%f, maxy:%f", maxx, maxy)
// }
