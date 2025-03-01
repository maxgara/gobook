package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

const (
	filename = "z3.bin"
)

func main() {
	f, _ := os.Create("out.png")
	rect := image.Rect(0, 0, 1024, 1024)
	img := image.NewRGBA(rect)
	for i := range 1024 * 1024 / 4 {
		green := color.RGBA{0, 255, 0, 0}
		fmt.Println(green)
		img.Pix[i*4] = 0xff
	}

	png.Encode(f, img)
	return

	bytes, err := os.ReadFile(filename)
	//bytes = bytes[:0x1000]
	if err != nil {
		log.Fatal(err)
	}
	//for i := range bytes {
	//	fmt.Printf("%v ", bytes[i])
	//}
	rect = image.Rect(0, 0, 0x1000, len(bytes)/0x1000)
	img = image.NewRGBA(rect)
	mode := 0
	for {
		switch mode {
		//RG row
		case 0:
			for i := 0; i < 0x1000 && i+1 < len(bytes); i += 2 {
				j := i / 0x1000
				//b1 := bytes[i]
				b2 := bytes[i+1]
				//red := color.RGBA{b1, 0, 0, 0}
				green := color.RGBA{0, b2, 0, 0}
				//img.SetRGBA(i, j, red)
				img.SetRGBA(i+1, j, color.RGBA{0, 255, 0, 0})
				img.SetRGBA(i+1, j, green)
			}
			mode = 1
		//GB row
		case 1:
			for i := 0; i < 0x1000 && i+1 < len(bytes); i += 2 {
				b1 := bytes[i]
				b2 := bytes[i+1]
				green := color.RGBA{0, b1, 0, 0}
				blue := color.RGBA{0, 0, b2, 0}
				j := i / 0x1000
				img.SetRGBA(i, j, green)
				img.SetRGBA(i+1, j, blue)
			}
			mode = 0
		}
		if len(bytes) <= 0x1000 {
			//all bytes read
			break
		}
		//drop processed row
		bytes = bytes[0x1000:]
	}
	png.Encode(os.Stdout, img)
}
