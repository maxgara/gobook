package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

const RECURSION_LIMIT = 2

type Page struct {
	Url      string
	Children []*Page
	err      error
}

func main() {
	url := os.Args[1]
	tree := makeTree(url)
	// if nDepth > RECURSION_LIMIT {
	// 	fmt.Printf("Hit Recursion limit %d\n", RECURSION_LIMIT)
	// }
	printTree(tree)
}

// accept url and recursively visit pages. Return Page-tree.
// if url cannot be retrieved or if link scraping fails,
// returned page will have non-nil err val and nil Children val.
func makeTree(url string) *Page {
	return makeBranch(url, 0)
}
func makeBranch(url string, dep int) *Page {
	dep++
	p := Page{url, nil, nil}
	if dep > RECURSION_LIMIT {
		p.err = fmt.Errorf("Past recursion limit")
		return &p
	}
	links, err := extractLinks(url)
	if err != nil {
		p.err = err //if page cannot be retrieved or link extraction otherwise fails, set err
		return &p
	}
	for _, l := range links {
		ch := makeBranch(l, dep)
		p.Children = append(p.Children, ch)
	}
	return &p
}

// recursively extracts all links (<a href=...> elements) from HTML Doc into slice
func extractLinks(url string) ([]string, error) {
	var links = []string{}
	resp, err := Fetch(url)
	if err != nil {
		// fmt.Fprintf(os.Stderr, "Fetch %s:%s", url, err)
		return nil, fmt.Errorf("ExtractLinks %s Fetch:%s", url, err)
	}
	reader := NewReader(resp)      // make io.Reader from string for Parse function call
	doc, err := html.Parse(reader) // skip page on parsing error to avoid runtime error
	if err != nil {
		return nil, fmt.Errorf("ExtractLinks %s HTMl.Parse:%s", url, err)
	}
	getlink := func(n *html.Node) error {
		if n.Type != html.ElementNode || n.Data != "a" {
			return nil
		}
		for _, prop := range n.Attr {
			if prop.Key != "href" {
				continue
			}
			links = append(links, prop.Val)
		}
		return nil
	}
	ForEachNode(doc, getlink, nil)
	return links, nil
}

func printTree(t *Page) {
	dep := 0
	printTreeD(t, dep)
}

// print tree with depth arg
func printTreeD(t *Page, d int) {
	for i := 0; i < d; i++ {
		fmt.Printf("-")
	}
	if t.err != nil {
		fmt.Printf("ERROR:%s %s\n", t.err, t.Url)
		return
	}
	fmt.Printf("%s\n", t.Url)
	for _, c := range t.Children {
		printTreeD(c, d+1)
	}
}
