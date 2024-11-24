package main

import (
	"fmt"
	"regexp"

	"maxgara-code.com/workspace/code-analysis/parse"
)

const test = `func f1(c int, d float) {
	var x int
	var y int
	var z int
	z = x / 2
	z = z * 2
	y = z - 3
	fmt.Println(y)
}
	func f2(a, b string){
	return
	}`

func main() {
	// var data []byte
	// data, _ = io.ReadAll(os.Stdin)
	file := parse.NewParseNd(test)
	funcs := findFunctions(file)
	_ = funcs
	// farr := parseFunctions(s, fidxarr)
}

// construct function call tree
func tree(farr []fn) {
	//convert to map for convenience
	fs := make(map[string]fn)
	for _, f := range farr {
		fs[f.name] = f
	}
	for fname, f := range fs {
		// TODO: make map
		_ = fname
		_ = f
	}
}

// find all function blocks and return start,end idxs. s[start:end] is full function
func findFunctions(q *parse.ParseNd) parse.ParseG {
	funcs := q.Parse(`(?<funcs>func [^\(]+(?:.*\n.*)*)`)
	// adjust funcs to stop at correct }
	sbrackets := funcs.Parse("(?<startbracket>{)")
	for _, v := range sbrackets {
		stop := matchb(*v.Base, v.Idx)
		v.Val = (*v.Base)[v.Idx:stop]
	}
	// names := funcs.Parse(`func (?<func_name>[^\(\s]+)`)
	// funcs.Parse(`[\(,]\s?(?<argnames>[^),\s]+)`)
	fmt.Println(funcs)
	return funcs
}

type fn struct {
	name     string   //name
	args     []string //args (all ints)
	body     string   //body text
	vsetstrs []string //var modification or initialization strings
	vdecstrs []string //var declaration strings
}

func (f fn) String() string {
	return fmt.Sprintf("func %v (args:%v)\ndecs:%v\nsets:\n%v\n", f.name, f.args, f.body, f.vdecstrs, f.vsetstrs)
}

func findvardecs(s string) [][]int {
	re := regexp.MustCompile(`var\s+[\w\d]+(?:,\s[\w\d]+)*\s(?:int|bool)`)
	matchsets := re.FindAllStringIndex(s, -1)
	return matchsets
}

func findvarsets(s string) [][]int {
	re := regexp.MustCompile(`[\w\d]+\s?=.*$`)
	matchsets := re.FindAllStringIndex(s, -1)
	return matchsets
}

// match curly bracket at pos start in byte slice s. return pair of ints so that s[x[0]:x[1]] contains both brackets
func matchb(s []byte, start int) int {
	var opens, closes int
	for i, c := range s[start:] {
		if c == '{' {
			opens++
		} else if c == '}' {
			closes++
		}
		if opens == closes {
			return i + start + 1
		}
	}
	return -1
}
