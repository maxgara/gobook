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
	filename      = "african_head.obj"
	// filename = "african_head.obj"
	delay   = 1
	yrotd   = 0.01 // +y azis rotation per frame
	xrotset = 0    // +x azis rotation
	RED     = 0x0000ff00
	GREEN   = 0x00ff0000
	BLUE    = 0xff000000
	ALPHA   = 0x000000ff
)

type F3 [3]float64

var window *sdl.Window

var wireframe bool
var file *os.File
var fileVerts []F3
var fileFaces [][3]int

// var fileFaceNorms []F3 // normal vector for each face (normalized to 1)
var done bool //control program exit
var loops uint64
var dur time.Duration
var greyval float64
var zbuff []float64
var lightpos []F3
var lightcolors []uint32
var lightpower []float64
var lightrot float64
var colorEnabled bool
var shadingEnabled bool

// load vertex from string
func loadVertex(s string, verts *[]F3) {
	fs := strings.Fields(s)
	var vt F3
	for i, coordstr := range fs[1:] {
		coord, err := strconv.ParseFloat(coordstr, 64)
		if err != nil {
			log.Fatal(err)
		}
		vt[i] = coord
	}
	*verts = append(*verts, vt)
}
func loadface(s string, faces *[][3]int) {
	fs := strings.Fields(s)
	var f [3]int
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
func loadobjfile(filename string) {
	f, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	s := strings.Split(string(f), "\n")

	for _, v := range s {
		switch {
		case strings.HasPrefix(v, "v "):
			loadVertex(v, &fileVerts)
		case strings.HasPrefix(v, "f "):
			loadface(v, &fileFaces)
		}
	}
}

// rotate around +y axis by t(heta) radians
func yrot(v F3, t float64) F3 {
	st := math.Sin(t)
	ct := math.Cos(t)
	x := v[0]
	y := v[1]
	z := v[2]
	x1 := x*ct + z*st
	y1 := y
	z1 := z*ct - x*st
	return F3{x1, y1, z1}
}

// rotate around +x axis by t(heta) radians
func xrot(v F3, t float64) F3 {
	st := math.Sin(t)
	ct := math.Cos(t)
	x := v[0]
	y := v[1]
	z := v[2]
	x1 := x
	y1 := y*ct + z*st
	z1 := z*ct - y*st
	return F3{x1, y1, z1}
}
func update() {
	for i := range fileVerts {
		v := fileVerts[i]
		v = yrot(v, yrotd)
		for j := range lightpos {
			lightpos[j] = xrot(lightpos[j], lightrot)
		}
		fileVerts[i] = v
	}
	greyval -= 0.001
}
func benchStart(dur *time.Duration) {
	start := time.Now()
	go func() {
		// terminate process after 120 seconds and report loops
		<-time.After(time.Second * 120)
		*dur = time.Since(start)
	}()
}
func main() {
	zbuff = make([]float64, width*height)
	lightpos = append(lightpos, F3{2, 0.5, 1})
	lightpos = append(lightpos, F3{-2, 0.5, 1})
	//lightpos = append(lightpos, F3{0, 3.5, 1.5})
	//lightpos = append(lightpos, F3{0, -3.5, 1.5})
	lightcolors = append(lightcolors, RED|GREEN|BLUE)
	lightcolors = append(lightcolors, GREEN|BLUE)
	//lightcolors = append(lightcolors, RED)
	//lightcolors = append(lightcolors, RED|GREEN)
	lightpower = append(lightpower, 1)
	lightpower = append(lightpower, 1)
	//lightpower = append(lightpower, 0.5)
	//lightpower = append(lightpower, 0.5)
	loadobjfile(filename)
	for i, v := range fileVerts {
		fileVerts[i] = xrot(v, xrotset)
		fmt.Printf("vtop(%v)=%v\n", v, vtop(v))
	}
	fmt.Printf("vtop(%v)=%v\n", F3{-1, -1, -1}, vtop(F3{-1, -1, -1}))
	fmt.Printf("vtop(%v)=%v\n", F3{1, 1, 1}, vtop(F3{1, 1, 1}))
	//test zpixel
	//zpixeldebug = true
	v1 := F3{0, 0, 0}
	v2 := F3{1, 1, 0}
	v3 := F3{1, 0, 1}
	_, _ = zpixel(v1, v2, v3, [2]int{3 * width / 4, 4 * height / 6})
	_, _ = zpixel(v2, v1, v3, [2]int{3 * width / 4, 4 * height / 6})
	zval, err := zpixel(v3, v2, v1, [2]int{3 * width / 4, 4 * height / 6})

	fmt.Printf("zval: v1=%v, v2=%v, v3=%v\tz=%v\terr=%v\n", v1, v2, v3, zval, err)
	//test cross
	norm := cross(v2, v3)
	fmt.Printf("cross of %v %v = %v", v2, v3, norm)
	norm = vnormalize(norm)
	fmt.Printf("after normalization: %v\n", norm)
	//test dynamicNormalForFace
	dn := DynamicNormalForFace(v1, v2, v3)
	fmt.Printf("normal for triangle %v %v %v = %v\n", v1, v2, v3, dn)
	v3 = F3{1, 0, 0}
	dn = DynamicNormalForFace(v1, v2, v3)
	fmt.Printf("normal for triangle %v %v %v = %v\n", v1, v2, v3, dn)
	av := vavg(v1, v2, v3)
	fmt.Printf("vavg = %v\n", av)
	// fmt.Println(fileVerts)
	// fmt.Println(fileFaces)
	mainLoop()
}
func mainLoop() {
	//setup window
	var winTitle string = "TinyRenderer"
	var err error
	window, err = sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(width), int32(height), sdl.WINDOW_SHOWN|sdl.WINDOW_ALWAYS_ON_TOP)
	if err != nil {
		log.Fatal(err)
	}
	//get drawing surface from window
	surf, err := window.GetSurface()
	if err != nil {
		log.Fatal(err)
	}
	//create blank surface to blit before redrawing
	blanksurf, err := sdl.CreateRGBSurface(5, 800, 800, 32, 0, 0, 0, 0)
	if err != nil {
		log.Fatal(err)
	}

	defer window.Destroy()

	//benchmarking
	benchStart(&dur)

	//draw loop
	for {
		loops++
		draw(surf, blanksurf)
		takeKeyboardInput()
		if done {
			end()
			return
		}
		sdl.Delay(delay)
	}

}

