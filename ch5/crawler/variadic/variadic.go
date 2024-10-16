package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"

	"golang.org/x/net/html"
)

func main() {
	str, err := Fetch("google.com")
	if err != nil {
		log.Fatal("bad fetch")
	}
	reader := NewReader(str)
	doc, _ := html.Parse(reader)
	els := getElementsByTagName(doc, "a")
	for _, n := range els {
		nstr := SprintNode(n)
		fmt.Printf("%s,\n", nstr)
	}
	fmt.Println(els)
}
func SprintNode(n *html.Node) string {
	return fmt.Sprintf("[Data:%s]\n", n.Data)
}

// get all nodes with element type matching one of the strings in name....
func getElementsByTagName(doc *html.Node, name ...string) []*html.Node {
	if doc == nil || name == nil {
		return nil
	}
	var out []*html.Node
	for _, n := range name {
		if doc.Data == n {
			out = append(out, doc)
		}
		c := doc.FirstChild
		out = append(out, getElementsByTagName(c, n)...)
		s := doc.NextSibling
		out = append(out, getElementsByTagName(s, n)...)
	}
	return out
}

// Fetch retrieves (HTTP GET) the resource at url and returns the
// resulting HTTP response body as a string
func Fetch(url string) (string, error) {
	//add https:// if missing
	prefix, _ := regexp.Match("^https?://", []byte(url))
	if !prefix {
		url = "https://" + url
	}
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("fetch: Get %s:%s\n", url, err)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		resp.Body.Close()
		return "", fmt.Errorf("fetch: Read body %s:%s\n", url, err)
	}
	resp.Body.Close()
	return string(b), nil
}

type ezreader struct {
	s *string
}

func NewReader(s string) *ezreader {
	var r = ezreader{&s}
	return &r
}
func (r ezreader) Read(p []byte) (int, error) {
	n := copy(p, *r.s)
	if n < len(*r.s) {
		*r.s = (*r.s)[n:]
		return n, nil
	}
	return n, io.EOF
}
