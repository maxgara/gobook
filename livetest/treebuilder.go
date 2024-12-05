package main

//rename this file

type Node struct {
	id  string
	val string
	chl []*Node
	dep int
}

// TODO: when more of the program is done, this will become an instance of a HTMLCanvas struct.
// Actually word canvas might be confusing because of the existing <canvas> HTML tag..
// maybe think of something else.
// node processor, reconstructs tree from nodes
type NodeWriter struct {
	nodes    map[string]*Node //map of nodes: [id_string]*node
	parentof map[string]*Node //map of nodes [child_id_string]*parent_node
	roots    []*Node
}

func NewNodeWriter() *NodeWriter {
	return &NodeWriter{nodes: make(map[string]*Node), parentof: make(map[string]*Node)}
}

// todo: turn this into a post-processing func, only run when re-rendering and adding nodes.
// todo:  Write optimized version for first render
// should always construct tree in one pass, doesn't get depths or id root nodes.
// id = new node id; val = new node val; chl= IDs of child nodes
func (nw *NodeWriter) Write(id string, val string, chl []string) {
	nn := Node{id: id, val: val}
	nw.nodes[id] = &nn
	//find children
	for _, cid := range chl {
		if child, ok := nw.nodes[cid]; ok {
			nn.chl = append(nn.chl, child)
			continue
		}
		nw.parentof[cid] = &nn
	}
	//find parent
	if parent, ok := nw.parentof[id]; ok {
		parent.chl = append(parent.chl, &nn)
	}
}
