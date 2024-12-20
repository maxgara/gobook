package main

// import (
// 	"fmt"

// 	"maxgara-code.com/workspace/code-analysis/parse"
// )

// const test = `
// func f1(c int, d float) {
// 	var x int
// 	var y int
// 	var z int
// 	z = x / 2
// 	z = z * 2
// 	y = z - 3
// 	fmt.Println(y)
// }
// func f2(a, b string){
// 	return
// }`

// func main() {
// 	file := parse.NewParseNd(test)
// 	funcs := file.Parse(`(?<funcs>func \w[^\{]+\{)`)
// 	sbracks := funcs.Temp(`(?<funcsbracket>\{)`)
// 	file.Refresh()
// 	//adjust func ends
// 	for i, br := range sbracks {
// 		end := matchb(*br.Base, br.Idx) //get end idx
// 		funcs[i].Val = (*funcs[i].Base)[funcs[i].Idx:end]
// 	}
// 	vars := funcs.Temp(`\s+(?<var>(?:var .*)|(?:\w+\s\:=\s.*))`)
// 	_ = vars.Parse(`var (?<varname>\w+)`)
// 	vars.Parse(`(?<varname>\w+) := `)
// 	fmt.Println(file)
// 	// _ = funcs
// 	// // farr := parseFunctions(s, fidxarr)
// }

// // match curly bracket at pos start in byte slice s. return pair of ints so that s[x[0]:x[1]] contains both brackets
// func matchb(s []byte, start int) int {
// 	var opens, closes int
// 	for i, c := range s[start:] {
// 		if c == '{' {
// 			opens++
// 		} else if c == '}' {
// 			closes++
// 		}
// 		if opens == closes {
// 			return i + start + 1
// 		}
// 	}
// 	return -1
// }
