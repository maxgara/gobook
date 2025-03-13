package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

func testInterpText(pix []byte) {
	var f Face
	f.vidx = [3]int{0, 1, 2}
	f.tidx = [3]int{0, 1, 2}
	for i := range width {
		for j := range height {
			color := interpTexture(&f, i, j)
			// fmt.Printf("color = %x\n", color)
			putpixel(i, j, color, pix)
		}
	}
}

func testTextureAt(pix []byte) {
	count := 1000
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

// draw a frame
func drawFrame(surf *sdl.Surface, blank *sdl.Surface) {
	//clear screen
	rect := sdl.Rect{0, 0, width, height}
	blank.Blit(&rect, surf, &rect)
	//reset zbuffer
	for i := range zbuff {
		zbuff[i] = -1000
	}
	//get pixels
	surf.Lock()
	pix := surf.Pixels()
	//draw
	for _, f := range faces {
		if wireframe {
			wireframeFace(&f, pix)
		}
		drawFace(&f, pix)
	}
	//done drawing, set up window for display
	surf.Unlock()
	window.UpdateSurface()
}

func wireframeFace(f *Face, pix []byte) {
	//get vertices
	vs := f.V()
	v1, v2, v3 := vs[0], vs[1], vs[2]
	globalcolor = RED
	vline(v1, v2, pix)
	vline(v2, v3, pix)
	vline(v3, v1, pix)
}

// draw pixels where mask=1, for pixels contained by b
func fillMask(b box, pix []byte) {
	var i, j int
	for i = b.x0; i <= b.x1; i++ {
		for j = b.y0; j <= b.y1; j++ {
			midx := i + width*j
			if zmask[midx] == 0 {
				continue
			}
			putpixel(i, j, GREEN, pix)
		}
	}
}

// draw a face onto surface pixels
func drawFace(f *Face, pix []byte) {
	// get vertices for face
	v := f.V()
	a, b, c := v[0], v[1], v[2]
	//select pixels to be considered based on face bounds
	pbox := pixelbox(a, b, c)
	// populate zmask buffer
	//f.getMask(pbox)
	//fill solid color in where mask != 0
	n := f.Norm()
	fillMask(pbox, pix)
	for i := pbox.x0; i <= pbox.x1; i++ {
		var hit bool
		for j := pbox.y0; j <= pbox.y1; j++ {
			_, z, err := f.Project(i, j)
			//only draw pixel inside valid triangle face
			if err == offTriangleError && hit {
				break
			}
			if err != nil {
				continue
			}
			//only draw pixels for faces at front of screen
			if zbuff[i+j*width] > z {
				continue
			}
			zbuff[i+j*width] = z
			c := f.TexAt(i, j)
			if lightingEnabled {
				cmag := dot(n, F3{1, 0, 1})
				c = bright(c, cmag)
			}
			hit = true
			putpixel(i, j, c, pix)
		}
	}
	f.Unload()
	//TODO: update this
	return
	// calculate face normal for lighting
	var lightConts uint32
	for lidx, src := range lightpos {
		pow := lightpower[lidx]
		intensity := greyscale(dot(n, src) * pow)
		col := lightcolors[lidx]
		if !colorEnabled {
			col = RED | GREEN | BLUE
		}
		lightConts = intensity & col
		// TODO: put pixel drawing code here, so more than one light can be drawn
	}

	if !shadingEnabled {
		return
	}
	for i := pbox.x0; i <= pbox.x1; i++ {
		for j := pbox.y0; j <= pbox.y1; j++ {
			// putpixel(i, j, greyscale(math.Abs(vn1[2])), pix)
			maskval := zmask[i+width*j]
			if maskval == 0 {
				continue
			}
			if textureEnabled {
				// estimate texture color by using color at vertex 1.
				// TODO: replace this with extrapolation using barycentric cs.
				// texturecolor := textureFor(aidx)
				texturecolor := interpTexture(f, i, j)
				lightConts = texturecolor
				// v0col :=
			}
			lightConts = lightConts & maskval
			// lightConts = (RED | GREEN | BLUE) & maskval
			putpixel(i, j, lightConts, pix)
		}
	}
}

// build a zmask for face based on existing zbuffer, restricted to pixels in b
func (f *Face) getMask(b box) {
	// get z-pixel values
	var i, j int
	for i = b.x0; i <= b.x1; i++ {
		for j = b.y0; j <= b.y1; j++ {
			_, z, err := f.Project(i, j)
			// get masking index
			midx := i + j*width
			// skip pixels which are not inside of this face
			if err != nil {
				zmask[midx] = 0x0
				continue
			}
			// skip pixels where another triangle is in front of this face
			if zbuff[midx] > z {
				zmask[midx] = 0x0
				continue
			}
			//update mask, zbuffer to show position where this triangle is currently in front
			zmask[midx] = 0x1
			zbuff[midx] = z
		}
	}
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

// adjust color brightness
func bright(c uint32, i float64) uint32 {
	if i < 0 { // {{{
		i = 0
	}
	if i > 1 {
		i = 1
	}
	rc := c & RED
	gc := c & GREEN
	bc := c & BLUE
	r := uint32(i*float64(rc)) & RED
	g := uint32(i*float64(gc)) & GREEN
	b := uint32(i*float64(bc)) & BLUE
	return r | b | g // }}}
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

// put pixel with color in pixels array at pos x,y
func putpixel(x, y int, color uint32, pixels []byte) {
	if x >= width || y >= height || x < 0 || y < 0 {
		// fmt.Fprintf(os.Stderr, "Out of bounds putpixel: %v,%v\n", x, y)
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
