package main

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
		err := exec.Command(execfp).Run()
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
