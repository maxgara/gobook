package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"strings"
	"time"
)

func getAnagrams(s string) {
	lines := strings.Fields(s)
	fmt.Printf("%v lines\n", len(lines))
	charMap := make(map[byte]int) //char c seen on line chars[c]
	for _, l := range lines {
		for i := range l {
			c := l[i]
			charMap[c] = i
		}
	}
	for k, v := range charMap {
		fmt.Printf("%c [%d] %v\n", k, k, v)
	}
	fmt.Printf("%v distinct chars\n\n\n", len(charMap))
	// testRun(lines)

	m := dictToMap(lines)
	var i int
	for k, v := range m {
		if i > 100 {
			break
		}
		fmt.Printf("%v => %v\n", k, v)
		i++
	}
}
func main() {
	//cpu/memory profiling
	var pprofCpu = flag.Bool("profilecpu", false, "write Cpu profile")
	var pprofMem = flag.Bool("profilememory", false, "write memory profile")
	flag.Parse()
	if *pprofCpu {
		cpuProfile, _ := os.Create("cpuprofile")
		pprof.StartCPUProfile(cpuProfile)
		defer pprof.StopCPUProfile()
	}
	if *pprofMem {
		memProfile, _ := os.Create("memprofile")
		sig := make(chan bool)
		go func(c chan bool) {
			defer pprof.WriteHeapProfile(memProfile)
			<-c
		}(sig)
	}
	data, _ := os.ReadFile("./../british-english-insane.txt")
	dataPrep(data)
	// <-time.After(5 * time.Second)

}

func dataPrep(data []byte) {
	lines := strings.Fields(string(data))

}

func dictToMap(lines []string) map[string][]string {
	m := make(map[string][]string)
	for _, l := range lines {
		ccs := make([]byte, 255) //char counts
		for _, c := range l {
			ccs[c]++
		}
		m[string(ccs)] = append(m[string(ccs)], l)
	}
	return m
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
