package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/pprof"
	"strings"
)

var file = "../british-english-insane.txt"
var out io.Writer
var dout io.Writer //if debug logging on
var d io.Writer

func main() {
	//cpu/memory profiling
	var pprofCPU = flag.Bool("profilecpu", false, "write Cpu profile")
	var pprofMem = flag.Bool("profilememory", false, "write memory profile")
	var outFile = flag.String("o", "", "write anagrams to file")
	var inFile = flag.String("f", "", "read dictionary from text file")
	var debugLog = flag.Bool("v", false, "turn on debug level messages and logging. written to stdout or -f file if provided")
	flag.Parse()
	pStop := initAndProf(*pprofCPU, *pprofMem, *debugLog, *outFile, *inFile)
	words := loaddata(file)
	fmt.Printf("dict size (words):%v\n", len(words))
	gs := getAnagramsWithLogs(words)
	fmt.Printf("anagram count: %v groups\n", len(gs))
	if *outFile != "" {
		f, err := os.Create(*outFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "couldn't write anagram results: %v\n", err)
		}
		fmt.Fprint(f, gs)
	}
	// fmt.Println(gs[:min(100, len(gs))])
	pStop()
}

// create anagram groups from word list
func getAnagramsWithLogs(ws []string) [][]string {
	amap := make(map[string][]string, len(ws))
	keys := getKeys(ws)
	for i := range ws {
		w, key := ws[i], keys[i]
		amap[key] = append(amap[key], w)
	}
	//clean up results
	var anagrams [][]string
	for _, g := range amap {
		//dedup anagram groups based on lowercase equivalancy
		gLow := make([]string, len(g))
		for i := range g {
			gLow[i] = strings.ToLower(g[i])
		}
		var unDup []string //anagram group, deduped
		var dups []string
		for i := range g[:len(g)-1] {
			d := false
			for j := i + 1; j < len(g); j++ {
				if gLow[i] == gLow[j] {
					d = true
					continue // keep only last occurance of duplicate word
				}
			}
			if !d {
				unDup = append(unDup, g[i])
			} else {
				dups = append(dups, g[i])
			}
		}
		unDup = append(unDup, g[len(g)-1])
		if len(dups) != 0 {
			fmt.Printf("duplicate words removed from anagram group:%v in %v -> %v\n", dups, g, unDup)
		}
		g = unDup
		//drop trivial anagram groups of 0 or 1 element
		if len(g) <= 1 {
			continue
		}
		anagrams = append(anagrams, unDup)
	}
	return anagrams
}

// create keys so that if k1=k2 then words[1] is an anagram of words[2]
func getKeys(words []string) []string {
	//normalize to lowercase
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
func initAndProf(pprofCPU, pprofMem, debugLog bool, outFile, inFile string) (stop func()) {
	var stopCalls []func()
	if pprofCPU {
		cpuProfile, _ := os.Create("cpuprofile") //create cpu profile log file
		pprof.StartCPUProfile(cpuProfile)
		stopCalls = append(stopCalls, func() {
			pprof.StopCPUProfile()
		})
	}
	if pprofMem {
		memProfile, _ := os.Create("memprofile") //create memory profile log file
		stopCalls = append(stopCalls, func() {
			pprof.WriteHeapProfile(memProfile)
		})
	}
	if inFile != "" {
		file = inFile
	}
	if outFile != "" {
		var err error
		out, err = os.Create(outFile)
		if err != nil {
			log.Fatalf("outfile: %v\n", err)
		}
	}
	if debugLog {
		if outFile != "" {
			dout = out
		}
		dout = os.Stdout
	}
	return func() {
		for _, f := range stopCalls {
			f()
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
