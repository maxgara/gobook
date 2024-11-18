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

const teststr = `func f1(int a, int b) int{
    a = a+b
    return a*b
}`

func main(){
file := NewParseNd("f", teststr)

ftemp:= file.Temp("ft", "func.*")
fmt.Println(file)
fmt.Println(ftemp)
}

func NewParseNd(name, val string) *ParseNd{
 return &ParseNd{name:name, s:val, p:make(map[string]ParseG)}
}
func (q *ParseNd) Walk(f func(q *ParseNd)) {
	f(q)
	for _, g := range q.p {
		for _, next := range g {
			next.Walk(f)
		}
	}
}

// ex: Ob FileString -> Obg Function [3obs]-> Obg [fname] [1ob]
func (q *ParseNd) Parse(name, pattern string) ParseG {
	p := regexp.MustCompile(pattern)
	arr := p.FindAllString(q.s, MAXMATCH)
	var g ParseG
	for i, s := range arr {
		//keep Temp suffix at the end.
		var newname string
		if basename, temp := strings.CutSuffix(name, "Temp"); temp {
			newname = basename + fmt.Sprintf("%v", i) + "Temp"
		} else {
			newname = basename + fmt.Sprintf("%v", i)
		}
		newob := ParseNd{name: newname, s: s, p: make(map[string]ParseG)}
		g = append(g, newob)
	}
	q.p[name] = g
	return g
}
func (q ParseNd) String() string {
	pstr := "pEMPTY"
	if len(q.p) != 0 {
	    for k,v := range q.p{
		
		pstr = fmt.Sprintf("%v:%v\n", k,v)
		}
	}
	return fmt.Sprintf("%v:\"%v\"; props:%v", q.name, q.s, pstr)
}
func (q ParseG) String() string {
	var s string
	for _, v := range q {
		s += v.String() + "\n"
	}
	return s
}

// should just be named Parse
func (g ParseG) ParseEach(name, pattern string) ParseG {
	var newg ParseG
	for _, v := range g {
		p := v.Parse(name, pattern)
		newg = append(newg, p...)
	}
	return newg
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
	f := func(current *ParseNd) {
		current.rSave(newp, q)
	}
	q.Walk(f)
}

// non-public helper func
func (current *ParseNd) rSave(newp ParseG, anc *ParseNd) {
	for i := range newp {
		q := &newp[i]
		if current == q {
			//I should never have started doing this numbering thing
			nl := strings.IndexAny(q.name, "0123456789")
			name := q.name[0:nl]
			qcopy := *q
			qcopy.name= name
			anc.p[name] = append(anc.p[name], qcopy)
		}
	}
}
