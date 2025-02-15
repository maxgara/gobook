package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	width, height = 800, 800 //window dims
)

var file *os.File

func main() {
	// file, _ = os.Create("cpuprofile")
	// pprof.StartCPUProfile(file)
	// defer pprof.StopCPUProfile()

	// file, _ = os.Create("memprofile")
	setupAndDraw()
}

func loadobjfile(filename string) {
	os.ReadFile("african_head.obj")

}

// bgra
func draw(pixels []byte) {
	//simpler line drawing func for convenience
	line := func(x0, y0, x1, y1 int) {
		DrawLine(x0, y0, x1, y1, pixels)
	}
	line(width/3, height, width, height/3)
	DrawLine(0, 0, width, height, pixels)
	DrawLine(0, 0, width/2, height, pixels)
	DrawLine(0, 0, width, height/2, pixels)
	DrawLine(0, height, width, 0, pixels)
	// putpixel(width/2, height/2, 0xffffff00, pixels)
	// line(height*1/6, 0xff000000, pixels)
	// line(height*2/6, 0x00ff0000, pixels)
	// line(height*3/6, 0x0000ff00, pixels)
	// line(height*4/6, 0x000000ff, pixels)

}
func DrawLine(x0, y0, x1, y1 int, pixels []byte) {
	var color uint32 = 0x0000ff00
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
	if x >= width || y >= height {
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
		<-time.After(time.Second * 60)

		dur = time.Since(start)
		done = true
	}()

	var loops uint64
	// squares := make([][]sdl.Point, 0, cells*cells)
	var first bool
	first = true
	//draw loop
	for {
		loops++
		rect := sdl.Rect{0, 0, height, width}
		blanksurf.Blit(&rect, surf, &rect)
		surf.Lock()
		pix := surf.Pixels()
		// drawLines(arr, pix)
		draw(pix)
		surf.Unlock()
		first = !first
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
		sdl.Delay(100)
	}

}
func drawLines(arr []sdl.Point, pix []byte) {
	// fmt.Println(arr)

	setp := func(p sdl.Point) {
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
			setp(p)
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
