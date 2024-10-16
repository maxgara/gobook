package stringreader

import (
	"io"
)

type Sreader struct{ Str string }

func Newreader(s string) *Sreader {
	return &Sreader{s}
}

// read from sreader
func (r *Sreader) Read(p []byte) (n int, err error) {
	n = copy([]byte(r.Str), p)
	if n < len(r.Str) {
		r.Str = r.Str[n:]
		return n, nil
	}
	return n, io.EOF
}

// func main() {
// 	r := newreader("<html></html>")
// 	html.Parse(r)
// }
