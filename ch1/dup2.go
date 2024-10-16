package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	counts := make(map[string]int)
	origins := make(map[string]string)
	files := os.Args[1:]
	if len(files) == 0 {
		countLines(os.Stdin, counts, origins)
	} else {
		for _, fname := range files {
			f, err := os.Open(fname)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error reading file %s. Error description:%s", fname, err)
				continue
			}
			countLines(f, counts, origins)
		}
	}

	for line, n := range counts {
		if n > 1 {
			// if line == "" {
			// 	line = "%BLANK_LINE%"
			// }
			fmt.Printf("%10s: %d \t(origins:%s)\n", line, n, origins[line])
		}
	}
}

// count occurances of lines in f, store counts in counts, track origins of each line in origins
func countLines(f *os.File, counts map[string]int, origins map[string]string) {
	org := f.Name()
	input := bufio.NewScanner(f)
	for input.Scan() {
		txt := input.Text()
		origins[txt] = add_origin(origins[txt], org)
		counts[txt]++
	}
}

// add file origin for line of text if not present. if present already do nothing. return new string of origins
func add_origin(s string, org string) string {
	if strings.Contains(s, org) {
		return s
	} else if len(s) == 0 {
		return org
	}
	return s + ", " + org
}
