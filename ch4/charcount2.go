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
	m := make(map[rune]int)
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
			//fmt.Println("io.EOF hit")
			break
		} else if err != nil {
			fmt.Println("unexpected err hit: " + err.Error())
			os.Exit(1)
		}
		totalRead++
		m[r]++
	}
	var minv = make(map[int][]rune) //m-inverse.
	for k, v := range m {
		if k == '\uFFFD' { //utf-8 replacement char
			continue
		}
		_, ok := minv[v]
		if ok {
			minv[v] = append(minv[v], k) //already a slice for this key (two chars have the same count)
		} else {
			minv[v] = []rune{k} //no slice for this key, need to make one
		}
	}
	//sort by char count for printing
	cts := make([]int, 0, len(minv))
	for k := range minv {
		cts = append(cts, k)
	}
	sort.Ints(cts)
	for i := len(cts) - 1; i >= 0; i-- {
		ct := cts[i]
		for _, chr := range minv[ct] {
			fmt.Printf("%q\t%d\n", chr, ct)
		}
	}

	secs := time.Since(start).Seconds()

	fmt.Printf("total characters read:%d\n", totalRead)
	fmt.Printf("Execution time:%fs\n", secs)

}
