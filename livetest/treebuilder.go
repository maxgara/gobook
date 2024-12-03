package live

//rename this file
import "fmt"

type Node struct {
	id  string
	chl []*Node
	dep int //depth of node in tree, root(s) have 0
}

// TODO: when more of the program is done, this will become an instance of a HTMLCanvas struct.
// Actually word canvas might be confusing because of the existing <canvas> HTML tag..
// maybe think of something else.
// node processor, reconstructs tree from nodes
type NodeWriter struct {
	nodes    map[string]*Node //map of nodes: [id_string]*node
	parentof map[string]*Node //map of nodes [child_id_string]*parent_node
}

func NewNodeWriter() *NodeWriter {
	return &NodeWriter{nodes: make(map[string]*Node), parentof: make(map[string]*Node)}
}

// should always work in one pass, doesn't get depths
func (n *NodeWriter) Write(id any, cn []any) {
	idstr := fmt.Sprint(id)
	nn := Node{id: idstr}
	//find children
	for _, id := range cn {
		chldIdStr := fmt.Sprint(id)
		if chlnode, ok := n.nodes[chldIdStr]; ok {
			nn.chl = append(nn.chl, chlnode)
		}
		n.parentof[chldIdStr] = &nn
		//find parents
		if parent, ok := n.parentof[idstr]; ok {
			parent.chl = append(parent.chl, &nn)
		}
	}
}
