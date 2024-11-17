package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

func main() {
	var data []byte
	data, _ = io.ReadAll(os.Stdin)
	s := string(data)
	fidxarr := findFunctions(s)
	farr := parseFunctions(s, fidxarr)

	tree(farr)
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
func findFunctions(s string) [][]int {
	// re := regexp.MustCompile(`func.*\(.*\).*\{.*}`)
	re := regexp.MustCompile(`func\s*[^\s\(]+.*?\(.*?\).*?{`)
	matchsets := re.FindAllStringIndex(s, -1)
	for i, match := range matchsets {
		cbrack := matchb(s, match[1]-1)
		if cbrack == -1 {
			panic("unmatched bracket")
		}
		matchsets[i][1] = cbrack
	}
	return matchsets
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

// extract funcs from file string, idxs returned by findFunctions
func parseFunctions(s string, idxs [][]int) []fn {
	var fns []fn
	for _, fidx := range idxs {
		fstr := s[fidx[0]:fidx[1]]
		fsplit := strings.Split(fstr, "\n")
		head := fsplit[0]
		bodyarr := fsplit[1:]
		//trim body
		for i, v := range bodyarr {
			bodyarr[i] = strings.TrimSpace(v)
		}
		body := strings.Join(bodyarr, "\n")
		body = body[:len(body)-1] //remove trailing }
		//extract
		headp := regexp.MustCompile(`func\s+([^(]+?)\s*\((.*)\)`) //file name, argstring pattern
		headms := headp.FindStringSubmatch(head)                  //matches
		name := headms[1]
		argstr := headms[2]
		argp := regexp.MustCompile(`(\w+)\b`) //arg pattern
		args := []string{}
		//extract args from argstring
		for _, v := range strings.Split(argstr, ",") {
			argm := argp.FindStringSubmatch(v) //argp matches.
			//account for 0 arg case
			if len(argm) == 0 {
				break
			}
			arg := argm[1] //0 would work too...
			args = append(args, arg)
			varsetidxs := findvarsets(body)
			vardecidxs := findvardecs(body)
			//TODO: convert above 2 vars to strings, add to f struct
			_ = varsetidxs
			_ = vardecidxs
			//TODO: write and call here func to find function calls
		}
		fns = append(fns, fn{name: name, args: args, body: body})
	}
	return fns
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

// match curly bracket at pos start in string s. return pair of ints so that s[x[0]:x[1]] contains both brackets
func matchb(s string, start int) int {
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
