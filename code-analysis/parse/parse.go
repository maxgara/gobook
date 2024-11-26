// this package allows easier parsing of complex strings using regular expressions. NewParseNd creates a new root node in a parse tree,
// which is then expanded using calls to Parse(regex). Important: the regular expression must contain 1 named capture group!
// calls to Parse can be chained to quickly build up a tree of property nodes. Parse can be called on a single node or a group (ParseG) of nodes,
// allowing calls to Parse to be chained in order to quickly build a parse tree.
// Extracting information from the tree structure can be done using the .Val property of individual nodes, calling Walk to walk the whole tree,
// or using a pReader to get direct property values for the current node. Temp and Refresh methods support more complex parsing
// by facilitating easy creation of temporary parse nodes. In the future there may be support for extraction of multiple properties
// using a single regex string.
// Future: There is also potential for a live extraction viewer application to see extraction results in real time,
// as well as a generating command to create optimized static structures and extractions from the dynamic parse tree.
package parse

import (
	"fmt"
	"regexp"
	"strings"
)

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
	Idx  int    //relative to base
	fqn  string //fully qualified name
}

func NewParseNd(val string) *ParseNd {
	base := []byte(val)
	return &ParseNd{Val: base, p: make(map[string]ParseG), Idx: 0, Base: &base}
}
func (q *ParseNd) String() string {
	return q.rString(0, "NODE")
}
func (q *ParseNd) rString(ind int, name string) string {
	var s string
	if name == "" {
		name = "Node"
	}
	s += q.stringnd(ind, name)
	ind++
	for pname, g := range q.p {
		for idx, nd := range g {
			name := fmt.Sprintf(".%v[%v]", pname, idx)
			s += nd.rString(ind, name)
		}
	}
	return s
}

// simple node string func without props, using specified indent and name
func (q *ParseNd) stringnd(ind int, name string) string {
	valstr := "\"" + string(q.Val) + "\""
	arr := strings.Split(valstr, "\n")
	inds := strings.Repeat("\t", ind)
	out := inds + name + ":"
	//one line value: keep on one line
	if len(arr) == 1 {
		out += valstr + "\n"
		return out
	}
	//multiline value: place val and node name on different lines
	for _, l := range arr {
		out += "\n" + inds + l
	}
	out += "\n"
	return out
}
func (q ParseG) String() string {
	var out string
	for _, v := range q {
		out += v.String()
	}
	return out
}

func (q *ParseNd) Walk(f func(q *ParseNd) (cont bool)) {
	if !f(q) {
		return
	}
	for _, g := range q.p {
		for _, next := range g {
			next.Walk(f)
		}
	}
}

// property reader. Does not read sub-properties.
type preader struct {
	idx int
	Slc []pinfo
	pinfo
}
type pinfo struct {
	Nd    *ParseNd
	Pname string // current property name
	Gidx  int    //index of node in group Nd.p[Pname]
}

// create a new property reader for node q.
func (q *ParseNd) newPreader() *preader {
	var slc []pinfo
	for name, g := range q.p {
		for gidx, nd := range g {
			info := pinfo{Nd: nd, Pname: name, Gidx: gidx}
			slc = append(slc, info)
		}
	}
	return &preader{Slc: slc}
}
func (pr *preader) Read() *ParseNd {
	if pr.idx >= len(pr.Slc) {
		return nil
	}
	inf := pr.Slc[pr.idx]
	pr.pinfo = inf
	pr.idx++
	return inf.Nd
}
func (pr *preader) Text() string {
	return pr.Nd.String()
}
func (pr *preader) String() string {
	return fmt.Sprintf("%v[%v]", pr.Pname, pr.Gidx)
}

func (q *ParseNd) Refresh() {
	// clear all temp properties
	q.Walk(func(q *ParseNd) bool {
		if !q.temp {
			return true
		}
		parent := q.anc
		for k, g := range parent.p {
			for _, v := range g {
				if v == q {
					delete(parent.p, k)
				}
			}
		}
		return true
	})
	//add fqns
	q.Walk(func(q *ParseNd) bool {
		if q.anc == nil {
			return true
		}
		q.fqn = q.anc.fqn + "/" + q.Name()
		return true
	})
}

// retrieve node's "name" from ancestor in format: propertyname[groupindex] = q.Name(); anc.p[propertyname][groupindex] = q
func (q *ParseNd) Name() string {
	r := q.anc.newPreader()
	for r.Read() != nil {
		if r.Nd == q {
			return r.String()
		}
	}
	return "" //root node or misconfiguration
}

// parse named subgroups as properties of q. if q is temp, then add properties to first non-temp ancestor.
func (q *ParseNd) Parse(pattern string) ParseG {
	p := regexp.MustCompile(pattern)
	arr := p.FindAllSubmatchIndex(q.Val, MAXMATCH) // only first submatch of each match is used
	name := pname(pattern)
	g := slcgrp(arr, q, name)
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
	nstr := np.FindString(p)            // name
	if len(nstr) == 0 {
		return "" // no name found
	}
	nstr = nstr[2 : len(nstr)-1] //drop <? and >
	return nstr
}

// convert arr to ParseG and add to q as prop name
func slcgrp(arr [][]int, q *ParseNd, name string) ParseG {
	var g ParseG
	target := q
	//find first non-temp ancestor as target for new property
	for target.temp {
		target = target.anc
	}
	for _, match := range arr {
		if len(match) <= 2 {
			return nil //no submatches in pattern, group undefined
		}
		//match 2, 3 are the start, end idxs of first submatch
		start := match[2]
		stop := match[3]
		val := q.Val[start:stop]
		newnode := ParseNd{Val: val, p: map[string]ParseG{}, anc: target, Base: q.Base, Idx: start + q.Idx}
		g = append(g, &newnode)
	}

	target.p[name] = g
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
