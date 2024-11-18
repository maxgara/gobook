package main

import (
	"fmt"
	"regexp"
)

// object group
type Obs []Ob

// object
type Ob struct {
	name string //if name matches Temp$, prop may eventually be purged
	val  string
	p    map[string]Obs //props
}

const MAXMATCH = -1 //control maximum matches per parse

// ex: Ob FileString -> Obg Function [3obs]-> Obg [fname] [1ob]
func (o *Ob) Parse(name, pattern string) Obs {
	p := regexp.MustCompile(pattern)
	arr := p.FindAllString(o.val, MAXMATCH)
	var g Obs = Obs{}
	for i, s := range arr {
		newname := name + fmt.Sprintf("%v", i)
		newob := Ob{name: newname, val: s, p: make(map[string]Obs)}
		g = append(g, newob)
	}
	o.p[name] = g
	return g
}

func (o Obs) ParseEach(name, pattern string) Obs {
	var all Obs
	for _, v := range o {
		newobs := v.Parse(name, pattern)
		all = append(all, newobs...)
	}
	return all
}
func (o *Ob) Temp(name, pattern string) Obs {
	return o.Parse(name+"Temp", pattern)
}
func (o Obs) Temp(name, pattern string) Obs {
	return o.ParseEach(name+"Temp", pattern)
}

func (o Obs) Push(p Obs) []string {
	// TODO
	return nil
}
