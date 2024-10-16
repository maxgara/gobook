package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	counts := make(map[string]int)
	files := os.Args[1:]
	if len(files) == 0 {
		countLines(os.Stdin, counts)
	} else {
		for _, arg := range files {
			fp, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "err:%v", err)
				continue //skip to next file
			}
			countLines(fp, counts)
		}
	}

	for line, n := range counts {
		if n > 1 {
			if line == "" {
				line = "%BLANK_LINE%"
			}
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}

func countLines(fp *os.File, counts map[string]int) {
	input := bufio.NewScanner(fp)
	for input.Scan() {
		counts[input.Text()]++
	}
}
