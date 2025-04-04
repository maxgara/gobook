package main

import (
	"errors"
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
	texturefilename = "african_head_diffuse.tga"
	delay           = 0     // delay between update calls
	yrotDelta       = 0.010 // +y azis rotation per frame
	xrotset         = 0     // +x azis rotation
	ALPHA           = 0x000000ff
	RED             = 0x0000ff00
	GREEN           = 0x00ff0000
	BLUE            = 0xff000000
	width, height   = 800, 800 // window dims
	filename        = "african_head.obj"
)

// offsets for bit shifting color channels to fit in byte
const (
	AOFF = 8 * iota
	ROFF
	GOFF
	BOFF
)

var (
	T0M               *M4
	ID                *M4  // useful for testing matrix ops
	dotsEnabled       bool //only draw vertices
	textureEnabled    bool //draw texture
	lightingEnabled   bool
	altFShaderEnabled bool
	window            *sdl.Window
	tstride           int
	done              bool // control graceful program exit
	loops             uint64
	texture           []uint32
	blanksurf         *sdl.Surface
	surf              *sdl.Surface
	zbuff             []float64
	start             time.Time
	cpuprofile        bool
	yrotTot           float64
)

type (
	//4D matrix
	M4 [4][4]float64 // matrix[row][col]
	//4D vector
	V4 struct {
		x float64
		y float64
		z float64
		m float64 //"magic" coordinate
	}
)

// 2D bounding box
type box struct {
	x0 float64
	x1 float64
	y0 float64
	y1 float64
}

func main() {
	ob := setup()
	// transform vertices so that they are in the screen area
	for !done {
		loops++
		drawFrame(surf, blanksurf, ob)
		update(ob)
		takeKeyboardInput()
		sdl.Delay(delay)
	}
	end()
}
func bgraToBytes(c uint32) [4]byte {
	b := byte(c >> 24)
	g := byte(c >> 16 & 0xff)
	r := byte(c >> 8 & 0xff)
	a := byte(c & 0xff)
	return [4]byte{b, g, r, a}
}
func takeKeyboardInput() {
	if event := sdl.PollEvent(); event != nil {
		if event, ok := event.(*sdl.KeyboardEvent); ok {
			if event.State == 0 {
				return
			}
			fmt.Printf("event.Keysym.Scancode: %v %[1]c\n", event.Keysym.Scancode)
			switch event.Keysym.Scancode {
			case 20: // q
				done = true
			case 7: // d
				dotsEnabled = !dotsEnabled
			case 23: //t
				textureEnabled = !textureEnabled
			case 15: //l
				lightingEnabled = !lightingEnabled
			}
		}
		if _, ok := event.(*sdl.QuitEvent); ok {
			return
		}
	}
}

// get matrix for rotation around +y axis by t(heta) radians
func getYRot(t float64) M4 {
	st := math.Sin(t)
	ct := math.Cos(t)
	rx := [4]float64{ct, 0, st, 0}
	ry := [4]float64{0, 1, 0, 0}
	rz := [4]float64{-st, 0, ct, 0}
	rm := [4]float64{0, 0, 0, 1}
	//	x1 := v.x*ct + v.z*st
	//	y1 := v.y
	//	z1 := v.z*ct - v.x*st
	return M4{rx, ry, rz, rm}
}

var vertexRTM M4 //vertex rotation and translation matrix
// call this before each frame render
func update(ob *Obj) {
	//get rotation Matrix
	M := getYRot(yrotTot)
	//get translation matrix (move object away from camera)
	TM := *getTransM(0, 0, 2)
	vertexRTM = mmMult(TM, M) // rotate before translation (rotation on the right)
	yrotTot += yrotDelta
}

// put pixel with color in pixels array at pos x,y
func putpixel(x, y int, r, g, b, a byte, pixels []byte) {
	if x >= width || y >= height || x < 0 || y < 0 {
		// fmt.Fprintf(os.Stderr, "Out of bounds putpixel: %v,%v\n", x, y)
		return
	}
	idx := 4 * (x + width*y)
	pixels[idx] = b
	pixels[idx+1] = g
	pixels[idx+2] = r
	pixels[idx+3] = a
}

