package main

// import (
// 	"fmt"
// 	"io"
// 	"os"
// 	"strings"
// 	"testing"
// )

// //	func TestFindFuncs(t *testing.T) {
// //		arr := findFunctions(fstr1)
// //		// want := []string{}
// //		// for _, v := range arr {
// //		// 	// fmt.Printf("func:##%v##\n", fstr1[v[0]:v[1]])
// //		// }
// //	}
// func TestFindFuncs2(t *testing.T) {
// 	s := get("./excode/ex1.go")
// 	arr := findFunctions(s)
// 	var wantstrs = []string{`func main() {
// 	var x int
// 	var y int
// 	var z int
// 	z = x / 2
// 	z = z * 2
// 	y = z - 3
// 	fmt.Println(y)
// }`, `func f1(x int) int {
// 	return 3 * x
// }`}

// 	for i, v := range arr {
// 		// fmt.Printf("##%v##\n", s[v[0]:v[1]])
// 		if wantstrs[i] != s[v[0]:v[1]] {
// 			t.Fail()
// 			t.Log(s[v[0]:v[1]])
// 			t.Log(wantstrs[i])
// 		}
// 	}
// }
// func TestFindVarDecs(t *testing.T) {
// 	s := get("./excode/ex1.go")
// 	arr := findvardecs(s)
// 	var wantstrs = []string{`var x int`, `var y int`, `var z int`}
// 	for i, v := range arr {
// 		if wantstrs[i] != s[v[0]:v[1]] {
// 			t.Fail()
// 			t.Log("got:\n" + "\"" + s[v[0]:v[1]] + "\"")
// 			t.Log("want:\n" + "\"" + wantstrs[i] + "\"")
// 		}
// 	}
// }
// func TestParseFunctions(t *testing.T) {
// 	s := get("./excode/ex1.go")
// 	f := parseFunctions(s, findFunctions(s))
// 	fmt.Print(f)
// }
// func Examplefindvardecs() {
// 	s := "var x int\n//comment"
// 	arr := findvardecs(s)
// 	fmt.Printf("%s", s[arr[0][0]:arr[0][1]])
// 	//Output: var x int
// }
// func TestFindVarSets(t *testing.T) {
// 	s := get("./excode/ex1.go")
// 	arr := findvarsets(s)
// 	want := []string{`z = x / 2`,
// 		`z = z * 2`,
// 		`y = z - 3`}
// 	if !check(s, arr, want) {
// 		t.Fail()
// 	}
// }
// func TestMatchb(t *testing.T) {
// 	start := strings.Index(fstr1, "{")
// 	m := matchb(fstr1, start)
// 	mstr := fstr1[start:m]
// 	// fmt.Printf("start idx: %v match idx: %v\n", start, m)
// 	// fmt.Printf("match:%v", mstr)
// 	if mstr[0] != '{' || mstr[len(mstr)-1] != '}' {
// 		t.Fail()
// 	}
// }

// func get(s string) string {
// 	f, err := os.Open(s)
// 	if err != nil {
// 		fmt.Printf("%v", err.Error())
// 		return ""
// 	}
// 	b, _ := io.ReadAll(f)
// 	s = string(b)
// 	return s
// }

// func check(s string, gotidx [][]int, want []string) bool {
// 	for i, j := range gotidx {
// 		ws := want[i]
// 		gs := s[j[0]:j[1]]
// 		if ws != gs {
// 			fmt.Printf("got:\n[%v]\n want:\n[%v]\n", gs, ws)
// 			return false
// 		}
// 	}
// 	return true
// }

// const fstr1 = `

// // initialize curve in box b
// func (b *svg) NewCurve() {
// 	cc := len(b.Curves)
// 	c := curve{Col: b.colors.new(), Label: label}
// 	b.Curves = append(b.Curves, c)
// }

// // initialize new svg
// func newsvg() svg {
// 	return svg{Curves: make([]curve, 0), Xmin: math.MaxFloat64, Ymin: math.MaxFloat64, Xmax: -math.MaxFloat64, Ymax: -math.MaxFloat64, colors: colorset{}}
// }

// // polyline curve data
// type curve struct {
// 	P     []point //points on curve
// 	Col   color   //color of line
// 	Label string  //label for curve
// }
// `
