package live

// walk node tree and assign depths to nodes
func walk(q *Node, dep int) {
	q.dep = dep
	dep++
	for _, chld := range q.chl {
		walk(chld, dep)
	}
}

// actually create HTML from object
func (nw *NodeWriter) String() {

}
func (dr *drawer) close() {
	dr.String(-1, "")
}
