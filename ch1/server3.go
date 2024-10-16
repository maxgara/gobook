package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/lisa", l_handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)
	for k, v := range r.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	for k, v := range r.Form {
		fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
	}
}
func l_handler(w http.ResponseWriter, r *http.Request) {
	lissajous(w)
}

var palette = []color.Color{color.White, color.RGBA{0, 100, 0, 255}, color.RGBA{255, 0, 0, 255}}

// var palette = []color.Color{color.White, color.Black}

const (
	whiteIndex = 0
	blackIndex = 1
	magicIndex = 2
)

func lissajous(out io.Writer) {
	const (
		cycles  = 5
		res     = 0.001
		size    = 100
		nframes = 64
		delay   = 8
	)
	freq := rand.Float64() * 3 //relative freq of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 //phase diff
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t * freq * float64(phase))
			var index uint8
			if t < cycles*math.Pi {
				index = blackIndex
			} else {
				index = magicIndex
			}
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), index)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim)
}
