package main

import (
	"fmt"
	"log"

	"github.com/veandco/go-sdl2/sdl"
)

var window *sdl.Window
var surf *sdl.Surface

const (
	width, height = 800, 800 //window dims
)

func main() {
	var winTitle string = "TinyRenderer"
	var err error
	blank, err := sdl.CreateRGBSurface(5, 800, 800, 32, 0, 0, 0, 0)
	if err != nil {
		log.Fatal(err)
	}
	rect := sdl.Rect{0, 0, width, height}
	blank.Blit(&rect, surf, &rect)
	window, err = sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(width)+1, int32(height)+1, sdl.WINDOW_SHOWN|sdl.WINDOW_ALWAYS_ON_TOP)
	if err != nil {
		log.Fatal(err)
	}
	defer window.Destroy()
	surf, err = window.GetSurface()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("vim-go")
	for {
		run()
		sdl.Delay(10)
	}
}

func run() {
	//set up window
	//get drawing surface from window
	err := surf.Lock()
	if err != nil {
		log.Fatal(err)
	}
	pix := surf.Pixels()
	//draw
	//      for _, v := range ob.vs {
	//              //fmt.Printf("putpixel @ %v %v (%v)\n", int(v.x), int(v.y), v)
	//              putpixel(int(v.x/2), int(v.y/2), 0xffffff00, pix)
	//      }
	pix[10] = 0xff
	//done drawing, set up window for display
	surf.Unlock()
	err = window.UpdateSurface()
	if err != nil {
		log.Fatal(err)
	}
}
