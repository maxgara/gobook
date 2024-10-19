package mapstack

import (
	"fmt"
	"runtime"
	"strings"
)

const STACKDEPTH = 20

// a frame becomes a vertex in the graph of function calls
type Vertex struct {
	name    string
	pack    string
	file    string
	calls   map[*Vertex]bool //child nodes
	callers map[*Vertex]bool //parent nodes
	flags   int              //flags representing the kind of caller (direct, indirect, active, inactive, etc.)
}

func getcFrames() (f *runtime.Frames) {
	var pcs = make([]uintptr, STACKDEPTH) //program counters for calling funcs
	n := runtime.Callers(2, pcs)
	pcs = pcs[:n] //remove extra buffer space
	return runtime.CallersFrames(pcs)
}

// break a function call like github.com/x/mutextrace.getlockcallers into package path and func name values
func parseFunc(fun string) (p string, f string) {
	pathend := strings.LastIndex(fun, "/")
	if pathend == -1 {
		pathend = 0
	}
	dot := strings.Index(fun[pathend:], ".") + pathend
	if dot != -1 {
		p = fun[:dot]
		f = fun[dot+1:]
	}
	return p, f
}
func run() {
	vlookup := make(map[string]*Vertex) //collection of all frames and their addresses
	addstack(vlookup)
}

// adds the most direct caller from the top of the stack to the graph, linked to all previous vertices
func addstack(vlookup map[string]*Vertex) {
	frames := getcFrames()
	var prev *Vertex //previous frame, when applicable
	for {
		f, more := frames.Next()
		zero := runtime.Frame{}
		//should not happen
		if f == zero {
			panic("no frames in stack trace?")
		}
		v := getVertex(&f)
		// link node in previous iteration, if applicable
		if prev != nil {
			v.calls[prev] = true
			prev.callers[v] = true
		}
		prev = v
		//if v represents a vertex (p) already in graph,  combine them. Else, add p
		if p, ok := vlookup[v.pack+"/"+v.name]; ok {
			merge(p, v)
		} else {
			vlookup[v.pack+"/"+v.name] = v
		}
		if !more {
			break
		}
	}
}

// verge vertex v into p, so v can be deleted
func merge(v *Vertex, p *Vertex) {
	for s, _ := range v.callers {
		p.callers[s] = true
	}
	for s, _ := range v.calls {
		p.calls[s] = true
	}
}

func getVertex(f *runtime.Frame) *Vertex {
	p, fname := parseFunc(f.Function)
	return &Vertex{fname, p, f.File, make(map[*Vertex]bool), make(map[*Vertex]bool), 0}
}
func printmap(l map[string]*Vertex) {
	for s, p := range l {
		fmt.Printf("%v:\ncallers:%v\ncalls:%v\n", s, p.callers, p.calls)
	}
}
