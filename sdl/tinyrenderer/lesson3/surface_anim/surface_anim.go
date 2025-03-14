package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	width, height = 800, 800        //window dims
	cells         = 160             //resolution
	xyrange       = 15.0            //XMAX-XMIN, YMAX-YMIN, for function (assumed they are equal)
	xyscale       = width / xyrange //convert abstract coordinates to screen dims.
	// angle = 0
)

var angle = math.Pi / 6 // angle of x,y axes (30 degrees at pi/6) (+z axis rotation)
// var tilt = math.Pi * 2 / 5 //second rotation, clockwise +x axis
var tilt float64 = math.Pi / 6

var t int
var file *os.File

func main() {
	// file, _ = os.Create("cpuprofile")
	// pprof.StartCPUProfile(file)
	// defer pprof.StopCPUProfile()

	// file, _ = os.Create("memprofile")
	draw()
}

// return canvas point x, y at corner of cell i,j. last return val indicates status (ok/err).
func corner(i, j int) (float64, float64, int) {
	var sinAng, cosAng = math.Sin(angle), math.Cos(angle)
	var sinTilt, cosTilt = math.Sin(tilt), math.Cos(tilt)
	//get x,y,z coords for idxs
	x := xyrange * (float64(i)/(cells-1) - 0.5)
	y := xyrange * (float64(j)/(cells-1) - 0.5)
	z := f(x, y)
	// fmt.Printf("corner (i=%v,j=%v) has x=%v, y=%v, z=%v; -> ", i, j, x, y, z)
	if math.IsNaN(z) {
		return 0, 0, 1
	}
	//calculate color
	// come back to this
	// rcomp := uint8(z * 255 / xyrange * 3)
	// bcomp := uint8(255 - rcomp)

	//use right handed coordinate system so:
	/*
			y
			│
		    │
		    │
		    │--------> x
			⊙
			z (coming out of the screen)
	*/

	//now rotate counter-clockwise around +z-axis:

	/*
		  ← y
			│
		    │
		    │		   ↑
		    │--------> x
			⊙ ↺
			z
	*/
	x = x*cosAng - y*sinAng
	y = x*sinAng + y*cosAng

	// now, tilt y axis back (rotate clockwise around +x axis)
	/*
			↑z
		    │
		    ↑y
		    │
		    │--------> x
	*/
	y = y*cosTilt + z*sinTilt

	//project x,y,z onto screen
	x = width/2 + (x * xyscale)
	y = height/2 + (y * xyscale)
	// fmt.Printf("pixels x=%v y=%v\n", x, y)
	return x, y, 0
}
func f0(x, y float64) float64 {
	return 0
}

// change theta with time so that mesh moves
var theta float64

