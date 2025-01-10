package main

import "net/http"

func main() {
	http.HandleFunc("/", h)
	http.ListenAndServe("localhost:8000", nil)
}

func h(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("testing"))
}
