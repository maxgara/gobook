package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"time"
)

func main() {
	start := time.Now()
	m := make(map[rune]bool)
	f, err := os.Open("The Go Programming Language.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	buf := bufio.NewReader(f)
	totalRead := 0
	for {
		r, _, err := buf.ReadRune()
		if err == io.EOF {
			fmt.Println("io.EOF hit")
			break
		} else if err != nil {
			fmt.Println("unexpected err hit: " + err.Error())
			os.Exit(1)
		}
		totalRead++
		// fmt.Printf("%c", r)
		m[r] = true
	}
	var chars = []string{}
	for c, _ := range m {
		if c != '\uFFFD' {
			chars = append(chars, string(c))
		}
	}
	sort.Strings(chars)
	secs := time.Since(start).Seconds()
	fmt.Println(chars)
	fmt.Printf("total characters read:%d\n", totalRead)
	fmt.Printf("Execution time:%fs\n", secs)

}
