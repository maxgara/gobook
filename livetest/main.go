package main

import (
	"fmt"
	"math/rand"
)

func main() {
	oldroot := treegen2()
	// fmt.Println(oldroot)
	nodes := MakeTree(*oldroot)
	// fmt.Println()
	nw := NewNodeWriter()
	nw.WriteAll([]*Node(nodes))
	// for _, n := range nodes {
	// 	var chstrs []string
	// 	for _, ch := range n.chl {
	// 		chstrs = append(chstrs, ch.id)
	// 	}
	// 	nw.Write(n.id, n.val, chstrs)
	// }
	fmt.Println(nw)
}

type Vertex struct {
	name             string    `NodeVal:"first_name"`
	descendent_names []*Vertex `ChildNodes:"kids"`
}

func treegen2() *Vertex {
	var people2 = [...]string{"arthur", "blake", "charles", "drew", "edith", "francisco", "gwen", "xarthur", "xblake", "xcharles", "xdrew", "xedith", "xfrancisco", "xgwen", "zarthur", "zblake", "zcharles", "zdrew", "zedith", "zfrancisco", "zgwen"}
	var tree []*Vertex
	tree = append(tree, &Vertex{name: people2[0]})
	for _, name := range people2[1:] {
		nn := &Vertex{name: name}
		r := rand.Intn(len(tree))
		tree[r].descendent_names = append(tree[r].descendent_names, nn)
		tree = append(tree, nn)
	}
	return tree[0]
}
