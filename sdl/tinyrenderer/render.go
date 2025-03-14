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
	filename      = "african_head.obj"
	//filename = "square.obj"
	//filename        = "tri.obj"
	texturefilename = "african_head_diffuse.tga"
	delay           = 0     //delay between update calls
	yrotd           = 0.010 // +y azis rotation per frame
	xrotset         = 0     // +x azis rotation
	ALPHA           = 0x000000ff
	RED             = 0x0000ff00
	GREEN           = 0x00ff0000
	BLUE            = 0xff000000
)

// offsets for bit shifting color channels to fit in byte
const (
	AOFF = 8 * iota
	ROFF
	GOFF
	BOFF
)

var window *sdl.Window
var file *os.File
var verts []F3
var faces []Face
var texFaces [][3]int
var textureVerts []F3
var texw, texh int
var done bool //control graceful program exit
var loops uint64
var texture []uint32
var blanksurf *sdl.Surface
var surf *sdl.Surface
var zbuff []float64
var zmask []uint32 //0xffffff00 where triangle is visible, otherwise 0x0
var start time.Time
var wireframe bool
var colorEnabled bool
var shadingEnabled bool
var cpuprofile bool
var textureEnabled bool
var lightingEnabled bool

type F3 [3]float64

// cache important properties calculated for Face
type FaceCache struct {
	v0 F3      //vertex 0
	uz float64 //V0->V1 delta z
	wz float64 //V0->V2 delta z
	//Matrix Coefficients
	A1 float64
	A2 float64
	A3 float64
	A4 float64
}

// invalidate cache
func (f *Face) Unload() {
	f.CacheOk = false
}

type Face struct {
	vidx    [3]int //vertex indices (in vertices global array)
	tidx    [3]int //texture-vertex indices (in textureVertices global array)
	CacheOk bool
	Cache   *FaceCache
}

// get texture color for pixel x,y based on texture coordinate interpolation
func (f *Face) TexAt(x, y int) uint32 {
	//get pixel as linear combination of vertices
	b, _, err := f.Project(x, y)
	if err == flatTriangleError {
		return 0x0 //invalid projection gives black color, whatever
	}
	//combine 2D texture coordinates, weighted by barycentric coords of pixel, to get final texture coordinate
	vs := f.T()
	v0, v1, v2 := vs[0], vs[1], vs[2]
	vx := b[0]*v0[0] + b[1]*v1[0] + b[2]*v2[0]
	vy := b[0]*v0[1] + b[1]*v1[1] + b[2]*v2[1]
	return textureAt(1-vy, vx)
}

// get vertices for face
func (f *Face) V() [3]F3 {
	var v [3]F3
	for i := range 3 {
		idx := f.vidx[i]
		v[i] = verts[idx-1]
	}
	return v
}

// get texture vertices for face (in this case our texture has z=0 for all points in the texture)
func (f *Face) T() [3]F3 {
	var t [3]F3
	for i := range 3 {
		idx := f.tidx[i]
		t[i] = textureVerts[idx-1]
	}
	return t
}

func (f *Face) Norm() F3 {
	vs := f.V()
	v0, v1, v2 := vs[0], vs[1], vs[2]
	u := vdiff(v1, v0)
	w := vdiff(v2, v0)
	norm := cross(u, w)
	return vnormalize(norm)
}

// TODO: check performance vs Project
// Project when basis-change matrix already calculated, z0, Uz, Wz cached.
func (f *Face) ProjectCached(x, y int) (bc F3, z float64, err error) {
	//convert pixel to vert
	pv := ptov([2]int{x, y})
	//get pixel vector relative to v0
	pv = vdiff(pv, f.Cache.v0)
	cx := pv[0]
	cy := pv[1]
	//(calculate coefficients for vectors v,w relative to z0)
	cu := f.Cache.A1*cx + f.Cache.A2*cy
	cw := f.Cache.A3*cx + f.Cache.A4*cy
	//make sure pv is inside of the triangle
	if cu < 0 || cw < 0 || cu+cw > 1 {
		return F3{}, 0, offTriangleError
	}
	//adjust z-coordinates for pv based on new basis 2D basis, interpolating across triangle
	z0 := f.Cache.v0[2]
	z = z0 + f.Cache.uz*cu + f.Cache.wz*cw
	//barycentric coords:
	bc = [3]float64{1 - cu - cw, cu, cw}
	return bc, z, nil
}

