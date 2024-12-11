package main

import (
	"fmt"
	"os"
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

// create string from object
func (nw *NodeWriter) String() string {
	get_roots(nw)
	var buf strings.Builder
	for _, q := range nw.roots {
		pipetr := make([]bool, 25) // pipetr[i] == 0 if no background pipe at offset i.
		treestr(q, &buf, 0, pipetr)
	}
	return buf.String()
}

// convenience func. -horrible formatting but good enough for debugging
func basicString(q *Node) string {
	s := q.val + "\nchildren:[\n"
	for _, ch := range q.chl {
		s += basicString(ch)
	}
	s += fmt.Sprintf("]{end %v}\n", q.val)
	return s
}

const vpipe = "│"
const hpipe = "─"
const rightbranch = "├"
const lshape = "└"

// compose a string displaying a tree from root node q
func treestr(q *Node, w *strings.Builder, off int, pipetr []bool) {
	if off == len(pipetr) {
		panic("recursion limit (default limit=25)")
	}
	//string of spaces and pipes continuing "behind" column off
	backpipes := func() string {
		s := ""
		for _, show := range pipetr[:off] {
			if show {
				s += vpipe
			} else {
				s += " " //this is to prevent extra offset from backpipes
			}
			s += "  "
		}
		return s
	}
	w.WriteString(q.val)
	var linefstr = "\n" + backpipes() + "%v" + hpipe + hpipe //general format string for child
	switch len(q.chl) {
	case 0:
		//do nothing
	case 1:
		s := fmt.Sprintf(linefstr, lshape)
		w.WriteString(s)
		treestr(q.chl[0], w, off+1, pipetr)
	default:
		l := len(q.chl)
		for _, cnode := range q.chl[:l-1] {
			pipetr[off] = true //toggle pipetr so that deeper levels have pipe at current off
			s := fmt.Sprintf(linefstr, rightbranch)
			w.WriteString(s)
			treestr(cnode, w, off+1, pipetr)
		}
		pipetr[off] = false //toggle off pipetr
		s := fmt.Sprintf(linefstr, lshape)
		w.WriteString(s)
		treestr(q.chl[l-1], w, off+1, pipetr)
	}
}
func (nw *NodeWriter) HTMLString() string {
	s := nw.String()
	arr := strings.Split(s, "\n")
	s = ""
	for _, line := range arr {
		line = strings.ReplaceAll(line, " ", "&nbsp")
		s += "<div>" + line + "</br></div>"
	}
	temp, _ := os.ReadFile("tree.tmpl")
	s = fmt.Sprintf(string(temp), s)
	return s
}
