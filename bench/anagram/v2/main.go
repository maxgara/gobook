package main

import (
	"flag"
	"io"
	"log"
	"os"
	"runtime/pprof"
	"strings"
)

var file = "../british-english-insane.txt"
var out io.Writer

func main() {
	//cpu/memory profiling
	var pprofCPU = flag.Bool("profilecpu", false, "write Cpu profile")
	var pprofMem = flag.Bool("profilememory", false, "write memory profile")
	flag.Parse()
	var stopCalls []func()
	if *pprofCPU {
		cpuProfile, _ := os.Create("cpuprofile") //create cpu profile log file
		pprof.StartCPUProfile(cpuProfile)
		stopCalls = append(stopCalls, func() {
			pprof.StopCPUProfile()
		})
	}
	if *pprofMem {
		memProfile, _ := os.Create("memprofile") //create memory profile log file
		stopCalls = append(stopCalls, func() {
			pprof.WriteHeapProfile(memProfile)
		})
	}
	words := loaddata(file)
	_ = getAnagrams(words)
	for _, f := range stopCalls {
		f()
	}
}

// create anagram groups from word list
func getAnagrams(ws []string) [][]string {
	amap := make(map[string][]string, len(ws))
	keys := getKeys(ws)
	for i := range ws {
		w, key := ws[i], keys[i]
		amap[key] = append(amap[key], w)
	}
	//clean up results
	var anagrams [][]string
	for _, g := range amap {
		//dedup anagram groups based on lowercase equivalance
		gLow := make([]string, len(g))
		for i := range g {
			gLow[i] = strings.ToLower(g[i])
		}
		g = dedup(g, gLow)
		//drop trivial anagram groups of 0 or 1 element
		if len(g) <= 1 {
			continue
		}
		anagrams = append(anagrams, g)
	}
	return anagrams
}

// dedup g based on equivalence of keys k
func dedup(g []string, k []string) []string {
	var dd []string //g, deduped
	for i := range g[:len(g)-1] {
		d := false
		// keep only last occurance of duplicate
		for j := i + 1; j < len(g); j++ {
			if k[i] == k[j] {
				d = true
				continue
			}
		}
		if !d {
			dd = append(dd, g[i])
		}
	}
	dd = append(dd, g[len(g)-1])
	return dd
}

// create keys so that if k1=k2 then words[1] is an anagram of words[2]
func getKeys(words []string) []string {
	lwords := make([]string, len(words))
	for i := range words {
		lwords[i] = strings.ToLower(words[i])
	}
	//make keys by bubble sort of letters
	var keys []string
	for _, w := range lwords {
		wslc := []rune(w)
		bsort(wslc)
		keys = append(keys, string(wslc))
	}
	return keys
}
func bsort(s []rune) {
	wl := len(s)
	if wl < 2 {
		return
	}
	for {
		var swap bool
		for i := range s[:wl-1] {
			if s[i] > s[i+1] {
				swap = true
				s[i], s[i+1] = s[i+1], s[i]
			}
		}
		if !swap {
			break //done sorting
		}
	}
}

// create list of words
func loaddata(f string) []string {
	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	s := string(data)
	words := strings.Fields(s)
	return words
}
