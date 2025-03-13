package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dblezek/tga"
)

//load texture, texture helper functions

// load texture vertex coordinates

// load texture file, output is an array of BGRA uint32 colors + the width of the image in pixels
func loadTexture(fs string) (texture []uint32) {
	//f, err := os.Open("african_head_diffuse.tga")
	f, err := os.Open(fs)
	if err != nil {
		log.Fatal(err)
	}
	img, err := tga.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	for i := img.Bounds().Min.X; i < img.Bounds().Max.X; i++ {
		for j := img.Bounds().Min.Y; j < img.Bounds().Max.Y; j++ {
			r, g, b, a := img.At(i, j).RGBA()
			//keep most significant bits of 16-bit color channels
			r = r >> 8
			g = g >> 8
			b = b >> 8
			a = a >> 8
			//put them in the right place for final uint32 color BGRA
			r = r << 8
			g = g << 16
			b = b << 24
			//a = a
			color := b | g | r | a
			texture = append(texture, color)
		}
	}
	texw = img.Bounds().Max.X - img.Bounds().Min.X
	texh = img.Bounds().Max.Y - img.Bounds().Min.Y
	return texture
}

// get texture color at coordinates x,y in space of [0,1) X [0,1).
// TODO: figure out why this is not working ( see testTextureAt )
// seems like the x,y are swapped AND the y is inverted?
func textureAt(x, y float64) uint32 {
	xidx := int(x * float64(texw-1))
	yidx := int(y * float64(texh-1))
	idx := int(xidx + yidx*texw)
	if idx < 0 || idx >= len(texture) {
		fmt.Printf("textureAt (%v %v) called\n", x, y)
		panic("textureAt: texture coordinates out of bounds")
	}
	col := texture[idx]
	return col
}

// get texture color for vertex idx
func textureFor(vidx int) uint32 {
	t := textureVerts[vidx]
	c := textureAt(t[0], t[1])
	return c
}

// gets texture color for pixel x,y based on texture coordinates of triangle ABC (where x,y should fall on ABC in screen space)
func interpTexture(f *Face, x, y int) uint32 {
	//get triangle in texture space
	tvs := f.T()
	vs := f.V()

	a, b, c := tvs[0], tvs[1], tvs[2]
	//get triangle in screen space
	va, vb, vc := vs[0], vs[1], vs[2]
	//fmt.Printf("textureVerts for %v %v %v are\n %v\n %v\n %v\n", aidx, bidx, cidx, a, b, c)
	//balance vertices in screen space to get (x,y)
	br, err := bary(va, vb, vc, x, y)
	if err != nil {
		return 0
	}
	//project new linear combination of vertices into texture space
	as := vscale(a, br[0])
	bs := vscale(b, br[1])
	cs := vscale(c, br[2])
	pos := vadd(as, bs, cs)
	//fmt.Printf("pixel %v %v falls within this triangle in screen-space, and has barycentric coordinates %v\n which when combined yield new texture-space point %v %v\n", x, y, br, pos[0], pos[1])
	return (textureAt(pos[0], pos[1]) & 0xffffff00)
}