func draw(surf *sdl.Surface, blank *sdl.Surface) {
	rect := sdl.Rect{0, 0, width, height}
	blank.Blit(&rect, surf, &rect)
	surf.Lock()
	pix := surf.Pixels()
	DrawLine(0, 0, width, height, greyscale(greyval), pix)
	// fmt.Println(pix)
	update()
	zbuff = make([]float64, width*height)
	for i := range zbuff {
		zbuff[i] = -1000
	}
	drawFrame(pix)
	surf.Unlock()
	window.UpdateSurface()

}
func drawFrame(pix []byte) {
	// for i := range width {
	// 	for j := range height {
	// 		putpixel(i, j, uint32(i*j), pix)
	// 	}
	// }
	// DrawLine(0, 0, width, height, RED|BLUE|GREEN, pix)
	//draw line between vertices
	for _, f := range fileFaces {
		globalcolor = RED | GREEN | BLUE
		i1, i2, i3 := f[0], f[1], f[2]
		v1, v2, v3 := fileVerts[i1-1], fileVerts[i2-1], fileVerts[i3-1]
		//b := pixelbox(v1, v2, v3)
		globalcolor = RED
		if wireframe {
			vline(v1, v2, pix)
			vline(v2, v3, pix)
			vline(v3, v1, pix)
		}
		vn1 := DynamicNormalForFace(v1, v2, v3)
		//vn0 := vavg(v1, v2, v3)
		if vn1[2] < 0 {
			globalcolor = BLUE
		}
		//vline(vn0, vadd(vn0, vn1), pix)
		triangleBoxShader(v1, v2, v3, pix)
		//	DrawLine(b.x0, b.y0, b.x0, b.y1, GREEN|BLUE, pix)
		//	DrawLine(b.x0, b.y1, b.x1, b.y1, GREEN|BLUE, pix)
		//	DrawLine(b.x1, b.y1, b.x1, b.y0, GREEN|BLUE, pix)
		//	DrawLine(b.x1, b.y0, b.x0, b.y0, GREEN|BLUE, pix)
	}
}

