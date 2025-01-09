package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	http.HandleFunc("localhost:8000", handler)
	execfp := os.Args[1]
	readfp := os.Args[2]
	handler := func(http.ResponseWriter, *http.Request) {
		exec.Command(execfp).Run()
		b, err := os.ReadFile(readfp)
		if err != nil {
			fmt.Println(err)
			return
		}

	}

}
