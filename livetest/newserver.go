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
