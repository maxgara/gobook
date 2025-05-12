package main

import (
	"fmt"
	"os"
	"sort"
)

func main() {
	db := createDB()
	fmt.Println(db)
	fmt.Println("vim-go")
	//getExts()
}
func getExts() {
	exts := make(map[string]int)
	extsR("/Users/maxgara/", &exts)
	var extSlice []struct {
		ex    string
		count int
	}
	for k, v := range exts {
		extSlice = append(extSlice, struct {
			ex    string
			count int
		}{k, v})
	}
	sort.Slice(extSlice, func(i int, j int) bool {
		return extSlice[i].count > extSlice[j].count
	})
	fmt.Println(extSlice[:100])
}
func extsR(dir string, exts *map[string]int) {
	fs, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, f := range fs {
		name := f.Name()
		if f.Type().IsDir() {
			extsR(dir+name+"/", exts)
			continue
		}
		e := ext(name)
		(*exts)[e]++
	}
}

// get extension from filename
func ext(f string) string {
	for i := len(f) - 1; i >= 0; i-- {
		if f[i] == '.' {
			return f[i:]
		}
	}
	return ""
}
