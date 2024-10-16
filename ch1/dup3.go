package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	counts := make(map[string]int)
	file_origins := make(map[string]string)
	for _, fname := range os.Args[1:] {
		bytes, err := ioutil.ReadFile(fname)
		if err != nil {
			fmt.Fprintf(os.Stderr, "err:%v", err)
			continue //skip to next file
		}
		lines := strings.Split(string(bytes), "\n")
		str_present := make(map[string]bool)
		for _, line := range lines {
			counts[line]++
			str_present[line] = true
		}
		var sep string
		for line, _ := range counts {
			file_origins[line] += sep + fname
			sep = " "
		}
	}

	for line, n := range counts {
		if n > 0 {
			if line == "" {
				line = "%BLANK_LINE%"
			}
			fmt.Printf("%d\t%s\t\t%s\n", n, file_origins[line], line)
		}
	}
}