// stop profiling + benchmarking
func end() {
	dur := time.Since(start)
	loopsPerSec := float64(loops) / float64(dur.Seconds())
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

// get pixel bounding box for vertices
func pixelbox(vs ...V4) box {
	b := box{}
	b.x0, b.y0 = 100000, 100000
	b.x1, b.y1 = -100000, -100000
	for _, v := range vs {
		b.x0, b.y0 = min(b.x0, v.x), min(b.y0, v.y)
		b.x1, b.y1 = max(b.x1, v.x), max(b.y1, v.y)
	}
	return b
}

func VtoF(v V4) [4]float64 {
	return [4]float64{v.x, v.y, v.z, v.m}
}

func FtoV(f [4]float64) V4 {
	return V4{f[0], f[1], f[2], f[3]}
}

type Face struct {
	vidx   [3]int // vertex indices (in vertices global array)
	tidx   [3]int // texture-vertex indices (in textureVertices global array)
	txidxs [3]float64
	tyidxs [3]float64
	lvals  [3]float64 //lighting per vertex
}

// represents a 3D object and associated texture. Affine transformation matrix can be applied via Transform and will affect the drawn image, while keeping original data intact. Additional Transforms will use original data.
type Obj struct {
	fileVs []V4 // original vertex data
	fileNs []V4 //original vertex-normal data
	vs     []V4 // vertex data Transform
	vns    []V4 // vertex normal data
	//tex    []uint32 // texture data
	tvs []V4 // texture vertex data
	//tw     int      // texture width
	//th     int      // texture height
	fs []Face // faces
}

func (A M4) Transform(vsOut, vsIn []V4) {
	for i, v := range vsIn {
		vsOut[i] = mvMult(A, v)
	}
}

// get a transformation matrix to translate a vertex by <x,y,z>
func getTransM(x, y, z float64) *M4 {
	A := M4{
		[4]float64{1, 0, 0, x},
		[4]float64{0, 1, 0, y},
		[4]float64{0, 0, 1, z},
		[4]float64{0, 0, 0, 1},
	}
	return &A
}

// populate global T0 Matrix, which moves points into the visibile area of the screen by mapping < -1, -1, *, *> to <0, 0, *, *> and <1, 1, *, *> to <width, height, *, *>
func getT0M() *M4 {
	w := float64(width)
	h := float64(height)
	A := M4{[4]float64{w / 2, 0, 0, w / 2}, // x scales to width
		[4]float64{0, h / 2, 0, h / 2}, // y scales to height
		[4]float64{0, 0, 1, 0},         // z unchanged
		[4]float64{0, 0, 0, 1}}         // translate right and down
	return &A
}

// transform verts to fill screen space
func (ob *Obj) T0() {
	w := float64(width)
	h := float64(height)
	A := M4{[4]float64{w / 2, 0, 0, w / 2}, // x scales to width
		[4]float64{0, h / 2, 0, h / 2}, // y scales to height
		[4]float64{0, 0, 1, 0},         // z unchanged
		[4]float64{0, 0, 0, 1}}         // translate right and down
	A.Transform(ob.vs, ob.fileVs)
}

// multiply vector by matrix
func mvMult(A M4, v V4) V4 {
	var nv [4]float64
	for i := range 4 {
		r := A[i]
		nv[i] = v.x*r[0] + v.y*r[1] + v.z*r[2] + v.m*r[3]
	}
	return V4{nv[0], nv[1], nv[2], nv[3]}
}

// matrix entry i,j -> j,i
func (A *M4) idxSwap() M4 {
	var out M4
	for i := range 4 {
		for j := range 4 {
			out[j][i] = A[i][j]
		}
	}
	return out
}

// matrix multiplication
func mmMult(A M4, B M4) M4 {
	var Csw M4
	Bsw := B.idxSwap()
	for i := range 4 { // range over cols of B to make vectors
		bcol := Bsw[i]
		newcol := mvMult(A, FtoV(bcol)) // get a new column from each A x B_i
		Csw[i] = VtoF(newcol)           // add new column to c-swap matrix as a row
		// place each new entry in the output matrix
	}
	return Csw.idxSwap()
}

// add two vectors together
func vadd(u, v V4) V4 {
	var out V4
	out.x = u.x + v.x
	out.y = u.y + v.y
	out.z = u.z + v.z
	out.m = u.m + v.m
	return out
}

// subtract v from u
func vsub(u, v V4) V4 {
	var out V4
	out.x = u.x - v.x
	out.y = u.y - v.y
	out.z = u.z - v.z
	out.m = u.m - v.m
	return out
}

// get matrix representing interpolation of a quantity q between face vertices
// matrix represents transformation (u,w) -> (u,v,q)
func getUVInterpolationM(q0, q1, q2 float64) M4 {
	qu := q1 - q0
	qw := q2 - q0
	r0 := [4]float64{1, 0, 0, 1}
	r1 := [4]float64{0, 1, 0, 0}
	r2 := [4]float64{qu, qw, 0, q0}
	r3 := [4]float64{0, 0, 0, 1}
	return M4{r0, r1, r2, r3}
}

// get matrix to change basis for point to barycentric coords
// matrix represents transformation (x,y) -> (c0,c1,c2)
func getBaryM(v0, v1, v2 V4) (M M4, err error) {
	// get AB and AC vectors
	u := vsub(v1, v0)
	w := vsub(v2, v0)
	// get matrix entry helper vars
	a := u.x
	b := w.x
	c := u.y
	d := w.y
	// make sure determinant is ! = 0
	if a*d-b*c < 0.01 && a*d-b*c > -0.01 {
		err = errors.New("triangle does not define a plane")
		return
	}
	det := 1 / (a*d - b*c)

	// get translation matrix for x,y offset from v0
	T := getTransM(-v0.x, -v0.y, 0)

	// get partial change of basis matrix (x,y only, z set to 0)
	r0 := [4]float64{d * det, -b * det, 0, 0}
	r1 := [4]float64{-c * det, a * det, 0, 0}
	r2 := [4]float64{0, 0, 0, 0}
	r3 := [4]float64{0, 0, 0, 1}
	M = M4{r0, r1, r2, r3}

	// change u,w coords to barycentric coords c0, c1, c2
	r0 = [4]float64{-1, -1, 0, 1} //1 - c1 - c2 (p = c0v0 + c1v1 + c2v2); c0 + c1 + c2 = 1
	r1 = [4]float64{1, 0, 0, 0}
	r2 = [4]float64{0, 1, 0, 0}
	r3 = [4]float64{0, 0, 0, 1}
	ZB := M4{r0, r1, r2, r3}

	M = mmMult(ZB, M)
	return mmMult(M, *T), nil
}

// some useful snippits for testing
func testcalls() {
	v0 := V4{1, 0, 0, 0}
	v1 := V4{0, 1, 0, 0}

	r0 := [4]float64{1, 0, 0, 0}
	r1 := [4]float64{0, 2, 0, 0}
	r2 := [4]float64{0, 0, 1, 0}
	r3 := [4]float64{0, 0, 0, 1}
	A := M4{r0, r1, r2, r3}
	fmt.Println(mvMult(A, vadd(v0, v1)))
}

// change x,y coords of vs to project onto a plane where z=0, based on perspective from origin
func perspectiveProject(v V4) V4 {
	v.x = v.x / v.z
	v.y = v.y / v.z
	return v
}

// dot product
func dot(a, b V4) float64 {
	x, y, z := a.x*b.x, a.y*b.y, a.z*b.z
	return x + y + z
}

// draw a frame
var lInterpM M4
var tInterpM M4

func drawFrame(surf *sdl.Surface, blank *sdl.Surface, ob *Obj) {
	// clear screen
	rect := sdl.Rect{0, 0, width, height}
	blank.Blit(&rect, surf, &rect)
	// reset zbuffer
	for i := range zbuff {
		zbuff[i] = 1000
	}
	// get pixel array
	err := surf.Lock()
	if err != nil {
		log.Fatal(err)
	}
	pix := surf.Pixels()
	// draw
	if dotsEnabled {
		for i := range ob.vs {
			x, y := ob.vs[i].x, ob.vs[i].y
			putpixel(int(x), int(y), 0, 0xff, 0xff, 0, pix)
		}
		surf.Unlock()
		err = window.UpdateSurface()
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if altFShaderEnabled {

		var size float64 = 30
		for _, v := range ob.vs {
			dummynorm := V4{0, 1, 0, 0}
			vShader(&v, &dummynorm)
			x := v.x
			y := v.y
			for j := max(0, y-size); j <= min(y+size, height-1); j++ {
				for i := max(0, x-size); i <= min(x+size, width-1); i++ {
					altfShader(int(i), int(j), v, texture, size, pix)
				}
			}
		}
		surf.Unlock()
		err = window.UpdateSurface()
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	for _, f := range ob.fs {
		for i := range 3 {
			f.txidxs[i] = ob.tvs[f.tidx[i]-1].x
			f.tyidxs[i] = ob.tvs[f.tidx[i]-1].y
		}
		vxs := f.vidx                                            //vertex indices
		ovs := ob.vs                                             //object vertices
		vs := []V4{ovs[vxs[0]-1], ovs[vxs[1]-1], ovs[vxs[2]-1]}  //face vertices
		vns := []V4{ovs[vxs[0]-1], ovs[vxs[1]-1], ovs[vxs[2]-1]} //face vertex normals
		vShader(&vs[0], &vns[0])
		vShader(&vs[1], &vns[1])
		vShader(&vs[2], &vns[2])
		bbox := pixelbox(vs...)
		M, err := getBaryM(vs[0], vs[1], vs[2])
		if err != nil {
			continue
		}
		it0, it1, it2 := vs[0].m, vs[1].m, vs[2].m
		if lightingEnabled {
			//get intensities for each vertex in face
			f.lvals[0] = it0
			f.lvals[1] = it1
			f.lvals[2] = it2
		}
		//tcs := f.tidx[0], f.tidx[1], f.tidx[2]
		//tInterpM = getUVInterpolationM(f.tidx)

		for j := int(bbox.y0); j <= int(bbox.y1); j++ {
			for i := int(bbox.x0); i <= int(bbox.x1); i++ {
				// get barycentric coords for <i,j>
				v := V4{float64(i), float64(j), 0, 1}
				bcs := mvMult(M, v)
				fshader(i, j, bcs, f, vs[0], vs[1], vs[2], texture, pix)
			}
		}
	}

	//done drawing, set up window for display
	surf.Unlock()
	err = window.UpdateSurface()
	if err != nil {
		log.Fatal(err)
	}
}
