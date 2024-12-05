package main

import (
	"fmt"
	"math/rand"
	"testing"
)

func main() {
	r := treegen()
	w := NewNodeWriter()
	WriteTreeToWriter(r, w)

	fmt.Printf("after treestr:\n%v\n", w.String())
}

var people = [...]string{"arthur", "blake", "charles", "drew", "edith", "francisco", "gwen", "xarthur", "xblake", "xcharles", "xdrew", "xedith", "xfrancisco", "xgwen", "zarthur", "zblake", "zcharles", "zdrew", "zedith", "zfrancisco", "zgwen"}

func treegen() *Node {
	var tree []*Node
	tree = append(tree, &Node{id: people[0], val: people[0]})
	for _, name := range people[1:] {
		nn := &Node{id: name, val: name}
		r := rand.Intn(len(tree))
		tree[r].chl = append(tree[r].chl, nn)
		tree = append(tree, nn)
	}
	return tree[0]
}
func TestTreeBuilder(t *testing.T) {
	r := treegen()
	fmt.Printf("treegen: %v\n", r)
	w := NewNodeWriter()
	WriteTreeToWriter(r, w)
	fmt.Printf("treewriter, after writing: %v\n", w)
	fmt.Println(w)
}

func TestNodeWriter_String(t *testing.T) {
	r := treegen()
	w := NewNodeWriter()
	WriteTreeToWriter(r, w)
	fmt.Print(w.String())
}
func TestSimpleServer(t *testing.T) {
	r := treegen()
	w := NewNodeWriter()
	WriteTreeToWriter(r, w)
	SimpleServer(w.String())
}

func WriteTreeToWriter(r *Node, w *NodeWriter) {
	var chlstrs []string
	for _, c := range r.chl {
		chlstrs = append(chlstrs, c.id)
	}
	w.Write(r.id, r.val, chlstrs)
	for _, chld := range r.chl {
		WriteTreeToWriter(chld, w)
	}
}
