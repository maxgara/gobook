package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	width, height = 800, 800 // window dims
	filename      = "african_head.obj"
)

var (
	T0M             *M4
	ID              *M4  // useful for testing matrix ops
	dotsEnabled     bool //only draw vertices
	textureEnabled  bool //draw texture
	lightingEnabled bool
)

type (
	F3 [3]float64
	M4 [4][4]float64 // matrix[row][col]
	V4 struct {
		x float64
		y float64
		z float64
		m float64 //"magic" coordinate
	}
)

// bounding box
type box struct {
	x0 float64
	x1 float64
	y0 float64
	y1 float64
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

// fragment shader
func fshader(i, j int, bcs V4, f Face, v0, v1, v2 V4, tex []uint32, pix []byte) {
	if bcs.x < 0 || bcs.y < 0 || bcs.z < 0 {
		return
	}
	if i+j*width >= len(zbuff) || i+j*width < 0 {
		return
	}
	z := bcs.x*v0.z + bcs.y*v1.z + bcs.z*v2.z
	if zbuff[i+j*width] < z {
		return
	}
	zbuff[i+j*width] = z
	//fmt.Printf("z=%v\n", z)
	// fmt.Printf("zval=%v\n", z)
	chv := byte(min(0xff, 0xff-z*0xff/2.5)) // channel val
	if lightingEnabled {
		//bcs = mvMult(lInterpM, bcs)
		it := bcs.x*f.lvals[0] + bcs.y*f.lvals[1] + bcs.z*f.lvals[2] //get lighting at point i,j from barycentric coords
		//if it < 0 {
		//	it = 0
		//}
		chv = byte(it / 4)
	}
	if textureEnabled {
		txidx := bcs.x*f.txidxs[0] + bcs.y*f.txidxs[1] + bcs.z*f.txidxs[2]
		tyidx := bcs.x*f.tyidxs[0] + bcs.y*f.tyidxs[1] + bcs.z*f.tyidxs[2]
		//tidx := int(float64(tstride)*txidx + tyidx*float64(tstride)*float64(tstride))
		xidxint := int(float64(tstride) * txidx)
		yidxint := int(float64(tstride) * tyidx)
		tidx := tstride*xidxint + tstride - yidxint
		//fmt.Printf("texture coords (%v, %v)\n", txidx, tyidx)
		tcol := tex[tidx]
		tcb := bgraToBytes(tcol)
		tb := tcb[0]
		tg := tcb[1]
		tr := tcb[2]
		putpixel(int(i), int(j), byte(float64(chv)*float64(tr)/255), byte(float64(chv)*float64(tg)/255), byte(float64(chv)*float64(tb)/255), 0, pix)
		return
	}
	putpixel(int(i), int(j), 0, chv, 0, 0, pix)
}
func bgraToBytes(c uint32) [4]byte {
	b := byte(c >> 24)
	g := byte(c >> 16 & 0xff)
	r := byte(c >> 8 & 0xff)
	a := byte(c & 0xff)
	return [4]byte{b, g, r, a}
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
	for _, f := range ob.fs {
		f.txidxs[0] = ob.tvs[f.tidx[0]-1].x
		f.txidxs[1] = ob.tvs[f.tidx[1]-1].x
		f.txidxs[2] = ob.tvs[f.tidx[2]-1].x
		f.tyidxs[0] = ob.tvs[f.tidx[0]-1].y
		f.tyidxs[1] = ob.tvs[f.tidx[1]-1].y
		f.tyidxs[2] = ob.tvs[f.tidx[2]-1].y
		vxs := f.vidx //vertex indices
		vs := []V4{ob.vs[vxs[0]-1], ob.vs[vxs[1]-1], ob.vs[vxs[2]-1]}
		vns := []V4{ob.vns[vxs[0]-1], ob.vns[vxs[1]-1], ob.vns[vxs[2]-1]}
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

	//for _, v := range ob.vs {
	//	putpixel(int(v.x), int(v.y), 0xff, 0xff, 0xff, 0, pix)
	//}
	//done drawing, set up window for display
	surf.Unlock()
	err = window.UpdateSurface()
	if err != nil {
		log.Fatal(err)
	}
}
