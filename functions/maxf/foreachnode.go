package maxf

import (
	"fmt"

	"golang.org/x/net/html"
)

func ForEachNode(n *html.Node, pre, post func(n *html.Node) error) error {
	var err error
	//apply pre func
	if pre != nil {
		err = pre(n)
	}
	if err != nil {
		return fmt.Errorf("forEachNode: pre: %v(%v):%v", pre, n, err)
	}
	ch, sb := n.FirstChild, n.NextSibling
	//recursive call
	for _, nextn := range []*html.Node{ch, sb} {
		if nextn == nil {
			continue
		}
		err = ForEachNode(nextn, pre, post)
		if err != nil {
			return err
		}
	}
	//apply post func
	if post != nil {
		err = post(n)
	}
	if err != nil {
		return fmt.Errorf("forEachNode: post: %v(%v):%v", post, n, err)
	}
	return nil
}