// func getVertexInterpShader(a,b,c F3, []uint32 cls) func(x,y int) float32{
//
// }
// triangle drawing func, does z-buffering
func triangleBoxShader(a, b, c F3, pix []byte) {
	tbox := pixelbox(a, b, c)
	for i := tbox.x0; i <= tbox.x1; i++ {
		for j := tbox.y0; j <= tbox.y1; j++ {
			z, err := zpixel(a, b, c, [2]int{i, j})
			if err != nil {
				//fmt.Fprintf(os.Stderr, "boxshader: %v\n", err)
				continue
			}
			//if z != 0 {
			//	fmt.Printf("z= %v\n", z)
			//}
			if zbuff[i+j*width] >= z {
				continue
			}
			zbuff[i+j*width] = z
			//putpixel(i, j, greyscale(z), pix)
			vn1 := DynamicNormalForFace(a, b, c)
			var lightConts uint32
			for lidx, src := range lightpos {
				pow := lightpower[lidx]
				intensity := greyscale(dot(vn1, src) * pow)
				col := lightcolors[lidx]
				if !colorEnabled {
					col = RED | GREEN | BLUE
				}
				lightConts = lightConts | (intensity & col)
			}

			//putpixel(i, j, greyscale(math.Abs(vn1[2])), pix)
			if !shadingEnabled {
				continue
			}
			putpixel(i, j, lightConts, pix)

		}
	}

}

var globalcolor uint32

// raw line from vertex a to vertex b using globalcolor
func vline(a, b F3, pixels []byte) {
	va, vb := vtop(a), vtop(b)
	p1x, p1y := va[0], va[1]
	p2x, p2y := vb[0], vb[1]
	DrawLine(p1x, p1y, p2x, p2y, globalcolor, pixels)
}

// bounding box type: 0 = min, 1 = max
type box struct {
	x0 int
	x1 int
	y0 int
	y1 int
}

// get pixel bounding box for vertices
func pixelbox(vs ...F3) box {
	b := box{}
	b.x0, b.y0 = 100000, 100000
	b.x1, b.y1 = -1, -1
	for _, v := range vs {
		p := vtop(v)
		b.x0, b.y0 = min(b.x0, p[0]), min(b.y0, p[1])
		b.x1, b.y1 = max(b.x1, p[0]), max(b.y1, p[1])
	}
	return b
}

type zpixelerror struct {
	err string
}

var offTriangleError = zpixelerror{"not on triangle"}
var flatTriangleError = zpixelerror{"Flat Triangle: vertices do not define a plane."}

func (e zpixelerror) Error() string {
	return e.err
}

var zpixeldebug bool

// pixel-to-vertex liner transform - assigning z=0
func ptov(p [2]int) F3 {
	x := float64(p[0] - width/2)
	y := float64(p[1] - height/2)
	x /= width / 2
	y /= width / 2
	return F3{x, y, 0}
}

// get z value of pixel px when projected onto triangle made of vertices v0,v1,v2. If px does not fall on the triangle, set err to OFFTRIANGLE
func zpixel(v0, v1, v2 F3, px [2]int) (z float64, err error) {
	if zpixeldebug {
		fmt.Printf("\n\n\nnzpixeldebug: zpixel called with v0=%v, v1=%v, v2=%v, px=%v\n", v0, v1, v2, px)
		defer func() {
			fmt.Printf("returned z=%v, err=%v\n", z, err)
		}()
	}
	//translate triangle verts so v0 -> 0
	u := vdiff(v1, v0)
	w := vdiff(v2, v0)
	//get z offset
	z0 := v0[2]
	//convert pixel to vert
	pv := ptov(px)
	if zpixeldebug {
		fmt.Printf("ptov(%v)=%v\n", px, pv)
	}
	//shift to match v0,1,2; after this pv has z=-p0
	pv = vdiff(pv, v0)
	if zpixeldebug {
		fmt.Printf("pv after xy shift: %v\n", pv)
	}
	if zpixeldebug {
		fmt.Printf("zpixeldebug: vectors: u=<%v>, w=<%v>", u, w)
		fmt.Printf("zpixeldebug: constant offset: v0=<%v>", v0)

	}
	a := u[0]
	b := w[0]
	c := u[1]
	d := w[1]
	x := pv[0]
	y := pv[1]
	//make sure determinant is ! = 0
	if a*d-b*c == 0 {
		return 0, flatTriangleError
	}
	det := 1 / (a*d - b*c)
	//change 2D basis for pv from X,Y to u,v.
	//(calculate coefficients for vectors v,w relative to z0)
	cu := d*det*x - b*det*y
	cw := -c*det*x + a*det*y
	if zpixeldebug {
		fmt.Printf("zpixeldebug: barycentric coordinates: [%v %v] => %v<%v %v> + %v<%v %v>\n", x, y, cu, u[0], u[1], cw, w[0], w[1])
	}
	//make sure pv is inside of the triangle
	if cu < 0 || cw < 0 || cu+cw > 1 {
		return 0, offTriangleError
	}
	//adjust z-coordinates for pv again, based on new basis relative to z0
	z = z0 + u[2]*cu + w[2]*cw
	if zpixeldebug {
		fmt.Printf("zpixeldebug: final zval = %v\n", z)
	}
	return z, nil
}

