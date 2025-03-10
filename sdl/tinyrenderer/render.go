package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/dblezek/tga"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	width, height = 800, 800 //window dims
	filename      = "african_head.obj"
	//filename = "square.obj"
	//filename        = "tri.obj"
	texturefilename = "african_head_diffuse.tga"
	delay           = 10
	yrotd           = 0.010 // +y azis rotation per frame
	xrotset         = 0     // +x azis rotation
	RED             = 0x0000ff00
	GREEN           = 0x00ff0000
	BLUE            = 0xff000000
	ALPHA           = 0x000000ff
)

type F3 [3]float64
type Face struct {
	vidx [3]int //vertex indices (in vertices global array)
	tidx [3]int //texture-vertex indices (in textureVertices global array)
}

// get vertices for face
func (f *Face) V() [3]F3 {
	var v [3]F3
	for i := range 3 {
		idx := f.vidx[i]
		v[i] = verts[idx-1]
	}
	return v
}

// get texture vertices for face
func (f *Face) T() [3]F3 {
	var t [3]F3
	for i := range 3 {
		idx := f.tidx[i]
		t[i] = textureVerts[idx-1]
	}
	return t
}

var window *sdl.Window

var wireframe bool
var file *os.File
var verts []F3
var faces []Face
var texFaces [][3]int
var textureVerts []F3
var texture []uint32
var texw, texh int

// var txxmin, txxmax, txymin, txymax float64
//var tstride int

// var fileFaceNorms []F3 // normal vector for each face (normalized to 1)
var done bool //control graceful program exit
var loops uint64

// var greyval float64
var zbuff []float64 //not currently used?
var zmask []uint32  //0xffffff00 where triangle is visible, otherwise 0x0
var lightpos []F3
var lightcolors []uint32
var lightpower []float64
var lightrot float64
var colorEnabled bool
var shadingEnabled bool
var cpuprofile bool
var start time.Time
var parallel int
var zmaskp [][]uint32 //??? only used for parallel case
var textureEnabled bool

