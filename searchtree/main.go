package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime/pprof"
	"strings"
	"time"
	"unicode"
)

const persistData = true   //save data in file
const toksEnabled = true   //tokenize data in files for fast content search
const cacheEnabled = false //don't load previous cache

var fileList []string
var tokens map[string]*Tok
var stop bool
var pause bool
var allclear chan bool

type Tok struct {
	s   string
	ref []int //idxs in data buckets
}

var r *regexp.Regexp

func main() {
	//setup
	f, _ := os.Create("cpuprof")
	pprof.StartCPUProfile(f)
	//r = regexp.MustCompile(`\.(txt|log|json|html|conf|go|c|xml)$`)
	tokens = make(map[string]*Tok)
	allclear = make(chan bool)

	// get data if cache not enabled or can't load cached data
	if !cacheEnabled || !loadIdxData() {
		//continuously update saved data files and print updates every 5 sec
		go checkForStopSig()
		if persistData {
			go indexStateUpdate()
		}
		buildSearchDB("/Users/maxgara/") //****do actual work
	}
	pprof.StopCPUProfile() //stop profiling after building data structs
	fmem, _ := os.Create("memprof")
	pprof.WriteHeapProfile(fmem)
	fmt.Printf("enter search:")
	//get search input
	for {
		buf := make([]byte, 1000)
		n, _ := os.Stdin.Read(buf)
		matches := searchPrim(string(buf[:n]))
		fmt.Printf("%v\n[%d matches]", matches, len(matches))
	}
}

// generate data for search -recursive
func buildSearchDB(dir string) {
	if stop {
		//user requested stop - move on to search phase
		fmt.Println("stopped")
		return
	}
	//pause work to print current map state
	if pause {
		allclear <- true
		<-allclear
	}
	fs, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, f := range fs {
		full := dir + "/" + f.Name()
		fileList = append(fileList, full)
		if toksEnabled {
			//if r.Match([]byte(full)) {
			if strings.HasSuffix(full, ".txt") {
				indexContents(full, len(fileList))
			}
		}
		if f.Type().IsDir() {
			buildSearchDB(full)
		}
	}
}

// store given idx in ref list for all tokens contained in f
func indexContents(f string, idx int) {
	b, err := os.ReadFile(f)
	if err != nil {
		return
	}
	ff := func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r) //split on any non-letter
	}
	ts := strings.FieldsFunc(string(b), ff) // fields of file contents
	for _, t := range ts {
		p, ok := tokens[t]
		if !ok {
			tokens[t] = &Tok{s: t}
			continue
		}
		p.ref = append(p.ref, idx)
	}
}

// primative free text search
func searchPrim(key string) []string {
	key = strings.Trim(key, "\n ")
	if len(key) < 3 {
		fmt.Println("search too short")
		return nil
	}

	p, ok := tokens[key]
	if !ok {
		return nil
	}
	idxs := p.ref
	var out []string
	for _, ix := range idxs {
		out = append(out, fileList[ix])
	}
	return out
}

// load fileList and tokens from file
func loadIdxData() bool {
	//get file list
	sfs, err := os.ReadFile("list.txt")
	if err != nil {
		return false
	}
	fileList = strings.Split(string(sfs), "\n")
	fmt.Println("loaded file list from save")
	//get token list
	// sts, err := os.ReadFile("toks.txt")
	// if err != nil {
	// 	return false
	// }
	// tokenPropList := strings.Split(string(sfs), "\n")
	// for tp := range tokenPropList {
	// 	tpset := strings.Split(tp, " ")
	// 	t := tpset[0]
	// 	ps := tpset[1:]

	// }
	return true
}

func checkForStopSig() {
	fmt.Printf("listening for stop...\n")
	d := make([]byte, 1)
	os.Stdin.Read(d)
	stop = true
	fmt.Printf("stopping..")
}

// save file list, print indexing update on regular interval
func indexStateUpdate() {
	fmt.Println("indexstateupdate running")
	//print update every 5 sec
	for {
		if stop {
			return
		}
		<-time.After(time.Second * 5)
		pause = true //request stop work on map
		<-allclear   //wait for confirmation of work stop
		f, err := os.Create("list.txt")
		if err != nil {
			log.Fatal(err)
		}
		for _, s := range fileList {
			f.WriteString(s + "\n")
			//fmt.Printf("filename %d =%v\n", i, s)
		}
		fmt.Printf("files:%d\n", len(fileList))
		f.Close()
		f, err = os.Create("toks.txt")
		if err != nil {
			log.Fatal(err)
		}
		for _, p := range tokens {
			s := fmt.Sprintf("%v [%d entries]\n", p.s, len(p.ref))
			f.WriteString(s)
		}
		fmt.Printf("toks:%d\n", len(tokens))
		f.Close()
		pause = false
		allclear <- true
	}
}
