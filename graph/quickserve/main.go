package main

//run a command then serve a file whenever an HTTP GET is received from localhost:8001

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	execfp := os.Args[1]
	readfp := os.Args[2]
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("exec command [%v]\n", execfp)
		cmd := exec.Command("/bin/zsh", execfp)
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		b, err := os.ReadFile(readfp)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		w.Write(b)
	}
	http.HandleFunc("/", handler)
	http.ListenAndServe("localhost:8001", nil)
}
