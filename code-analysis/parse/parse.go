package main

import (
	"fmt"
	"regexp"
	"strings"
)

//Next Up: add functionality to name a property in the group's regex expression.

const MAXMATCH = 100 //control maximum matches per parse

func main() {
	s := "abc abd abx bay pab"
	x := NewParseNd(s)
	words := x.NamedParse("abc (?<myname>..)")
	fmt.Print(words)

}

// object group
type ParseG []*ParseNd

// object
type ParseNd struct {
	val []byte            //substring ref. should eventually use rune?
	p   map[string]ParseG //props; child node group
	anc *ParseNd          //direct ancestor
}

func NewParseNd(val string) *ParseNd {
	return &ParseNd{val: []byte(val), p: make(map[string]ParseG)}
}
func (q *ParseNd) String() string {
	var pstr string
	for k := range q.p {
		pstr += fmt.Sprintf("\t%v\n", k)
	}
	return fmt.Sprintf("node:%s\n\tp:\n%v\tanc:%p\n", q.val, pstr, q.anc)
}

// func (g *ParseG) String() string {
// 	out := fmt.Sprintln("Group:")
// 	for _, q := range *g {
// 		out += fmt.Sprint(q)
// 	}
// 	return out
// }

func (q *ParseNd) Walk(f func(q *ParseNd) bool) {
	if !f(q) {
		return
	}
	for _, g := range q.p {
		for _, next := range g {
			next.Walk(f)
		}
	}
}

// name the properties you want parsed out using <?name> syntax in pattern. Uses previous Parse function to work
// supports multiple property extractions per call, probably could have simplified this a little by only supporting one
// I have just realized in the course of debugging that golang's regexp lib already supports named captures...  :(
func (q *ParseNd) NamedParse(pattern string) ParseG {
	pnode := NewParseNd(pattern)
	caps := pnode.Parse(`\(\?<\w+>[^\)]+\)`) //parse named capture from full pattern
	fmt.Println("named capture groups found:")
	fmt.Println(caps)
	namenodes := caps.Parse(`\w+`) //parse name from named capture
	fmt.Println("names from named capture groups:")
	fmt.Println(namenodes)
	var names []string
	for _, v := range caps {
		idx := strings.Index(pattern, string(v.val)) // find named capture idx in original pattern
		namenode := v.p["0"][0]                      //get name node attached named capture node
		namel := len(namenode.val)
		beforen := pattern[:1+idx] // pattern until ?
		fmt.Printf("beforen: '%v'\n", beforen)
		name := pattern[3+idx : namel+idx] // name
		fmt.Printf("name: '%v'\n", name)
		var aftern = "" //pattern after > ,Default val empty before bounds check
		if len(pattern) > len(v.val)+idx {
			aftern = pattern[len(v.val):] // pattern after n (should be checking for bounds here)******
		}
		pattern = beforen + aftern //modify pattern to exclude name part
		names = append(names, name)
	}
	fmt.Println("actual names, extracted:")
	fmt.Println(names)
	q.Parse(pattern) //return val
	var ret ParseG
	for i, name := range names {
		mapk := fmt.Sprint(i)
		q.p[name] = q.p[mapk]
		delete(q.p, mapk)
		ret = append(ret, q.p[name]...)
	}
	return ret
}

func (q *ParseNd) Parse(pattern string) ParseG {
	p := regexp.MustCompile(pattern)
	arr := p.FindAll(q.val, MAXMATCH)
	return slcgrp(arr, q)
}

// convert arr to ParseG and add to q
func slcgrp(arr [][]byte, q *ParseNd) ParseG {
	var g ParseG
	for _, s := range arr {
		newob := &ParseNd{val: s, p: make(map[string]ParseG), anc: q}
		g = append(g, newob)
	}
	name := fmt.Sprintf("%v", len(q.p))
	q.p[name] = g
	return g
}

func (g ParseG) Parse(pattern string) ParseG {
	var newg ParseG
	for _, v := range g {
		p := v.Parse(pattern)
		newg = append(newg, p...)
	}
	return newg
}
