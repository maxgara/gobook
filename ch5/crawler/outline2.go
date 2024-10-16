package main

// outline2 iterates thru an html node tree similarly to outline, but allows the caller to supply a pair
// of functions f1, f2 (*html.Node) to call before and after processsing child nodes.
import (
	"fmt"
	"log"
	"os"

	"golang.org/x/net/html"
)

var depth int
var nDepth int

func test() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	f1 := func(n *html.Node) error {
		if n.Type == html.ElementNode {
			fmt.Printf("%*s<%v>\n", nDepth, " ", n.Data)
			nDepth++
		}
		return nil
	}
	f2 := func(n *html.Node) error {
		if n.Type == html.ElementNode {
			fmt.Printf("%*s</%v>\n", nDepth, " ", n.Data)
			nDepth--
		}
		return nil
	}
	ForEachNode(doc, f1, f2)
}

func ForEachNode(n *html.Node, pre, post func(n *html.Node) error) error {
	var err error
	if pre != nil {
		err = pre(n)
	}
	if err != nil {
		return fmt.Errorf("ForEachNode: pre: %v(%v):%v", pre, n, err)
	}
	ch, sb := n.FirstChild, n.NextSibling
	for _, nextn := range []*html.Node{ch, sb} {
		if nextn == nil {
			continue
		}
		err = ForEachNode(nextn, pre, post)
		if err != nil {
			return err
		}
	}
	if post != nil {
		err = post(n)
	}
	if err != nil {
		return fmt.Errorf("ForEachNode: post: %v(%v):%v", post, n, err)
	}
	return nil
}
