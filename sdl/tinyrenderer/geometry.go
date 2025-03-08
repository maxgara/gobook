package main

import (
	"fmt"
	"math"
)

var globalcolor uint32

type zpixelerror struct {
	err string
}

var offTriangleError = zpixelerror{"not on triangle"}
var flatTriangleError = zpixelerror{"Flat Triangle: vertices do not define a plane."}

// zpixel error message
func (e zpixelerror) Error() string {
	return e.err
}

var zpixeldebug bool

// barycentric coordinates for pixel; err == nil if barycentric coordinates found and point is in the triangle.
// err != nil if point is outside of triangle, but coordinates will still be correct.
func bary(v0, v1, v2 F3, px, py int) (F3, error) {
	//translate triangle verts so v0 -> 0
	u := vdiff(v1, v0)
	w := vdiff(v2, v0)
	//convert pixel to vert
	pv := ptov([2]int{px, py})
	//shift to match v0,1,2; after this pv has z=-p0
	pv = vdiff(pv, v0)
	a := u[0]
	b := w[0]
	c := u[1]
	d := w[1]
	x := pv[0]
	y := pv[1]
	//make sure determinant is ! = 0
	if a*d-b*c == 0 {
		return F3{0, 0, 0}, flatTriangleError
	}
	det := 1 / (a*d - b*c)
	//change 2D basis for pv from X,Y to u,v.
	//(calculate coefficients for vectors v,w relative to z0)
	cu := d*det*x - b*det*y
	cw := -c*det*x + a*det*y
	// third barycentric coordinate (normalized so the at Sum = 1)
	c0 := 1 - cu - cw
	// if we are outside the triangle, return an error but also include the correct barycentric coordinates
	if cu < 0 || cw < 0 || c0 < 0 {
		return F3{c0, cu, cw}, offTriangleError
	}
	return F3{c0, cu, cw}, nil
}

// get z value of pixel px when projected onto triangle made of vertices v0,v1,v2. If px does not fall on the triangle, set err to OFFTRIANGLE
func zpixel(v0, v1, v2 F3, px [2]int) (z float64, err error) {
	if zpixeldebug {
		//fmt.Printf("\n\n\nnzpixeldebug: zpixel called with v0=%v, v1=%v, v2=%v, px=%v\n", v0, v1, v2, px)
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
func vadd(vecs ...F3) F3 {
	var nv F3
	for _, v := range vecs {
		nv = F3{nv[0] + v[0], nv[1] + v[1], nv[2] + v[2]}
	}
	return nv
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

// pixel-to-vertex liner transform - assigning z=0
func ptov(p [2]int) F3 {
	x := float64(p[0] - width/2)
	y := float64(p[1] - height/2)
	x /= width / 2
	y /= width / 2
	return F3{x, y, 0}
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

//interpolate color 3 color value based on barycentric coordinates
