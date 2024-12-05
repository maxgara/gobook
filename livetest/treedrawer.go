package main

import (
	"fmt"
	"strings"
)

// get_roots node tree and do post-processing
func get_roots(nw *NodeWriter) {
	for qstr, q := range nw.nodes {
		if _, ok := nw.parentof[qstr]; !ok {
			nw.roots = append(nw.roots, q)
		}
	}
}

// create HTML from object
func (nw *NodeWriter) String() string {
	get_roots(nw)
	var buf strings.Builder
	for _, q := range nw.roots {
		pipetr := make([]bool, 25) // pipetr[i] == 0 if no backgroup pipe at offset 0
		treestr(q, &buf, 0, pipetr)
		fmt.Print(basicprint(q))
	}

	return buf.String()
}

func basicprint(q *Node) string {
	s := q.val + "\nchildren:[\n"
	for _, ch := range q.chl {
		s += basicprint(ch)
	}
	s += "]\n"
	return s
}

const vpipe = "│"
const hpipe = "─"
const downbranch = "┬"
const upbranch = "┴"
const rightbranch = "├"
const lshape = "└"

func treestr(q *Node, w *strings.Builder, off int, pipetr []bool) {
	backpipes := func() string {
		s := ""
		for _, show := range pipetr[:off] {
			if show {
				s += vpipe
			}
			s += "\t"
		}
		return s
	}
	w.WriteString(q.val)
	off++
	switch len(q.chl) {
	case 0:
		//do nothing
	case 1:
		w.WriteString("\n" + backpipes() + lshape + hpipe + hpipe + hpipe + hpipe)
		treestr(q.chl[0], w, off, pipetr)
	default:
		l := len(q.chl)
		for _, cnode := range q.chl[:l-1] {
			w.WriteString("\n" + backpipes() + rightbranch + hpipe + hpipe + hpipe + hpipe)
			pipetr[off] = true
			treestr(cnode, w, off, pipetr)
		}
		pipetr[off] = false
		w.WriteString("\n" + backpipes() + lshape + hpipe + hpipe + hpipe + hpipe)
		treestr(q.chl[l-1], w, off, pipetr)
	}
}

type pnode struct {
	Dep int
	Val string
}

func getpnodes(q *Node, dep int, pns *[]pnode) {
	*pns = append(*pns, pnode{Dep: dep, Val: q.val})
	dep++
	for _, chld := range q.chl {
		getpnodes(chld, dep, pns)
	}
}

// ├── analyze_test.go
// ├──
