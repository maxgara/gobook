package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestColors(t *testing.T) {
	var cols = colorset{}
	// var newcols []color
	// newcols = append(newcols, cols.newcolor())
	fmt.Println(cols.newcolor())
	fmt.Println(cols.newcolor())
	fmt.Println(cols.newcolor())
	fmt.Println(cols.newcolor())
	fmt.Println(cols.newcolor())
	fmt.Println(cols.newcolor())
	fmt.Println(cols.newcolor())
	fmt.Println(cols.newcolor())
	fmt.Println(cols.newcolor())
	fmt.Println(cols.newcolor())
	fmt.Println(cols.newcolor())
	fmt.Println(cols.newcolor())
	fmt.Println(cols.newcolor())
	fmt.Println(cols.newcolor())
	fmt.Println(cols.newcolor())
	fmt.Println(cols.newcolor())
}
func TestSvgGen(t *testing.T) {
	const data = `1 5
	2 10
	3 30
	4 60`
	fmt.Sprintf(printSVGs(data))
}
func TestMultiColorSVG(t *testing.T) {
	s := ""
	for _, d := range data {
		s += fmt.Sprintf("%s\n\n", d)
	}
	printSVGs(s)
}
func TestPrintBoundsSvg(t *testing.T) {
	fmt.Printf("bounds:%s\n", printBoundsSVG([]float64{0, 0, 40, 50}))
}
func TestParsep(t *testing.T) {
	for _, s := range strings.Split(data[0], "\n") {
		p, n := parsep(s)
		fmt.Printf("p:%v n:%v\n", p, n)
	}
}
func TestParse(t *testing.T) {
	boxes := parse(data[0])
	fmt.Printf("boxes:%v\n", boxes)
}
func TestPrint(t *testing.T) {
	boxes := parse(data[0])
	print(boxes)
}

var data = []string{`1 5
2 12
3 25
4 60
5 80
6 45
7 30
8 90
9 70
10 55`,

	`1 8
2 15
3 20
4 50
5 60
6 35
7 28
8 85
9 65
10 50`,

	`1 10
2 18
3 27
4 62
5 77
6 40
7 36
8 95
9 72
10 58`,

	`1 7
2 13
3 24
4 55
5 82
6 47
7 34
8 92
9 68
10 52`,

	`1 6
2 14
3 22
4 57
5 75
6 42
7 32
8 88
9 66
10 54`,

	`1 9
2 17
3 26
4 61
5 79
6 44
7 29
8 93
9 73
10 53`,

	`1 4
2 11
3 23
4 58
5 78
6 41
7 31
8 89
9 69
10 56`,

	`1 3
2 16
3 21
4 63
5 76
6 39
7 33
8 87
9 67
10 51`,
}
