package main

import (
	"fmt"
	"net/http"
	"os"
)

var t []byte

func main() {
	var err error
	t, err = os.ReadFile("/Users/maxgara/Desktop/go-code/gobook/workspace/quickserve/testsrv/text.txt")
	if err != nil {
		panic(fmt.Sprintf("text file read: %v", err))
	}
	http.HandleFunc("/", h)
	http.ListenAndServe("localhost:8000", nil)
}

func h(w http.ResponseWriter, r *http.Request) {
	s := fmt.Sprintf("read %d bytes: %v", len(t), string(t))
	w.Write([]byte(s))
}
