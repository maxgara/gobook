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
	return texture, img.Bounds().Max.X - img.Bounds().Min.X
}

// get texture color at coordinates x,y.
func textureAt(x, y float64) uint32 {
	idx := x + y*float64(tstride)
	col := texture[int(idx)]
	return col
}

// get texture color for vertex idx
func textureFor(vidx int) uint32 {
	t := textureVerts[vidx]
	c := textureAt(t[0], t[1])
	return c
}
