package main

import (
	"fmt"
	"regexp"
)

const MAXMATCH = 100 //control maximum matches per parse

func main() {
	_ = NewParseNd("f")
	ExampleParseNd_Parse()
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

// ex: Ob FileString -> Obg Function [3obs]-> Obg [fname] [1ob]
func (q *ParseNd) Parse(pattern string) ParseG {
	p := regexp.MustCompile(pattern)
	arr := p.FindAll(q.val, MAXMATCH)
	return slcgrp(arr, q)
}

// convert arr to ParseG and add to q
func slcgrp(arr [][]byte, q *ParseNd) ParseG {
	var g ParseG
	for _, s := range arr {
		newob := &ParseNd{val: s, p: make(map[string]ParseG)}
		g = append(g, newob)
	}
	name := fmt.Sprintf("%v", len(q.p))
	q.p[name] = g
	return g
}

func (g ParseG) Parse(name, pattern string) ParseG {
	var newg ParseG
	for _, v := range g {
		p := v.Parse(pattern)
		newg = append(newg, p...)
	}
	return newg
}
