// Server for Tree, Graph, etc.
package main

import (
	"log"
	"net/http"
)

func SimpleServer(s string) {
	simpleserve := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(s))
	}
	http.HandleFunc("/", simpleserve)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
func CanvasServer(c *Canvas) {
	rootHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(c.String()))
	}
	applyHandler := func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		id := r.Header.Get("HX-Trigger")
		arg := r.Form.Get("input")
		freturn := c.Apply(id, arg)
		w.Write([]byte(freturn))
	}
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/apply/", applyHandler)
	http.ListenAndServe("localhost:8000", nil)
	// http.HandleFunc(	)
}
