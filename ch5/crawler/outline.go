package main

import (
	"fmt"
	"golang.org/x/net/html"
	"log"
	"os"
)

func outlineTest() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	outline(nil, doc)
}

func outline(stack []string, n *html.Node) {
	if n.Type == html.ElementNode {
		stack = append(stack, n.Data)
		fmt.Println(stack)
	}
	if ch := n.FirstChild; ch != nil {
		outline(stack, ch)
	}
	if sb := n.NextSibling; sb != nil {
		outline(stack, sb)
	}

}
