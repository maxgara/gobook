package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"
)

func main() {
	start := time.Now()
	f, err := os.Open("The Go Programming Language.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	buf := bufio.NewReader(f)
	totalRead := 0
	for {
		_, _, err := buf.ReadRune()
		if err == io.EOF {
			break
		}
		totalRead++
	}
	secs := time.Since(start).Seconds()
	fmt.Printf("total characters read:%d\n", totalRead)
	fmt.Printf("Execution time:%fs\n", secs)

}