func update() {
	for i := range verts {
		v := verts[i]
		v = yrot(v, yrotd)
		for j := range lightpos {
			lightpos[j] = xrot(lightpos[j], lightrot)
		}
		verts[i] = v
	}
	//greyval -= 0.001
	//this is a test
}
func benchStart() {
	start = time.Now()
	go func() {
		// terminate process after 120 seconds and report loops
		<-time.After(time.Second * 120)
		end()
	}()
}
func mainLoop() {
	//setup window
	var winTitle string = "TinyRenderer"
	var err error
	window, err = sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(width), int32(height), sdl.WINDOW_SHOWN|sdl.WINDOW_ALWAYS_ON_TOP)
	if err != nil {
		log.Fatal(err)
	}
	//get drawing surface from window
	surf, err := window.GetSurface()
	if err != nil {
		log.Fatal(err)
	}
	//create reusable blank surface to blit before redrawing
	blanksurf, err := sdl.CreateRGBSurface(5, 800, 800, 32, 0, 0, 0, 0)
	if err != nil {
		log.Fatal(err)
	}

	defer window.Destroy()

	//draw loop
	for {
		loops++
		draw(surf, blanksurf)
		takeKeyboardInput()
		if done {
			end()
			return
		}
		//sdl.Delay(delay)
	}
}
func testfunctions() {
	//test zpixel
	//zpixeldebug = true
	v1 := F3{-1, 0, 0}
	v2 := F3{1, 1, 0}
	v3 := F3{1, 0, 1}
	_ = vscale(v1, 1)
	_ = vadd(v1, v1)
	_ = vavg(v1, v1)
	_, _ = zpixel(v1, v2, v3, [2]int{3 * width / 4, 4 * height / 6})
	_, _ = zpixel(v2, v1, v3, [2]int{3 * width / 4, 4 * height / 6})
	zval, err := zpixel(v3, v2, v1, [2]int{3 * width / 4, 4 * height / 6})

	fmt.Printf("zval: v1=%v, v2=%v, v3=%v\tz=%v\terr=%v\n", v1, v2, v3, zval, err)
	//test cross
	norm := cross(v2, v3)
	//	fmt.Printf("cross of %v %v = %v", v2, v3, norm)
	norm = vnormalize(norm)
	_ = norm
	//	fmt.Printf("after normalization: %v\n", norm)
	//test dynamicNormalForFace
	dn := DynamicNormalForFace(v1, v2, v3)
	_ = dn
	//	fmt.Printf("normal for triangle %v %v %v = %v\n", v1, v2, v3, dn)
	//v3 = F3{1, 0, 0}
	dn = DynamicNormalForFace(v1, v2, v3)
	fmt.Printf("normal for triangle %v %v %v = %v\n", v1, v2, v3, dn)
	//av := vavg(v1, v2, v3)
	//fmt.Printf("vavg = %v\n", av)
	//test zpixelmask
	zbuff = make([]float64, width*height+801)
	for i := range zbuff {
		zbuff[i] = -1000
	}
	_ = zpixelboxmask(v1, v2, v3, zmask)
	for i := 0; i < height-1; i++ {
		//fmt.Printf("zmask:%v\n", zmask[i*width:i*width+width])
	}

}
func testDrawTextureImg(pix []byte) {
	f, err := os.Open("african_head_diffuse.tga")
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
			putpixel(i, j, color, pix)
		}
	}
}
func main() {
	textureEnabled = true
	parallel = 0 //false
	shadingEnabled = true
	cpuprofile = false
	if cpuprofile {
		f, err := os.Create("cpuprofile")
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	zbuff = make([]float64, width*height) //keep track of what's in front of scene

	//set up zmask (reusable) and 1 light
	zmask = make([]uint32, width*height+801)
	lightpos = append(lightpos, F3{2, 2.5, 20})
	//lightpos = append(lightpos, F3{-2, 1.5, 2})
	//lightpos = append(lightpos, F3{-2, -0.5, 2})
	lightcolors = append(lightcolors, RED|GREEN|BLUE)
	//lightcolors = append(lightcolors, GREEN|BLUE)
	//lightcolors = append(lightcolors, RED)
	lightpower = append(lightpower, 1)
	//lightpower = append(lightpower, 1)
	//lightpower = append(lightpower, 0.5)

	loadobjfile(filename)
	texture = loadTexture(texturefilename)
	for i, v := range verts {
		verts[i] = xrot(v, xrotset)
	}

	//testfunctions()
	mainLoop()
}
func takeKeyboardInput() {
	if event := sdl.PollEvent(); event != nil {
		if event, ok := event.(*sdl.KeyboardEvent); ok {
			if event.State == 0 {
				return
			}
			fmt.Printf("event.Keysym.Scancode: %v %[1]c\n", event.Keysym.Scancode)
			switch event.Keysym.Scancode {
			case 20: //q
				done = true
			case 82: //up arrow
			case 81: //down arrow
			case 80: //left arrow
			case 79: //right arrow
			case 26: //'w'
				wireframe = !wireframe
			case 6: //'c'
				colorEnabled = !colorEnabled
			case 22: //'s'
				shadingEnabled = !shadingEnabled
			}

		}
		if _, ok := event.(*sdl.QuitEvent); ok {
			done = true
			return
		}
	}
}
func end() {
	dur := time.Since(start)
	var loopsPerSec = float64(loops) / float64(dur.Seconds())
	// fmt.Println(surf.Pixels())
	runtime.GC() // get up-to-date statistics
	file, _ = os.Create("memprofile")
	pprof.WriteHeapProfile(file)
	fmt.Printf("dur=%v; loops=%v; lps=%v\n", dur.Seconds(), loops, loopsPerSec)
}
