// temporary file for testing while in development
package main

import (
	"fmt"
	"math/rand"
)

type exampleitem struct {
	val string
}

func (ex *exampleitem) String() string {
	return ex.val
}
func main() {
	// var strs string
	// // for range 1 {
	// oldroot := treegen2()
	// nodes := MakeTree(*oldroot)
	// nw := NewNodeWriter()
	// nw.WriteAll([]*Node(nodes))
	var ex = &exampleitem{"this\nis\na\ntest."}
	c := NewCanvas("my_canvas")
	c.NewItem("my_test", ex)
	c.NewInputTextArea("testinput", func(s string) string { ex.val += "XXX"; return "ok function called." }, "my_test")
	fmt.Printf("%#v", c)
	// strs += nw.HTMLString()
	// fmt.Println(s)
	// }
	// fmt.Println(nw)
	// fmt.Println(c.String())
	CanvasServer(c)
}

type Vertex struct {
	name             string    `NodeVal:"-"`
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
