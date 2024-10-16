package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
	"maxgara-code.com/workspace/functions/maxf"
)

func main() {
	site := os.Args[1]
	err := title(site)
	if err != nil {
		fmt.Printf("err:%s\n", err)
	}
	maxf.Fetch("example.com")
}
func title(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	//check content-type is html
	ct := resp.Header.Get("Content-Type")
	if ct != "text/html" && !strings.HasPrefix(ct, "text/html;") {
		return fmt.Errorf("%s has content type %s, not text/html", url, ct)
	}
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return fmt.Errorf("parsing %s as html: %s", url, err)
	}
	visitNode := func(n *html.Node) error {
		if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
			fmt.Println(n.FirstChild.Data)
		}
		return nil
	}
	maxf.ForEachNode(doc, visitNode, nil)
	return nil
}
