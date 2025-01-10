package main

//run a command then serve a file whenever an HTTP GET is received from localhost:8001

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	LOCAL_FILE = iota
	LOCALHOST_NETCONN
)

func main() {
	execfp := os.Args[1]
	readfp := os.Args[2]
	var mode = LOCAL_FILE //default
	if strings.HasPrefix(readfp, "localhost") {
		mode = LOCALHOST_NETCONN
	}
	fmt.Printf("mode %v\n", mode)
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("exec command [%v]\n", execfp)
		cmd := exec.Command(execfp)
		err := cmd.Start()
		if err != nil {
			fmt.Printf("exec err:%v\n", err)
			panic("")
		}
		b := read(mode, readfp)
		cmd.Process.Kill()
		w.Write(b)
	}
	http.HandleFunc("/", handler)
	http.ListenAndServe("localhost:8001", nil)
}

// read all bytes from fp
func read(mode int, fp string) []byte {
	if mode == LOCAL_FILE {
		b, err := os.ReadFile(fp)
		if err != nil {
			fmt.Printf("readfile err:%v\n", err)
			panic("")
		}
		return b
	}
	// LOCALHOST_NETCONN case
	// allow execfp proc time to open network interface, otherwise the first page load will get TCP_CONN_REFUSED
	// + the associated quickserve goroutine will crash and prevent the execfp function from being closed.
	// further page loads will appear to work normally but will be communicating with the original execfp proc, not
	// updated versions of the same. (confusing error)
	time.Sleep(time.Millisecond * 100)
	r, err := http.Get("http://" + fp + "/")
	if err != nil {
		fmt.Printf("http get err:%v\n", err)
		panic("")
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("resp read: %v\n", err)
		panic("")
	}
	return b
}
