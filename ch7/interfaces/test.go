package main

import (
	"fmt"
	"io"

	"maxgara-code.com/workspace/ch7/stringreader"
)

func main() {
	ptest()

}

//	func f(out io.Writer) {
//		if out == nil {
//			fmt.Printf("out == nil\n")
//		}
//		out.Write([]byte{1, 2, 3})
//	}
func ptest() {
	var s = make([]io.Reader, 0, 5)
	for i := 0; i < 3; i++ {
		s = append(s, stringreader.Newreader("test"))
	}
	t := trivialReader{}
	s = append(s, t)
	s = append(s, stringreader.Newreader("test"))

	print(&s)
}

func print(s *[]io.Reader) {
	var maxlen int
	for _, v := range *s {
		vl, pl := len(fmt.Sprint(v)), len(fmt.Sprintf(&v))
		if len(v) > flen {
			flen = len(v)
		}
	}
	for _, v := range *s {
		fmt.Printf("%*p\t", &v)
	}
	fmt.Println()
	for _, v := range *s {
		fmt.Printf("%*16v\t", v)
	}
}

type trivialReader struct {
	A float64
	B float64
	C float64
}

func (t trivialReader) Read([]byte) (int, error) {
	return 0, io.EOF
}
