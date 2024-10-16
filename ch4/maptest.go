package main

import (
	"fmt"
	"sort"
)

func main() {
	// m := make(map[string]int, 0)
	var m = map[string]int{
		"max":     28,
		"rebecca": 28,
	}
	// delete(m, "rebecca")
	m["max"] = 29
	names := []string{}
	for n, a := range m {
		names = append(names, n)
		fmt.Printf("name=%s\tage=%d\n", n, a)
	}
	sort.Strings(names)
	fmt.Println(names)
	if _, ok := m["sam"]; !ok {
		fmt.Println("bad lookup")
	}
}
