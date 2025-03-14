package main

import (
	"log"
	"os"
	"strconv"
	"strings"
)

// load texture vertex from string
func loadTextureVertex(s string, tvs *[]F3) {
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

// load vertex from string
func loadVertex(s string, verts *[]F3) {
	fs := strings.Fields(s)
	var vt F3
	for i, coordstr := range fs[1:] {
		coord, err := strconv.ParseFloat(coordstr, 64)
		if err != nil {
			log.Fatal(err)
		}
		vt[i] = coord
	}
	*verts = append(*verts, vt)
}

// load face from string
func loadface(s string, faces *[]Face) {
	fs := strings.Fields(s)
	var face Face
	var vidxs [3]int //face vert idxs
	var tidxs [3]int //face texture-vert idxs
	for i, field := range fs[1:] {
		vidxstr := strings.Split(field, "/")[0] //vertex index
		tidxstr := strings.Split(field, "/")[1] //texture vertex index
		vidx, err := strconv.ParseInt(vidxstr, 10, 32)
		if err != nil {
			log.Fatal(err)
		}
		tidx, err := strconv.ParseInt(tidxstr, 10, 32)
		if err != nil {
			log.Fatal(err)
		}
		vidxs[i] = int(vidx)
		tidxs[i] = int(tidx)
	}
	face.vidx = vidxs
	face.tidx = tidxs
	*faces = append(*faces, face)
}
func loadobjfile(filename string) {
	f, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	s := strings.Split(string(f), "\n")

	for _, v := range s {
		switch {
		case strings.HasPrefix(v, "v "):
			loadVertex(v, &verts)
		case strings.HasPrefix(v, "f "):
			loadface(v, &faces)
		case strings.HasPrefix(v, "vt "):
			loadTextureVertex(v, &textureVerts)
		}
	}
}
