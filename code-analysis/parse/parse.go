package parse

import (
	"fmt"
	"regexp"
	"strings"
)

//Next Up: add functionality to name a property in the group's regex expression.

const MAXMATCH = 100 //control maximum matches per parse

// parse node group
type ParseG []*ParseNd

// parse node
type ParseNd struct {
	Val  []byte            //slice of base node
	p    map[string]ParseG //props; child node group
	anc  *ParseNd          //direct ancestor
	Base *[]byte
	temp bool
	Idx  int //relative to base
}

func NewParseNd(val string) *ParseNd {
	base := []byte(val)
	return &ParseNd{Val: base, p: make(map[string]ParseG), Idx: 0, Base: &base}
}
func (q *ParseNd) String() string {
	var pstr string
	for k, v := range q.p {
		pstr += fmt.Sprintf("\t%v\n%v\n", k, v)
	}
	return fmt.Sprintf("node:%s\noffset:%v\n\tp:\n%v", q.Val, q.Idx, pstr)
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

// parse first named subgroup as a property of q. if q is temp, then add properties to first non-temp
func (q *ParseNd) Parse(pattern string) ParseG {
	p := regexp.MustCompile(pattern)
	arr := p.FindAllSubmatchIndex(q.Val, MAXMATCH) // only first submatch is actually used
	name := pname(pattern)
	g := slcgrp(arr, q)
	//find non-temp ancestor of q (can be q itself)
	for q.temp {
		q = q.anc
	}
	q.p[name] = append(q.p[name], g...)
	return g
}

// temporarily parse property from q.
func (q *ParseNd) Temp(pattern string) ParseG {
	g := q.Parse(pattern)
	for _, t := range g {
		t.temp = true
	}
	return g
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
func slcgrp(arr [][]int, q *ParseNd) ParseG {
	var g ParseG
	for _, bounds := range arr {
		if len(bounds) <= 2 {
			continue //no submatch group (either user didn't include one or it didn't match)
		}
		val := q.Val[bounds[2]:bounds[3]] //bounds 2, 3 are the start, end idxs of first submatch
		newnode := ParseNd{Val: val, p: map[string]ParseG{}, anc: q, Base: q.Base, Idx: bounds[2] + q.Idx}
		g = append(g, &newnode)
	}
	return g
}

// Parse(pattern) for each member of g, unless g is a temporary group. If g is temporary, Parse as normal
// but assign properties to first non-temp ancestor.
// If pattern matches but no subgroup does then do nothing.
func (g ParseG) Parse(pattern string) ParseG {
	var newg ParseG
	for _, v := range g {
		p := v.Parse(pattern)
		newg = append(newg, p...)
	}
	return newg
}
func (g ParseG) Temp(pattern string) ParseG {
	var newg ParseG
	for _, v := range g {
		p := v.Temp(pattern)
		newg = append(newg, p...)
	}
	return newg
}
