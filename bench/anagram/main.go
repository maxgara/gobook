package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"strings"
	"time"
)

func main() {
	//cpu profiling
	memProfile, _ := os.Create("memprofile")
	cpuProfile, _ := os.Create("cpuprofile")

	pprof.StartCPUProfile(cpuProfile)
	defer pprof.WriteHeapProfile(memProfile)
	defer pprof.StopCPUProfile()

	f, _ := os.ReadFile("./../british-english-insane.txt")
	s := string(f)
	lines := strings.Fields(s)
	fmt.Printf("%v lines\n", len(lines))
	chars := make(map[byte]int) //char c seen on line chars[c]
	for _, l := range lines {
		for i := range l {
			c := l[i]
			chars[c] = i
		}
	}
	for k, v := range chars {
		fmt.Printf("%c [%d]: %v\n", k, k, v)
	}
	fmt.Printf("%v distinct chars\n\n\n", len(chars))
	// testRun(lines)
	go run(lines)
	<-time.After(5 * time.Second)
}

func comp(s1, s2 string) bool {
	if len(s1) != len(s2) {
		return false
	}
	cslc := [255]byte{}
	for i := range s1 {
		c1, c2 := s1[i], s2[i]
		cslc[c1]++
		cslc[c2]--
	}
	return cslc == [255]byte{}
}

func run(lines []string) {
	dl := len(lines)
	//send progress updates
	var due bool
	go func() {
		for {
			<-time.After(time.Second)
			due = true
		}

	}()
	var groups [][]string
	var t int //seconds passed
	for i, s1 := range lines[:dl-1] {
		var group []string
		for _, s2 := range lines[i+1:] {
			if due {
				t++
				p := float32(i) * 100 / float32(dl) //percent done
				fmt.Printf("\r%d. i=%v; %.3f%%; [avg %.3f%%/s]\t\n", t, i, p, p/float32(t))
				due = false
			}
			if comp(s1, s2) {
				group = append(group, s2)
			}
		}
		if len(group) != 0 {
			group = append(group, s1)
			groups = append(groups, group)
		}
	}
	fmt.Println(groups)
}
