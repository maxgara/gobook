package live

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
)

var people = [...]string{"arthur", "blake", "charles", "drew", "edith", "francisco", "gwen"}

func TestDrawer(t *testing.T) {
	for range 5 {
		root := treegen()
		fmt.Println(root)
		dr := drawer{w: os.Stdout}
		walk(root, 0, &dr)
	}
}
func treegen() *Node {
	var tree []*Node
	tree = append(tree, &Node{id: people[0]})
	for _, name := range people[1:] {
		nn := &Node{id: name}
		r := rand.Intn(len(tree))
		tree[r].chl = append(tree[r].chl, nn)
		tree = append(tree, nn)
	}
	return tree[0]
}
