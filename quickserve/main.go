package main

//run a command then serve a file whenever an HTTP GET is received from localhost:8001
//usage: quickserve <executable_path> <output_file_or_interface>
//ex. quickserve ./myproc localhost:8000

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
	procName := os.Args[1]
	procRes := os.Args[2]
	var mode = LOCAL_FILE //default
	if strings.HasPrefix(procRes, "localhost") ||
		strings.HasPrefix(procRes, "http://localhost") {
		mode = LOCALHOST_NETCONN
	}
	done := make(chan *exec.Cmd)
	killLastProc := func() {

		var c *exec.Cmd
		c = <-done
		err := c.Process.Kill()
		if err != nil {
			panic("couldn't kill process.")
		}
	}
	var first = true
	handler := func(w http.ResponseWriter, r *http.Request) {
		if !first {
			killLastProc()
		}
		first = false
		cmd := exec.Command(procName)
		errp, err := cmd.StderrPipe()
		if err != nil {
			panic(fmt.Sprintf("exec stderr pipe connection:%v", err))
		}
		outp, err := cmd.StdoutPipe()
		if err != nil {
			panic(fmt.Sprintf("exec stdout pipe connection:%v", err))
		}
		go io.Copy(os.Stdout, errp)
		go io.Copy(os.Stdout, outp)

		err = cmd.Start()
		if err != nil {
			panic(fmt.Sprintf("exec err (proc=[%v]):%v\n", procName, err))
		}
		b := read(mode, procRes)
		done <- cmd
		_, err = w.Write(b)
		if err != nil {
			panic("couldn't write to testserver client")
		}
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
		time.Sleep(time.Millisecond * 1000)
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
