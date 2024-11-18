package main

import (
	"fmt"
	"regexp"
	"strings"
)

// object group
type ParseG []ParseNd

// object
type ParseNd struct {
	name string //if name matches Temp$, prop may eventually be purged. Don't add numbers
	s    string
	p    map[string]ParseG //props; child node group
}

const MAXMATCH = 100 //control maximum matches per parse

// ex: Ob FileString -> Obg Function [3obs]-> Obg [fname] [1ob]
func (q *ParseNd) Parse(name, pattern string) ParseG {
	p := regexp.MustCompile(pattern)
	arr := p.FindAllString(q.s, MAXMATCH)
	var g ParseG = ParseG{}
	for i, s := range arr {
		//keep Temp suffix at the end.
		var newname string
		if newname, temp := strings.CutSuffix(s, "Temp"); temp {
			newname += fmt.Sprintf("%v", i) + "Temp"
		} else {
			newname += fmt.Sprintf("%v", i)
		}
		newob := ParseNd{name: newname, s: s, p: make(map[string]ParseG)}
		g = append(g, newob)
	}
	q.p[name] = g
	return g
}

func (g ParseG) ParseEach(name, pattern string) ParseG {
	var all ParseG
	for _, v := range g {
		newobs := v.Parse(name, pattern)
		all = append(all, newobs...)
	}
	return all
}

// parse a temporary property
func (q *ParseNd) Temp(name, pattern string) ParseG {
	return q.Parse(name+"Temp", pattern)
}

// parse a temporary property
func (g ParseG) Temp(name, pattern string) ParseG {
	return g.ParseEach(name+"Temp", pattern)
}

// assign property p to ancestor g, minus elements of p not descended from o. remove Temp suffix if present.
func (q *ParseNd) Save(newp ParseG) {
	q.rSave(newp, q)
}

// non-public helper func
func (current *ParseNd) rSave(newp ParseG, anc *ParseNd) {
	for _, v := range newp {
		if current == &v {
			//remove old suffix
			cut := current.Temp("nodigits", `^\d+(Temp)?$`)
			newname := strings.TrimSuffix(current.name, cut[0].s) //base name, no suffix
			ls := fmt.Sprintf("%v", len(anc.p[newname]))          //get new numbering
			//assign prop using base name
			anc.p[newname] = append(anc.p[newname], v)
			//add new suffix to attached node
			newname += ls
			v.name = newname
			break
		}
	}
	for _, v := range current.p {
		for _, w := range v {
			w.rSave(newp, anc)
		}
	}
}

// Push each Ob in p to an ancestor in O
func (g ParseG) Push(p ParseG) []string {
	return nil
}