// scale vector
func vscale(u F3, c float64) F3 {
	return F3{u[0] * c, u[1] * c, u[2] * c}
}

// add vectors
func vadd(u, v F3) F3 {
	return F3{u[0] + v[0], u[1] + v[1], u[2] + v[2]}
}

// invert vector
func vinv(v F3) F3 {
	return F3{-v[0], -v[1], -v[2]}
}

// average of vectors
func vavg(vs ...F3) F3 {
	var avg F3
	for _, v := range vs {
		avg[0] += v[0]
		avg[1] += v[1]
		avg[2] += v[2]
	}
	div := float64(len(vs))
	avg[0] = avg[0] / div
	avg[1] = avg[1] / div
	avg[2] = avg[2] / div
	return avg
}

// dot product
func dot(a, b F3) float64 {
	x, y, z := a[0]*b[0], a[1]*b[1], a[2]*b[2]
	return x + y + z
}

// vector cross product
func cross(a, b F3) F3 {
	a1, a2, a3 := a[0], a[1], a[2]
	b1, b2, b3 := b[0], b[1], b[2]
	x := a2*b3 - a3*b2
	y := a3*b1 - a1*b3
	z := a1*b2 - a2*b1
	return F3{x, y, z}
}

// normalize vector
func vnormalize(v F3) F3 {
	div := math.Sqrt(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])
	for i := range v {
		v[i] = v[i] / div
	}
	return v
}

// vector subtraction u-v
func vdiff(u, v F3) F3 {
	x := u[0] - v[0]
	y := u[1] - v[1]
	z := u[2] - v[2]
	return F3{x, y, z}
}

// get normal to triangle face - orient towards +z
func DynamicNormalForFace(v1, v2, v3 F3) F3 {
	u := vdiff(v2, v1)
	v := vdiff(v3, v1)
	c := cross(u, v)
	if c[2] > 0 {
		c = vinv(c)
	}
	cn := vnormalize(c)
	return cn
}

// vertex to pixel conversion
func vtop(v F3) [2]int {
	x := int((v[0] + 1) * float64(width) / 2)
	y := int((v[1] + 1) * float64(height) / 2)
	return [2]int{x, y}
}
func DrawLine(x0, y0, x1, y1 int, color uint32, pixels []byte) {
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

// return color: <=0->black, 1->white
func greyscale(i float64) uint32 {
	if i < 0 {
		i = 0
	}
	if i > 1 {
		i = 1
	}
	r := uint32(i*RED) & RED
	g := uint32(i*GREEN) & GREEN
	b := uint32(i*BLUE) & BLUE
	return r | b | g
}
func putpixel(x, y int, color uint32, pixels []byte) {
	if x >= width || y >= height || x < 0 || y < 0 {
		//fmt.Fprintf(os.Stderr, "Out of bounds putpixel: %v,%v\n", x, y)
		return
	}
	idx := 4 * (x + width*y)
	b := byte(color & 0xff000000 >> 24)
	g := byte(color & 0x00ff0000 >> 16)
	r := byte(color & 0x0000ff00 >> 8)
	a := byte(color & 0x000000ff)
	pixels[idx] = b
	pixels[idx+1] = g
	pixels[idx+2] = r
	pixels[idx+3] = a
}
func takeKeyboardInput() {
	if event := sdl.PollEvent(); event != nil {
		if event, ok := event.(*sdl.KeyboardEvent); ok {
			if event.State == 0 {
				return
			}
			fmt.Printf("event.Keysym.Scancode: %v %[1]c\n", event.Keysym.Scancode)
			switch event.Keysym.Scancode {
			case 20: //q
				done = true
			case 82: //up
			case 81: //down
			case 80: //left
			case 79: //right
			case 26: //'w'
				wireframe = !wireframe
			case 6: //'c'
				colorEnabled = !colorEnabled
			case 22: //'s'
				shadingEnabled = !shadingEnabled
			}

		}
		if _, ok := event.(*sdl.QuitEvent); ok {
			done = true
			return
		}
	}
}
func end() {
	var loopsPerSec = float64(loops) / float64(dur.Seconds())
	// fmt.Println(surf.Pixels())
	runtime.GC() // get up-to-date statistics
	pprof.WriteHeapProfile(file)
	fmt.Printf("dur=%v; loops=%v; lps=%v\n", dur.Seconds(), loops, loopsPerSec)
}
