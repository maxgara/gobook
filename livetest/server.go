package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"maxgara-code.com/workspace/code-analysis/parse"
)

var templ *template.Template

func main() {
	test()
}

// web string interface.
// WString should return an HTML string that can be displayed in a <div> element within a webpage.
type WStringer interface {
	WString() string
}

// serve object's WString within a basic HTML template, server address is localhost:8000.
func LivePrint(ob WStringer) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		templ.Execute(w, ob)
	}

	http.HandleFunc("/", handler)
	var err error
	templ, err = template.ParseFiles("basic.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

type tester struct {
	parse.ParseNd
}

func (t tester) WString() string {
	return t.ParseNd.String()
}

func test() {
	root := parse.NewParseNd("this is my new testing parse node string. Check it out.")
	root.Parse(`(?<longwords>\w{4,})`)
	t := tester{ParseNd: *root}
	fmt.Println(t.WString())
	LivePrint(t)
}
