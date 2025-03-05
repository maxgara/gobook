package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/dblezek/tga"
)

//load texture, texture helper functions

// load texture vertex coordinates
func loadvtex(s string, tvs *[]F3) {
	fs := strings.Fields(s)
	var vt F3
	for i, coordstr := range fs[1:] {
		coord, err := strconv.ParseFloat(coordstr, 64)
		if err != nil {
			log.Fatal(err)
		}
		vt[i] = coord
	}
	*tvs = append(*tvs, vt)
}

// load texture file
func loadTexture(fs string) (texture []uint32, tstride int) {
	//f, err := os.Open("african_head_diffuse.tga")
	f, err := os.Open(fs)
	if err != nil {
		log.Fatal(err)
	}
	img, err := tga.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	for i := img.Bounds().Min.X; i < img.Bounds().Max.Y; i++ {
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
			a = a
			color := b | g | r | a
			texture = append(texture, color)
		}
	}
	texw = img.Bounds().Max.X - img.Bounds().Min.X
	texh = img.Bounds().Max.Y - img.Bounds().Min.Y
	return texture, img.Bounds().Max.X - img.Bounds().Min.X
}

// get texture color at coordinates x,y.
func textureAt(x, y float64) uint32 {
	xidx := (x + 1) / 2 * float64(texw)
	yidx := (y + 1) / 2 * float64(texh)
	idx := xidx + yidx*float64(tstride)
	//idx -= yidx * float64(tstride)
	if int(idx) >= len(texture) {
		panic("textureAt: texture coordinates out of bounds")
	}
	col := texture[int(idx)]
	return col
}

// get texture color for vertex idx
func textureFor(vidx int) uint32 {
	t := textureVerts[vidx]
	c := textureAt(t[0], t[1])
	return c
}

// stretch texture across triangle, map pixel (x,y) to color
func interpTexture(aidx, bidx, cidx, x, y int) uint32 {
	a, b, c := textureVerts[aidx], textureVerts[bidx], textureVerts[cidx]
	br, err := bary(a, b, c, x, y)
	if err != nil {
		//return 0
	}
	as := vscale(a, br[0])
	bs := vscale(b, br[1])
	cs := vscale(c, br[2])
	pos := vadd(as, bs, cs)
	return (textureAt(pos[0], pos[1]) & 0xffffff00)
}
