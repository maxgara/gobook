package main

import (
	"bufio"
	"image"
	"image/color"
	"image/gif"
	"io"
	"os"
	"text/template"
	"time"
)

const MODULEPATH = "maxgara-code.com"

func main() {

}
func git_notes() {
	//git push -u origin main
	//git status
	//git commit -m "message"
}
func notes() {
	//TIME
	start := time.Now()
	secs := time.Since(start).Seconds()
	_ = secs
	//FILE IO
	//rune by rune
	f, _ := os.Open("The Go Programming Language.txt")
	buf := bufio.NewReader(f)
	r, _, err := buf.ReadRune()
	_ = r
	//can do EOF check with:
	if err == io.EOF {
	}

	//print formatting
	//%q for quoted unicode, lets you see what the unprintable chars were (ex: '\n')
	//%v for general value (prints using object's method)

	//GIFs
	var palette = []color.Color{color.White, color.Black}
	const blackIndex = 1
	anim := gif.GIF{LoopCount: 100}
	rect := image.Rect(0, 0, 1, 1)          //create the image rectangle as a template for real image
	img := image.NewPaletted(rect, palette) //initialize image struct w/ palette
	img.SetColorIndex(1, 1, blackIndex)     // set a pixel to a color index from the image's palette
	anim.Delay = append(anim.Delay, 10)     // add a delay to the animation stack
	anim.Image = append(anim.Image, img)    //add an image to the animation stack
	gif.EncodeAll(os.Stdout, &anim)         //send gif to output
}

// HTML Templates
type ezwriter string

func (w ezwriter) Write(b []byte) (int, error) {
	w += ezwriter(b)
	return len(b), nil
}

func quickTempHTML(t string, data any) string {
	var w ezwriter
	template.Must(template.New("quickTemp").Parse(t)).Execute(w, data)
	return string(w)
}

const templ = `<h1>{{.TotalCount}} issues</h1>
<table>
<tr style='text-align: left'>
<th>#</th>
<th>State</th>
<th>User</th>
<th>Title</th>
</tr>
{{range .Items}}
<tr>
<td><a href='{{.HTMLURL}}'>{{.Number}}</td>
<td>{{.State}}</td>
<td><a href='{{.User.HTMLURL}}'>{{.User.Login}}</a></td>
<td><a href='{{.HTMLURL}}'>{{.Title}}</a></td>
</tr>
{{end}}
</table>
`
