package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"strings"
)

var file = "../british-english-insane.txt"

func main() {
	// debug.SetGCPercent(300)
	// go func() {
	// 	<-time.After(10 * time.Second)
	// 	panic("time's up!")
	// 	}()
	//cpu/memory profiling
	var pprofCPU = flag.Bool("cpu", false, "write Cpu profile")
	var pprofMem = flag.Bool("mem", false, "write memory profile")
	flag.Parse()
	if *pprofCPU {
		profile, _ := os.Create("cpuprofile") //create cpu profile log file
		pprof.StartCPUProfile(profile)
		defer func() {
			pprof.StopCPUProfile()
		}()
	}
	if *pprofMem {
		profile, _ := os.Create("memprofile") //create cpu profile log file
		defer func() {
			runtime.GC()
			pprof.WriteHeapProfile(profile)
		}()
	}
	words := loaddata(file)
	_ = getAnagrams(words)
	// fmt.Println(out)
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
	keys := make([]string, 0, len(words))
	//****init key
	// -
	//** keygen
	for _, w := range lwords {
		key := runefieldkey([]rune(w))
		keys = append(keys, key)
	}
	return keys
}

func bytefieldkey(w []rune, key []byte) {
	byte_occurances := key
	for _, r := range w {
		byte_occurances[r]++
	}
}

func runefieldkey(w []rune) string {
	occurances := make([]byte, 0xFF)
	for _, r := range w {
		occurances[r]++
	}
	return string(occurances)
}

func runeslicekey(w []rune) []rune {
	var rsym []rune //runes found in word
	var rn []rune   //count for runes in rsym
	for _, r := range w {
		var rfound bool
		for i := range rsym {
			if rsym[i] == r {
				rfound = true
				rn[i]++ //rune previously found case, increment existing count
			}
		}
		if !rfound {
			rsym = append(rsym, r)
			rn = append(rn, 1)
		}
	}
	return append(rsym, rn...)
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