func f(x, y float64) float64 {
	r := math.Hypot(x, y) // distance from (0,0)
	return 3 * math.Sin(r+theta) / r
}
func f2(x, y float64) float64 {
	// r := math.Hypot(x, y) // distance from (0,0)
	// return math.Sin(r+float64(t)/10) / r
	return math.Hypot(x, y)
}
func f1(x, y float64) float64 {
	return x * y / 1000
}
func getSquares(arr *[]sdl.Point) {

	*arr = make([]sdl.Point, 0, len(*arr))
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay, aerr := corner(i+1, j)
			bx, by, berr := corner(i, j)
			cx, cy, cerr := corner(i, j+1)
			dx, dy, derr := corner(i+1, j+1)
			//skip error polygons and continue
			if aerr|berr|cerr|derr > 0 {
				continue
			}
			a := sdl.Point{int32(ax), int32(ay)}
			b := sdl.Point{int32(bx), int32(by)}
			c := sdl.Point{int32(cx), int32(cy)}
			d := sdl.Point{int32(dx), int32(dy)}
			*arr = append(*arr, a, b, c, d)
		}
	}
	// fmt.Println(squares)
	// fmt.Printf("max x:%f, maxy:%f", maxx, maxy)
	// return squares
	fmt.Printf("len arr = %v\n", len(*arr))
}
func draw() {

	var window *sdl.Window
	var renderer *sdl.Renderer
	// var points []sdl.Point
	// var rect sdl.Rect
	// var rects []sdl.Rect
	var winTitle string = "Go-SDL2 Render"
	window, err := sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(width), int32(height), sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatal(err)
	}
	surf, err := window.GetSurface()
	// blanksurf, err := sdl.CreateRGBSurfaceFrom(unsafe.Pointer(&surf.Pixels()), width, height, 32, 0, 0, 0, 0, 0)
	// fmt.Println(sdl.GetPixelFormatName(uint(surf.Format.Format)))
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(surf.Pixels())

	blanksurf, err := sdl.CreateRGBSurface(5, 800, 800, 32, 0, 0, 0, 0)
	if err != nil {
		log.Fatal(err)
	}
	//5 0x6000002be880 800 800 3200 0x158028000 <nil> 0 <nil> {0 0 800 800} 0x6000039b8f80 1
	// fmt.Printf("surf:%v width:%v, height:%v\n\n", surf, surf.W, surf.H)

	defer window.Destroy()

	// renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		os.Exit(1)
	}
	// sdl.Surf
	// defer renderer.Destroy()
	// fmt.Printf("existing window surface:", surf)

	var dur time.Duration
	start := time.Now()
	var done bool
	var loopsPerSec float64
	go func() {
		<-time.After(time.Second * 60)

		dur = time.Since(start)
		done = true
	}()

	var loops uint64
	// squares := make([][]sdl.Point, 0, cells*cells)
	arr := make([]sdl.Point, 0, cells*cells)
	var first bool
	first = true
	for {
		// fmt.Println(arr)
		loops++
		getSquares(&arr)

		// for i3 := range pix[:len(pix)/4] {
		// 	i := i3 * 4
		// 	pix[i] = 255
		// }
		// if first {
		rect := sdl.Rect{0, 0, height, width}
		blanksurf.Blit(&rect, surf, &rect)
		surf.Lock()
		pix := surf.Pixels()
		drawLines(arr, pix)
		surf.Unlock()
		first = !first
		// }
		window.UpdateSurface()

		// renderer.SetDrawColor(0, 0, 0, 255)
		// err := renderer.Clear()
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// renderer.SetDrawColor(255, 255, 255, 255)
		// err = renderer.DrawLines(arr)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		renderer.Present()

		// <-time.After(5000 * time.Millisecond)

		if done {
			surf.Lock()
			// fmt.Println(surf.Pixels())
			runtime.GC() // get up-to-date statistics
			pprof.WriteHeapProfile(file)
			loopsPerSec = float64(loops) / float64(dur.Seconds())
			fmt.Printf("dur=%v; loops=%v; lps=%v\n", dur.Seconds(), loops, loopsPerSec)
			return
		}
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
		// sdl.Delay(50)
		theta += 0.1
		// tilt += (1.0 / 360.0)
		angle += 0.01
	}

}
func drawLines(arr []sdl.Point, pix []byte) {
	// fmt.Println(arr)

	drawPoint := func(p sdl.Point) {
		x := p.X
		y := p.Y
		if x < 0 || x >= width || y < 0 || y >= width {
			return
		}
		idx := x*4 + y*width*4
		pix[idx] = 255
	}
	var N = 100
	drawLine := func(x0, y0, x1, y1 int32) {
		for i := range N {
			perc := (100 * i) / (N - 1)
			// fmt.Println("percent=%v\n", perc)
			x := (x0*(100-int32(perc)) + x1*int32(perc)) / 100
			y := (y0*(100-int32(perc)) + y1*int32(perc)) / 100
			p := sdl.Point{x, y}
			drawPoint(p)
		}
	}
	drawLineSDLPoint := func(p0, p1 sdl.Point) {
		drawLine(p0.X, p0.Y, p1.X, p1.Y)
	}
	for i := 0; i < len(arr); i += 4 {
		a, b, c, d := arr[i], arr[i+1], arr[i+2], arr[i+3]
		drawLineSDLPoint(a, b)
		drawLineSDLPoint(b, c)
		drawLineSDLPoint(c, d)
		drawLineSDLPoint(d, a)
	}
	// for i := 0; i < width; i++ {
	// 	drawPoint(sdl.Point{X: int32(i), Y: height / 2})

	// }
	drawLine(0, height/2, width-1, height/2)
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
