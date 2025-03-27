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
	vidx [3]int // vertex indices (in vertices global array)
	tidx [3]int // texture-vertex indices (in textureVertices global array)
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

// get matrix to change basis for point to barycentric coords
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

	// get partial reverse change of basis (z extrapolated from u,v)
	R0 := [4]float64{1, 0, 0, 0}
	R1 := [4]float64{0, 1, 0, 0}
	R2 := [4]float64{u.z, w.z, 0, v0.z}
	R3 := [4]float64{0, 0, 0, 1}
	ZB := M4{R0, R1, R2, R3}

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
func perspectiveProject(vs []V4) {
	for i := range vs {
		v := vs[i]
		v.x = v.x / v.z
		v.y = v.y / v.z
		vs[i] = v
	}
}

// dot product
func dot(a, b V4) float64 {
	x, y, z := a.x*b.x, a.y*b.y, a.z*b.z
	return x + y + z
}

// draw a frame
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
		vxs := f.vidx //vertex indices
		vs := []V4{ob.vs[vxs[0]-1], ob.vs[vxs[1]-1], ob.vs[vxs[2]-1]}
		bbox := pixelbox(vs...)
		_ = bbox
		M, err := getBaryM(vs[0], vs[1], vs[2])
		if err != nil {
			continue
		}
		if lightingEnabled {
			//TODO create matrix to interpolate xyz coords for normals across face

		}
		for j := int(bbox.y0); j <= int(bbox.y1); j++ {
			for i := int(bbox.x0); i <= int(bbox.x1); i++ {
				// get barycentric coords for <i,j>
				v := V4{float64(i), float64(j), 0, 1}
				varrout := make([]V4, 1)
				M.Transform(varrout, []V4{v})
				bcs := varrout[0]
				if bcs.x < 0 || bcs.y < 0 || 1-bcs.x-bcs.y < 0 {
					continue
				}
				// vec0 := vsub(vs[1], vs[0])
				// vec1 := vsub(vs[2], vs[0])

				// z := bcs.x*vec0.z + bcs.y*vec1.z
				z := bcs.z
				if i+j*width >= len(zbuff) || i+j*width < 0 {
					continue
				}
				if zbuff[i+j*width] < z {
					continue
				}
				zbuff[i+j*width] = z
				// fmt.Printf("zval=%v\n", z)
				if z > 2 {
					z = 2
				}
				chv := byte(min(0xff, 0xff-z*0xff/2.5)) // channel val
				if lightingEnabled {
					vn0 := ob.vns[f.vidx[0]-1]
					it := 0xff * dot(vn0, V4{3, 0, -1, 0})
					if it < 0 {
						it = 0
					}
					chv = byte(it / 5)
				}
				putpixel(int(i), int(j), chv, chv, chv, 0, pix)
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
