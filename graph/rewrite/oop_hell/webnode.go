package main

import (
	"fmt"
	"hash"
	"io"
	"os"
	"strings"
)

// docbuilder builds a document, navigating through elements as it works.
type docBuilder struct {
	w     io.Writer
	idx   int //for new nodes
	loc   *Node
	atmap map[struct{string, int}]string
	root  *Node
}
type attribute struct {
	id  int //node idx
	key string
	h {int,}
}

// basic element in doc, can be wrapped inside a Content element.
type Node struct {
	name  string
	id    int //idx of node
	class string
	attrs *map[string]string //kv pairs for html opening tag
	chl   []*Node            //children
	par   *Node              //parent
	c     Content
	cstr  string //content string, for elements writing their own strings
	ok    bool   //ok after initialization and before any copying process
}

// interface used for rendering a Node.
// implementation will typically contain the corresponding Node, which is fixed size, and may make use of printNode in its String function.
// Node access functions should be used carefully.
type Content interface {
	String() string
	Node() *Node //retrieve corresponding Node for Content (inverse of Node.c)
}

// start a document
func newDocBuilder(w io.Writer) *docBuilder {
	b, err := os.ReadFile(style_file)
	if err != nil {
		panic(err)
	}
	style = string(b)
	r := &root{n: Node{}}
	r.n.c = r //god I hate this... why did I make it so complicated?
	return &docBuilder{w: w, loc: &r.n, root: &r.n}
}

// read completed buffer string
func (b *buffer) String() string {
	s := (*strings.Builder)(b).String()
	*b = buffer{} //reset buffer before returning
	return s
}

// Node String() -> Content String()
func (n Node) String() string {
	return (n.c).String()
}

// add Node n as a child of current Node, then navigate b to n. Do not use without initializing n.c
// if n.c is not initialized, program will panic on a later call to n.String()
func (b *docBuilder) StartNode(n *Node) {
	l := b.loc
	l.chl = append(l.chl, n) //add child to parent
	n.par = l                //add parent to child
	b.loc = n                //navigate to child
}

// link n and c, make n a child of current node at builder b, then navigate b to n
func (b *docBuilder) StartContentNode(c Content, n Node) {
	//copy n into c
	np := c.Node()
	*np = Node{name: n.name, class: n.class, c: c} //add non-tree vals
	copy(np.attrs, n.attrs)                        //deep copy attributes
	b.StartNode(np)                                // link builder loc <-> n , nav loc => n
}

// navigate b up one element in document tree
func (b *docBuilder) EndNode() {
	b.loc = b.loc.par
}

func (b *docBuilder) SendDoc() {
	const fstr = `<html><head><style>%v</style></head><body>%v</body></html>`
	cstr := b.root.String() //content string
	b.w.Write([]byte(fmt.Sprintf(fstr, style, cstr)))
}

// raw string content without html tags
type raw struct {
	n Node
	s string
}

func (r *raw) String() string {
	return r.s
}

type root struct {
	n Node
}

func (r *root) String() string {
	var s string
	for _, v := range r.n.chl {
		s += v.String()
	}
	return s
}
func (r *root) Node() *Node {
	return &r.n
}
