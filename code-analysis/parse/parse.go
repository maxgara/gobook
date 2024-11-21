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
	x.Parse("ab(?<myname>.)").Parse("(?<xgroup>x)")
	fmt.Print(x)

}

// parse node group
type ParseG []*ParseNd

// parse node
type ParseNd struct {
	val  []byte            //slice of base node
	p    map[string]ParseG //props; child node group
	anc  *ParseNd          //direct ancestor
	base *[]byte
	idx  int //relative to base
}

func NewParseNd(val string) *ParseNd {
	base := []byte(val)
	return &ParseNd{val: base, p: make(map[string]ParseG), idx: 0, base: &base}
}
func (q *ParseNd) String() string {
	var pstr string
	for k, v := range q.p {
		pstr += fmt.Sprintf("\t%v\n%v\n", k, v)
	}
	return fmt.Sprintf("node:%s\noffset:%v\n\tp:\n%v", q.val, q.idx, pstr)
}
func (q ParseG) String() string {
	//indent member node strings
	var out string
	for _, v := range q {
		s := v.String()
		out += "\t{\n\t" + strings.ReplaceAll(s, "\n", "\n\t") + "}\n"
	}
	return out
}

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

// parse first named subgroup as a property of q.
func (q *ParseNd) Parse(pattern string) ParseG {
	p := regexp.MustCompile(pattern)
	arr := p.FindAllSubmatchIndex(q.val, MAXMATCH) // only first submatch is actually used
	name := pname(pattern)
	return slcgrp(arr, q, name)
}

// extract name string from pattern
func pname(p string) string {
	np := regexp.MustCompile(`\?<\w+>`) // name pattern
	nstr := np.FindString(p)            // name indexes
	if len(nstr) == 0 {
		return "" // no name found
	}
	nstr = nstr[2 : len(nstr)-1] //drop <? and >
	return nstr
}

// convert arr to ParseG and add to q as prop name
func slcgrp(arr [][]int, q *ParseNd, name string) ParseG {
	var g ParseG
	for _, bounds := range arr {
		if len(bounds) <= 2 {
			continue //no submatch group (either user didn't include one or it didn't match)
		}
		val := q.val[bounds[2]:bounds[3]] //bounds 2, 3 are the start, end idxs of first submatch
		newnode := ParseNd{val: val, p: map[string]ParseG{}, anc: q, base: q.base, idx: bounds[2] + q.idx}
		g = append(g, &newnode)
	}
	q.p[name] = g
	return g
}

// parse first named subgroup as a property each member q of g matching pattern.
// If pattern matches but no subgroup does then do nothing.
func (g ParseG) Parse(pattern string) ParseG {
	var newg ParseG
	for _, v := range g {
		p := v.Parse(pattern)
		newg = append(newg, p...)
	}
	return newg
}
