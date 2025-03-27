package main

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

// load normal vector from string
func loadNormal(s string, vns *[]V4) {
	fs := strings.Fields(s)
	var arr [3]float64
	for i, coordstr := range fs[1:] {
		coord, err := strconv.ParseFloat(coordstr, 64)
		if err != nil {
			log.Fatal(err)
		}
		arr[i] = coord
	}
	vn := V4{arr[0], arr[1], arr[2], 0} //Normal is a vector, not a point, so it has M=0 magic coordinate
	*vns = append(*vns, vn)
}

// load texture vertex from string
func loadTextureVertex(s string, tvs *[]V4) {
	fs := strings.Fields(s)
	var arr [3]float64
	for i, coordstr := range fs[1:] {
		coord, err := strconv.ParseFloat(coordstr, 64)
		if err != nil {
			log.Fatal(err)
		}
		arr[i] = coord
	}
	vt := V4{arr[0], arr[1], arr[2], 0}
	*tvs = append(*tvs, vt)
}

// load vertex from string
func loadVertex(s string, verts *[]V4) {
	fs := strings.Fields(s)
	var arr [3]float64
	for i, coordstr := range fs[1:] {
		coord, err := strconv.ParseFloat(coordstr, 64)
		if err != nil {
			log.Fatal(err)
		}
		arr[i] = coord
	}
	vt := V4{arr[0], arr[1], arr[2], 1}
	*verts = append(*verts, vt)
}

// load face from string
func loadface(s string, faces *[]Face) {
	fs := strings.Fields(s)
	var face Face
	var vidxs [3]int //face vert idxs
	var tidxs [3]int //face texture-vert idxs
	T0M = getT0M()
	for i, field := range fs[1:] {
		// no texture case
		if !strings.Contains(field, "/") {
			vidx, err := strconv.ParseInt(field, 10, 32)
			if err != nil {
				log.Fatal(err)
			}
			vidxs[i] = int(vidx)
			continue
		}

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
func loadobjfile(filename string) *Obj {
	var vs []V4
	var tvs []V4
	var vns []V4
	var fs []Face
	f, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	s := strings.Split(string(f), "\n")
	for _, v := range s {
		switch {
		case strings.HasPrefix(v, "v "):
			loadVertex(v, &vs)
		case strings.HasPrefix(v, "f "):
			loadface(v, &fs)
		case strings.HasPrefix(v, "vt "):
			loadTextureVertex(v, &tvs)
		case strings.HasPrefix(v, "vn "):
			loadNormal(v, &vns)
		}
	}
	//deep copy of vertex slice needed for fileVs, since they will need to remain the same while normal vs subject to transformations.
	//same for vertex-normals
	vsCopy := make([]V4, len(vs))
	vnsCopy := make([]V4, len(vs))
	copy(vsCopy, vs)
	copy(vnsCopy, vns)
	ob := Obj{vs: vs, tvs: tvs, vns: vns, fs: fs, fileVs: vsCopy, fileNs: vnsCopy}
	return &ob
}
func setup() *Obj {
	//set default values for global options
	cpuprofile = true
	//allocate zbuffer
	zbuff = make([]float64, width*height) //keep track of what's in front of scene
	//load files
	ob := loadobjfile(filename)
	//set up window
	var winTitle string = "TinyRenderer"
	var err error
	window, err = sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(width), int32(height), sdl.WINDOW_SHOWN|sdl.WINDOW_ALWAYS_ON_TOP)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("window created")
	//get drawing surface from window
	surf, err = window.GetSurface()
	if err != nil {
		log.Fatal(err)
	}
	//allocate reusable blank surface to blit before redrawing
	blanksurf, err = sdl.CreateRGBSurface(5, 800, 800, 32, 0, 0, 0, 0)
	if err != nil {
		log.Fatal(err)
	}
	//profiling
	if cpuprofile {
		f, err := os.Create("cpuprofile")
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
	}
	//benchmarking
	start = time.Now()
	return ob
}
