package main

import (
	"log"
	"os"

	"image"
	"image/png"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	width, height = 800, 800 //window dims
	RED           = 0x0000ff00
	GREEN         = 0x00ff0000
	BLUE          = 0xff000000
	ALPHA         = 0x000000ff
	imagefile     = "z3.bin"
)

var done bool
var blank *sdl.Surface
var rect sdl.Rect

func testfunction() image.Image {
	//f, err := os.Open("rainbow.png")
	f, err := os.Open("/System/Library/CoreServices/Dock.app/Contents/Resources/finder.png")
	if err != nil {
		log.Fatal(err)
	}
	img, err := png.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	return img
}

// setup window and surface
func setup() (*sdl.Window, *sdl.Surface) {
	//create window
	var winTitle string = "Render"
	var err error
	window, err := sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(width), int32(height), sdl.WINDOW_SHOWN|sdl.WINDOW_ALWAYS_ON_TOP)
	if err != nil {
		log.Fatal(err)
	}

	//get drawing surface
	surf, err := window.GetSurface()
	if err != nil {
		log.Fatal(err)
	}

	//initialize new blank surface and framing rectangle for convenience of draw func
	rect = sdl.Rect{0, 0, width, height}
	blank, err = sdl.CreateRGBSurface(5, width, height, 32, 0, 0, 0, 0)
	if err != nil {
		log.Fatal(err)
	}
	return window, surf
}

func main() {
	//img := loadbin(imagefile)
	img := testfunction()
	window, surf := setup()
	defer window.Destroy()
	//done triggers graceful exit if user presses 'q' key
	for !done {
		draw(window, surf, img)
		takeKeyboardInput()
	}
}
func loadPNG(filename string) {
}

// load binary file into image
// file is assumed to have RGBG interpolated values with new rows beginning at intervals of 0x1000.
// trailing columns of 0x0 values are dropped
func loadbin(filename string) *image.RGBA {
	f, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	img := image.RGBA{}
	img.Pix = make([]uint8, len(f)*5)
	var rgMode bool //true=RG row, false=GB row
	for i := 0; i+1 < len(f); i += 2 {
		//load two bytes; either R,G or G,B
		b0, b1 := f[i], f[i+1]
		if i%0x1000 == 0 {
			rgMode = !rgMode
		}
		//bytes -> pixels
		var p0, p1 [3]byte
		if rgMode {
			p0 = [3]byte{b0, 0, 0}
			p1 = [3]byte{0, b1, 0}
		} else {
			p0 = [3]byte{0, b0, 0}
			p1 = [3]byte{0, 0, b1}
		}
		//set each rgba value in image.RGBA
		for j := range 3 {
			img.Pix[4*i+j] = p0[j]
			img.Pix[4*(i+1)+j] = p1[j]
		}
	}
	img.Stride = 0x1000
	img.Rect = image.Rect(0, 0, 0x400, len(img.Pix)/(0x4*0x1000)+1)
	return &img
}

// draw img on surf in window
func draw(window *sdl.Window, surf *sdl.Surface, img image.Image) {
	//clear surface
	blank.Blit(&rect, surf, &rect)
	//lock the surface for pixel editing
	surf.Lock()
	pix := surf.Pixels()
	for i := img.Bounds().Min.X; i < img.Bounds().Max.X; i++ {
		for j := img.Bounds().Min.Y; j < img.Bounds().Max.Y; j++ {
			x := i
			y := j
			r, g, b, a := img.At(x, y).RGBA()
			r = 255
			//convert 16-bit-per-channel values returned by interface function above to 8 bit-per-channel RGB
			//use BGRA format for SDL (since that is what the array of pixels uses)
			pix[i], pix[i+1], pix[i+2], pix[i+3] = uint8(b/0x100), uint8(g/0x100), uint8(r), uint8(a/0x100)
			//pix[i] = 255
		}
	}
	surf.Unlock()
	window.UpdateSurface()
}
func putpixel(x, y int, color uint32, pixels []byte) {
	if x >= width || y >= height || x < 0 || y < 0 {
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
	pixels[idx+3] = a
}

// currently only checks for quit signal
func takeKeyboardInput() {
	if event := sdl.PollEvent(); event != nil {
		if event, ok := event.(*sdl.KeyboardEvent); ok {
			switch event.Keysym.Scancode {
			case 20: //q
				done = true
			case 82: //up
			case 81: //down
			case 80: //left
			case 79: //right
			}

		}
		if _, ok := event.(*sdl.QuitEvent); ok {
			done = true
			return
		}
	}
}
