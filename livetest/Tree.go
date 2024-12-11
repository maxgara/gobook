// This file contains the methods to build a node tree from a root node, based on struct tags.
// The function which will be exposed in the API will be MakeTree.
package main

import (
	"fmt"
	"reflect"
)

func ExampleTranslation() {
	//example of an arbitrary linked list node type
	type fakenode struct {
		value    string      `NodeVal:"name"`
		other    string      `NodeVal:"info"`
		children []*fakenode `ChildNodes:"kids"`
	}

	child := fakenode{value: "child", other: "I do not yet exist"}
	node := fakenode{value: "max gara", other: "I wrote this program", children: []*fakenode{&child}}
	parent := fakenode{value: "alan gara", other: "I would have written this program in C", children: []*fakenode{&node}}
	conf := *config(parent)
	fmt.Println(conf)
	nodes := &([]*Node{})
	realparent := TranslateNodeR(parent, conf, nodes)
	fmt.Println(realparent)
	fmt.Println(realparent.chl[0])
	fmt.Println(realparent.chl[0].chl[0])
	// Output:
	// {id=0; name:alan gara; info:I would have written this program in C;  children=1}
	// {id=1; name:max gara; info:I wrote this program;  children=1}
	// {id=2; name:child; info:I do not yet exist;  children=0}

}
func (node *Node) String() string {
	id := fmt.Sprint(node.id)
	val := node.val
	chlct := fmt.Sprint(len(node.chl))
	return fmt.Sprintf("{id=%v; %v children=%v}", id, val, chlct)
}

const (
	ChildSliceScheme = iota
	FirstChildNextSiblingScheme
	LinkedListScheme
)

// supported tags are:
// NodeVal:"<NodeValLabel>" ; if NodeValLbel=="-", then the label part will be omitted during rendering
// ChildNodes:"<ChildNodesLabel>" ; Label ignored
// FirstChildNode: "<FCLabel>" ; Label ignored
// NextSiblingNode: "<NSLabel>" ; Label ignored
// *********Important: a Label cannot be "", or the tag is ignored.***********
type graphInputConfig struct {
	scheme             int      //values defined above. determines how graph is walked and what fields in this struct are used
	NodeValIdxs        []int    //index of val in struct fields. More than one field can have this tag. field can also have other tags
	NodeValLabels      []string //labels provided by NodeVal tags
	ChildNodesIdx      int      //index of slice or pointer to child node(s) within node struct
	FirstChildNodeIdx  int      //if FCNS scheme used, we need this field
	NextSiblingNodeIdx int      //if FCNS scheme used, we need this field
}

var nullconf = graphInputConfig{scheme: -1, ChildNodesIdx: -1, FirstChildNodeIdx: -1, NextSiblingNodeIdx: -1}

// see if node is a struct with tags allowing implementation of a supported graph traversal scheme.
func config(node any) *graphInputConfig {
	conf := nullconf
	nodetype := reflect.TypeOf(node)
	kind := nodetype.Kind()
	if kind != reflect.Struct {
		panic("not a struct")
	}
	n := nodetype.NumField()
	for i := 0; i < n; i++ {
		field := nodetype.Field(i)
		fieldConf(&conf, field, i)
	}

	if conf.ChildNodesIdx != -1 {
		conf.scheme = ChildSliceScheme
	} else if conf.FirstChildNodeIdx != -1 && conf.NextSiblingNodeIdx != -1 {
		conf.scheme = FirstChildNextSiblingScheme
	} else {
		fmt.Println(conf)
		panic("can't find matching scheme for config")
	}
	return &conf
}

// helper for config
func fieldConf(conf *graphInputConfig, field reflect.StructField, idx int) {
	k := field.Type.Kind()
	tag := field.Tag
	if _, ok := tag.Lookup("ChildNodes"); ok {
		if k == reflect.Slice {
			conf.ChildNodesIdx = idx
		} else {
			panic("ChildNodes Tag on wrong type.kind: must be slice")
		}
	}
	if _, ok := tag.Lookup("NodeVal"); ok {
		label := tag.Get("NodeVal")
		conf.NodeValIdxs = append(conf.NodeValIdxs, idx)
		conf.NodeValLabels = append(conf.NodeValLabels, label)
	}
	if _, ok := tag.Lookup("NextSiblingNode"); ok {
		conf.NextSiblingNodeIdx = idx
	}
}

var translated []string //track node values which have already been translated, this allows accounting for cycles

// TranslateNodeR translates any node implementation with supported scheme to a Node tree, using configuration conf.
// translation is recursive and walks the tree/graph, translating each node it hits and placing new Nodes in list.
// if there are any cycles in the graph, they will be broken.
// In this case, all nodes linked below node will still be translated and listed.
// return value is the translation of the original node, as is list[0].
func TranslateNodeR(node any, conf graphInputConfig, list *[]*Node) *Node {
	if conf.scheme != ChildSliceScheme {
		panic("currently unsupported") //todo: support more schemes
	}
	if conf.ChildNodesIdx == -1 {
		panic("bad configuration")
	}
	id := fmt.Sprint(len(*list))
	q := Node{id: id}
	var nodeReflectVal reflect.Value
	// when this function is called recursively, after the first call the node is a reflect.Value.
	// we have to account for the change in type here:
	//todo: improve this step
	if v, ok := node.(reflect.Value); ok {
		nodeReflectVal = v
	} else {
		nodeReflectVal = reflect.ValueOf(node)
	}
	//create a value string based on the specified config. maybe this should be less confusing.
	var val string //all values (NodeVal tagged fields) formatted and combined, with labels.
	sep := ";"
	l := len(conf.NodeValIdxs)
	for lidx, vidx := range conf.NodeValIdxs {
		if lidx == l-1 {
			sep = ""
		}
		vl := conf.NodeValLabels[lidx]   //value label
		vv := nodeReflectVal.Field(vidx) //value value
		var vs string                    //value string
		if vl == "-" {
			vs = fmt.Sprintf("%v%v", vv, sep)
		} else {
			vs = fmt.Sprintf("%v:%v%v", vl, vv, sep)
		}
		val += vs
	}

	q.val = val
	translated = append(translated, val) //don't translate the same node again.
	*list = append(*list, &q)            //allow the next index to be set correctly.
	childslc := nodeReflectVal.Field(conf.ChildNodesIdx)
	l = childslc.Len()
	for i := 0; i < l; i++ {
		chelem := childslc.Index(i).Elem()
		newchild := TranslateNodeR(chelem, conf, list)
		//link nodes
		q.chl = append(q.chl, newchild)
	}
	return &q
}

func MakeTree(node any) []*Node {
	conf := *config(node)
	nodes := &([]*Node{})
	TranslateNodeR(node, conf, nodes)
	return *nodes
}
