package main

import (
	"fmt"
	"log"

	"github.com/veandco/go-sdl2/sdl"
)

// draw a frame
func drawFrame(pix []byte) {
	for _, f := range faces {
		vs := f.V()
		v1, v2, v3 := vs[0], vs[1], vs[2]
		globalcolor = RED
		if wireframe {
			vline(v1, v2, pix)
			vline(v2, v3, pix)
			vline(v3, v1, pix)
		}
		triangleBoxShader(&f, pix, zmask)
	}
}
func testInterpText(pix []byte) {
	var f Face
	f.vidx = [3]int{0, 1, 2}
	f.tidx = [3]int{0, 1, 2}
	for i := range width {
		for j := range height {
			color := interpTexture(&f, i, j)
			//fmt.Printf("color = %x\n", color)
			putpixel(i, j, color, pix)
		}
	}
}
func testTextureAt(pix []byte) {
	var count = 1000
	for i := range count {
		for j := range count {
			x := float64(i) / float64(count)
			y := float64(j) / float64(count)
			color := textureAt(x, y) & 0xffffff00
			fmt.Printf("color = %x\n", color)
			putpixel(int(x*float64(width)), int(y*float64(height)), color, pix)
		}
	}
}

// draw a frame (high level)
func draw(surf *sdl.Surface, blank *sdl.Surface) {
	rect := sdl.Rect{0, 0, width, height}
	blank.Blit(&rect, surf, &rect)
	for i := range zbuff {
		zbuff[i] = -1000
	}
	surf.Lock()
	pix := surf.Pixels()
	//DrawLine(0, 0, width, height, greyscale(greyval), pix)
	// fmt.Println(pix)
	update()
	//testDrawTextureImg(pix)
	//testTextureAt(pix)
	//testInterpText(pix)
	drawFrame(pix)
	surf.Unlock()
	window.UpdateSurface() // }}}
}

// triangle drawing func, does z-buffering.
// zmask is scratch space for computing visible pixels, does not need to be zeroed-out but cannot be changed concurrently with this function.
func triangleBoxShader(f *Face, pix []byte, zmask []uint32) {
	v := f.V()
	a, b, c := v[0], v[1], v[2]
	tbox := pixelbox(a, b, c) // {{{
	err := zpixelboxmask(a, b, c, zmask)
	if err != nil {
		log.Fatal(err)
	}
	vn1 := DynamicNormalForFace(a, b, c)
	var lightConts uint32
	for lidx, src := range lightpos {
		pow := lightpower[lidx]
		intensity := greyscale(dot(vn1, src) * pow)
		col := lightcolors[lidx]
		if !colorEnabled {
			col = RED | GREEN | BLUE
		}
		lightConts = intensity & col
	}

	if !shadingEnabled {
		return
	}
	for i := tbox.x0; i <= tbox.x1; i++ {
		for j := tbox.y0; j <= tbox.y1; j++ {
			//putpixel(i, j, greyscale(math.Abs(vn1[2])), pix)
			maskval := zmask[i+width*j]
			if maskval == 0 {
				continue
			}
			if textureEnabled {
				//estimate texture color by using color at vertex 1.
				//TODO: replace this with extrapolation using barycentric cs.
				//texturecolor := textureFor(aidx)
				texturecolor := interpTexture(f, i, j)
				lightConts = texturecolor
				//v0col :=
			}
			lightConts = lightConts & maskval
			//lightConts = (RED | GREEN | BLUE) & maskval
			putpixel(i, j, lightConts, pix)
		}
	}
}

// build a zmask for v0,v1,v2
func zpixelboxmask(v0, v1, v2 F3, zmask []uint32) (err error) {
	//get bounding box to draw in{{{
	bds := pixelbox(v0, v1, v2)
	//create pixel square to blit
	//translate triangle verts so v0 -> 0
	u := vdiff(v1, v0)
	w := vdiff(v2, v0)
	//get z offset
	z0 := v0[2]

	a := u[0]
	b := w[0]
	c := u[1]
	d := w[1]
	//make sure determinant is ! = 0
	if a*d-b*c == 0 {
		return nil
	}
	det := 1 / (a*d - b*c)
	//get z-pixel values
	var i, j int
	for i = bds.x0; i <= bds.x1; i++ {
		for j := bds.y0; j <= bds.y1; j++ {
			midx := i + j*width
			zmask[midx] = 0
		}
	}
	for i = bds.x0; i <= bds.x1; i++ {
		for j = bds.y0; j <= bds.y1; j++ {
			//convert pixel to vert
			px := [2]int{i, j}
			pv := ptov(px)
			//shift to match v0,1,2
			pv = vdiff(pv, v0)
			x := pv[0]
			y := pv[1]
			//change 2D basis for pv from X,Y to u,v.
			//(calculate coefficients for vectors v,w relative to z0)
			//make sure pv is inside of the triangle
			cu := d*det*x - b*det*y
			cw := -c*det*x + a*det*y
			//get masking index
			midx := i + j*width
			//make sure px is inside of the triangle
			if cu < 0 || cw < 0 || cu+cw > 1 {
				//	zmask[midx] = 0
				//do not continue down y-axis once we have left the triangle
				continue
			}
			//if there is a triangle in front of the current position, don't draw
			z := z0 + u[2]*cu + w[2]*cw
			if zbuff[midx] > z {
				//zmask[midx] = 0
				continue
			}
			zbuff[midx] = z
			//set mask color bits for pixels where triangle should be drawn
			zmask[midx] = RED | GREEN | BLUE

		}

	}
	return nil // }}}
}

// draw line between vertices
func DrawLine(x0, y0, x1, y1 int, color uint32, pixels []byte) {
	if x1 < x0 { // {{{
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
	if i < 0 { // {{{
		i = 0
	}
	if i > 1 {
		i = 1
	}
	r := uint32(i*RED) & RED
	g := uint32(i*GREEN) & GREEN
	b := uint32(i*BLUE) & BLUE
	return r | b | g // }}}
}
func putpixel(x, y int, color uint32, pixels []byte) {
	if x >= width || y >= height || x < 0 || y < 0 { // {{{
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
	pixels[idx+3] = a // }}}
}
