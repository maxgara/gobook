package main

import (
	"io"
	"os"
	"regexp"
)

//goal: read flow and create some sort of representation of the following:
// -possible paths
// -range of start and end states
// -dependence of end states on initial states
// -variable symmetry - transformations on var space which produce the same result f(x) = f(T(x)). ex T= x mod 3
// T(x1) = T(x2) implies f(x1) = f(x2)
// -code “symmetry?” how can code be systematically transformed (analagous to geometric rotations, reflections etc.) to yield the same result (representation)?
// -boundary conditions for paths
//
//current plan is to primarily use the idea of symmetry "groups" for a function f,
// a symmetry represents the range of a function, and the pattern in which

// beginning steps:

// write very basic go source parser:
// - identify function blocks
// - identify function calls
// - identify var declarations + definitions (without short var defs)
// - identify operators
// - identify if

// begin to associate operations with symmetries (no edge conditions,
// eg. x*x/x=x)

// create func call graph
// -only function names for now

// begin to compute symmetry flows/results, predict outcomes
// - operators z
//assumptions:
// no overflow/underflow
// no loops
// only ints

func main() {
	var data []byte
	data, _ = io.ReadAll(os.Stdin)
	s := string(data)
	findFunctions(s)
}

func findFunctions(s string) [][]int {
	// re := regexp.MustCompile(`func.*\(.*\).*\{.*}`)
	re := regexp.MustCompile(`func\s*[^\s\(]+.*?\(.*?\).*?{`)
	matchsets := re.FindAllStringIndex(s, -1)
	// fmt.Println(matchsets)
	for i, match := range matchsets {
		cbrack := matchb(s, match[1]-1)
		if cbrack == -1 {
			panic("unmatched bracket")
		}
		matchsets[i][1] = cbrack
	}
	return matchsets
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