// get barycentric coordinates and z-coordinate for pixel x,y based on face-projection
func (f *Face) Project(x, y int) (bc F3, z float64, err error) {
	if f.CacheOk {
		return f.ProjectCached(x, y)
	}
	//get face verts
	vs := f.V()
	v0, v1, v2 := vs[0], vs[1], vs[2]
	//translate face verts so v0 -> 0
	u := vdiff(v1, v0)
	w := vdiff(v2, v0)
	//get z coord for base vert
	z0 := v0[2]
	//convert pixel to vert
	pv := ptov([2]int{x, y})
	//get pixel vector relative to v0
	pv = vdiff(pv, v0)
	//express pv in terms of new basis vectors u,w
	//c1*u + c2*w = |ux, wx| |c1| = |x|
	//		|uy, wy| |c2|   |y|
	//
	//We want to invert this to get
	//c1 = x * A^-1
	//c2   y
	cx := pv[0]
	cy := pv[1]
	a := u[0]
	b := w[0]
	c := u[1]
	d := w[1]
	//make sure determinant is ! = 0
	if a*d-b*c == 0 {
		return F3{}, 0, flatTriangleError
	}
	det := 1 / (a*d - b*c)
	//(calculate coefficients for vectors v,w relative to z0)
	cu := d*det*cx - b*det*cy
	cw := -c*det*cx + a*det*cy
	//make sure pv is inside of the triangle
	if cu < 0 || cw < 0 || cu+cw > 1 {
		return F3{}, 0, offTriangleError
	}
	//adjust z-coordinates for pv based on new basis 2D basis, interpolating across triangle
	z = z0 + u[2]*cu + w[2]*cw
	//barycentric coords:
	bc = [3]float64{1 - cu - cw, cu, cw}
	//cache values
	var C FaceCache
	C.v0 = v0
	C.uz = u[2]
	C.wz = w[2]
	C.A1 = d * det
	C.A2 = -b * det
	C.A3 = -c * det
	C.A4 = a * det
	f.Cache = &C
	f.CacheOk = true
	return bc, z, nil
}

func main() {
	setup()
	for !done {
		loops++
		drawFrame(surf, blanksurf)
		update()
		takeKeyboardInput()
		sdl.Delay(delay)
	}
	end()
}
func setup() {
	//set default values for global options
	textureEnabled = true
	shadingEnabled = true
	cpuprofile = true
	//allocate zbuffer
	zbuff = make([]float64, width*height) //keep track of what's in front of scene
	//allocate zmask
	zmask = make([]uint32, width*height+801)
	//load files
	loadobjfile(filename)
	texture = loadTexture(texturefilename)
	//set up window
	var winTitle string = "TinyRenderer"
	var err error
	window, err = sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(width), int32(height), sdl.WINDOW_SHOWN|sdl.WINDOW_ALWAYS_ON_TOP)
	if err != nil {
		log.Fatal(err)
	}
	//get drawing surface from window
	surf, err = window.GetSurface()
	if err != nil {
		log.Fatal(err)
	}
	//allocate reusable blank surface to blit before redrawing
	blanksurf, err = sdl.CreateRGBSurface(5, 800, 800, 32, 0, 0, 0, 0)
	if err != nil {
		log.Fatal(err)
	}
	//profiling
	if cpuprofile {
		f, err := os.Create("cpuprofile")
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
	}
	//benchmarking
	start = time.Now()
}

func update() {
	for i := range verts {
		v := verts[i]
		v = yrot(v, yrotd)
		verts[i] = v
	}
}

// draw loop
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
			case 82: //up arrow
			case 81: //down arrow
			case 80: //left arrow
			case 79: //right arrow
			case 26: //'w'
				wireframe = !wireframe
			case 6: //'c'
				colorEnabled = !colorEnabled //??
			case 22: //'s'
				shadingEnabled = !shadingEnabled
			case 15: //'l'
				lightingEnabled = !lightingEnabled
			case 23: //'t'
				textureEnabled = !textureEnabled
			}

		}
		if _, ok := event.(*sdl.QuitEvent); ok {
			done = true
			return
		}
	}
}

// stop profiling + benchmarking
func end() {
	dur := time.Since(start)
	var loopsPerSec = float64(loops) / float64(dur.Seconds())
	if cpuprofile {
		pprof.StopCPUProfile()
	}
	runtime.GC() // get up-to-date statistics
	file, err := os.Create("memprofile")
	if err != nil {
		log.Fatal(err)
	}
	pprof.WriteHeapProfile(file)
	fmt.Printf("dur=%v; loops=%v; lps=%v\n", dur.Seconds(), loops, loopsPerSec)
}
