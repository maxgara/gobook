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
	delay   = 200
	yrotd   = 0.01 // +y azis rotation per frame
	xrotset = 0    // +x azis rotation
	RED     = 0x0000ff00
	GREEN   = 0x00ff0000
	BLUE    = 0xff000000
	ALPHA   = 0x000000ff
)

type F3 [3]float64

var window *sdl.Window

var wireframe bool
var file *os.File
var fileVerts []F3
var fileFaces [][3]int
var textureVerts []F3
var texture []uint32
var tstride int

// var fileFaceNorms []F3 // normal vector for each face (normalized to 1)
var done bool //control program exit
var loops uint64

// var greyval float64
var zbuff []float64
var zmask []uint32
var lightpos []F3
var lightcolors []uint32
var lightpower []float64
var lightrot float64
var colorEnabled bool
var shadingEnabled bool
var cpuprofile bool
var start time.Time
var parallel int
var zmaskp [][]uint32
var textureEnabled bool

func update() {
	for i := range fileVerts {
		v := fileVerts[i]
		v = yrot(v, yrotd)
		for j := range lightpos {
			lightpos[j] = xrot(lightpos[j], lightrot)
		}
		fileVerts[i] = v
	}
	//greyval -= 0.001
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
	//create blank surface to blit before redrawing
	blanksurf, err := sdl.CreateRGBSurface(5, 800, 800, 32, 0, 0, 0, 0)
	if err != nil {
		log.Fatal(err)
	}

	defer window.Destroy()

	//benchmarking
	benchStart()

	//parallel zmask arrays so that masks do not change while they are being referenced
	if parallel > 1 {
		zmaskp = make([][]uint32, parallel)
		for i := range zmaskp {
			zmaskp[i] = make([]uint32, width*height+801)
		}
	}
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
	parallel = 0
	shadingEnabled = true
	cpuprofile = true
	if cpuprofile {
		f, err := os.Create("cpuprofile")
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	//zbuff = make([]float64, width*height)
	zmask = make([]uint32, width*height+801)
	lightpos = append(lightpos, F3{2, 2.5, 20})
	lightpos = append(lightpos, F3{-2, 1.5, 2})
	//lightpos = append(lightpos, F3{-2, -0.5, 2})
	//lightpos = append(lightpos, F3{0, 3.5, 1.5})
	//lightpos = append(lightpos, F3{0, -3.5, 1.5})
	lightcolors = append(lightcolors, RED|GREEN|BLUE)
	lightcolors = append(lightcolors, GREEN|BLUE)
	//lightcolors = append(lightcolors, RED)
	//lightcolors = append(lightcolors, RED|GREEN)
	lightpower = append(lightpower, 1)
	lightpower = append(lightpower, 1)
	//lightpower = append(lightpower, 0.5)
	//lightpower = append(lightpower, 0.5)
	loadobjfile(filename)
	texture, tstride = loadTexture("african_head_diffuse.tga")
	for i, v := range fileVerts {
		fileVerts[i] = xrot(v, xrotset)
		//		fmt.Printf("vtop(%v)=%v\n", v, vtop(v))
	}
	testfunctions()
	mainLoop()
}

// draw line between vertices
func drawFrame(pix []byte) {
	for _, f := range fileFaces {
		i1, i2, i3 := f[0], f[1], f[2]
		v1, v2, v3 := fileVerts[i1-1], fileVerts[i2-1], fileVerts[i3-1]
		//b := pixelbox(v1, v2, v3)
		globalcolor = RED
		if wireframe {
			vline(v1, v2, pix)
			vline(v2, v3, pix)
			vline(v3, v1, pix)
		}
		//vn1 := DynamicNormalForFace(v1, v2, v3)
		//vn0 := vavg(v1, v2, v3)
		//vline(vn0, vadd(vn0, vn1), pix)
		triangleBoxShader(i1-1, i2-1, i3-1, pix, zmask)
		//	DrawLine(b.x0, b.y0, b.x0, b.y1, GREEN|BLUE, pix)
		//	DrawLine(b.x0, b.y1, b.x1, b.y1, GREEN|BLUE, pix)
		//	DrawLine(b.x1, b.y1, b.x1, b.y0, GREEN|BLUE, pix)
		//	DrawLine(b.x1, b.y0, b.x0, b.y0, GREEN|BLUE, pix)
	}
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
			case 82: //up
			case 81: //down
			case 80: //left
			case 79: //right
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
