package main

import (
	"log"
	"os"
	"strconv"
	"strings"
)

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
func loadface(s string, faces *[][3]int) {
	fs := strings.Fields(s)
	var f [3]int
	for i, field := range fs[1:] {
		idxstr := strings.Split(field, "/")[0]
		fidx, err := strconv.ParseInt(idxstr, 10, 32)
		if err != nil {
			log.Fatal(err)
		}
		f[i] = int(fidx)
	}
	*faces = append(*faces, f)
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
			loadVertex(v, &fileVerts)
		case strings.HasPrefix(v, "f "):
			loadface(v, &fileFaces)
		case strings.HasPrefix(v, "vt "):
			loadvtex(v, &textureVerts)
		}
	}
}
