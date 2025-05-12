package main

import (
	"fmt"
	"os"
	"strings"
)

type Tok struct {
	s    string
	chl  []*Tok //larger tokens where this token is a substring - should be organized to keep search steps low
	par  *Tok
	fake bool
}

const first1000 = true

func main() {
	createDB()
	root := Tok{}
	fmt.Println("vim-go")
	dict := loadDict()
	for _, w := range dict {
		addWord(w, &root)
	}
	//ttprint(&root, "#root")
	search("test", &root)
	searchf("test", &root)
	profile(&root)
	//ttprint(&root, "$", 4)
	for {
		buf := make([]byte, 1000)
		n, _ := os.Stdin.Read(buf)
		search(string(buf[:n]), &root)
		searchf(string(buf[:n]), &root)
	}
}

// double-link parent to child
func claimChild(p, c *Tok) {
	p.chl = append(p.chl, c)
	c.par = p
}

// double-unlink parent and child
func disownChild(p, c *Tok) {
	l := len(p.chl)
	if l == 0 {
		fmt.Printf("error: disownChild: parent %v has no children\n", p.s)
	}
	for i, t := range p.chl[:l-1] {
		if t == c {
			copy(p.chl[i:], p.chl[i+1:])
		}
	}
	p.chl = p.chl[:l-1]
	c.par = nil //will need to change this if multiple parents are allowed
}

// sort word into tree based on relation [<], where a<b if a is a substring of b. then b is a descendent of a, or in a different branch
// return true if w was added a descendent of root.
func addWord(w string, root *Tok) bool {
	//fmt.Printf("adding %v at root=%v\n", w, root.s)
	//ttprint(root, "$", 2)
	if len(w) == 0 {
		fmt.Println("0 len word skip")
		return true
	}
	// add replace root node with w, re-add root as child of w
	if strings.Contains(root.s, w) {
		if root.par == nil {
			fmt.Printf("fail at root=%v\n", root.s)
		}
		nt := Tok{s: w}
		claimChild(root.par, &nt)
		disownChild(root.par, root)
		claimChild(&nt, root)
		//look for old siblings of root to claim as children after initial swap
		//		for _, t := range nt.par.chl {
		//			if strings.Contains(t.s, nt.s) {
		//				//TODO finish this
		//				disownChild(t.par, t)
		//				claimChild(&nt, t)
		//			}
		//		}
		return true
	}
	// if root !< w => cannot add
	if !strings.Contains(w, root.s) {
		return false
	}
	// try to add as a deeper desc.
	// if too many children, fill in fake nodes
	if len(root.chl) > 100 {
		reorg(root)
	}
	for _, c := range root.chl {
		if addWord(w, c) {
			return true
		}
	}
	// otherwise add w as root.w
	root.chl = append(root.chl, &Tok{s: w, par: root})
	return true
}

// reorganize children so that search path is shorter
func reorg(root *Tok) {
	fmt.Printf("triggered reorg")
	alph := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bins := []Tok{}
	for _, r := range alph {
		test := root.s + string(r)
		var hit bool
		for _, t := range root.chl {
			if !hit && strings.Contains(t.s, test) {
				hit = true
				bins = append(bins, Tok{s: test})
			}
			if strings.Contains(t.s, test) {
				disownChild(t.par, t)
				l := len(bins)
				claimChild(&bins[l-1], t)
			}
		}
	}
	for _, b := range bins {
		claimChild(root, &b)
	}
}

// depth-first walk until reaching a node n where f(n) == true
// returns false if f(n) == false for all n <= root
func walk(root *Tok, f func(*Tok) bool) bool {
	if f(root) {
		return true
	}
	for _, n := range root.chl {
		if walk(n, f) {
			return true
		}
	}
	return false
}

// get stats about tree
func profile(root *Tok) {
	levels := make([]int, 50) //count toks by depth
	var l int                 //current level
	lstat(root, &l, &levels)
	for i := range levels {
		if levels[i] == 0 {
			continue
		}
		fmt.Printf("level %d: %d\n", i, levels[i])
	}
}

func lstat(t *Tok, l *int, levels *[]int) {
	(*levels)[*l]++
	*l++
	for _, c := range t.chl {
		lstat(c, l, levels)
	}
	*l--
}

// find exact term match for w (fast)
func searchf(w string, root *Tok) bool {
	w = strings.Trim(w, "\n ")
	fmt.Printf("fast searching for [%v]\n", w)
	path := ""
	count := 0
	if rsf(w, root, &path, &count) {
		fmt.Printf("found match at %v !\nsteps=%v\n", path, count) //print only at end of walkback
		return true
	}
	fmt.Printf("no match found\n")
	return false
}

// find exact term match for w
func search(w string, root *Tok) bool {
	w = strings.Trim(w, "\n ")
	fmt.Printf("searching for [%v]\n", w)
	path := ""
	count := 0
	if rs(w, root, &path, &count) {
		fmt.Printf("found match at %v !\nsteps=%v\n", path, count) //print only at end of walkback
		return true
	}
	fmt.Printf("no match found\n")
	return false
}

func rs(w string, root *Tok, path *string, count *int) bool {
	*count++
	// if current node is match, start writing path from bottom.
	if w == root.s {
		*path = root.s
		return true
	}
	// if desc node is match, add current node to path and return
	for _, c := range root.chl {
		if rs(w, c, path, count) {
			*path = root.s + "." + *path
			return true
		}
	}
	return false
}

func rsf(w string, root *Tok, path *string, count *int) bool {
	fmt.Printf("%v ", *path+"."+root.s)
	*count++
	// if current node is match, start writing path from bottom.
	if w == root.s {
		*path = root.s
		fmt.Printf("FOUND\n")
		return true
	}
	if !strings.Contains(w, root.s) {
		fmt.Printf("MISS\n")
		return false
	}
	fmt.Printf("OK\n")
	// if desc node is match, add current node to path and return
	for _, c := range root.chl {
		if rsf(w, c, path, count) {
			*path = root.s + "." + *path
			return true
		}
	}
	fmt.Printf("BACKTRACK\n")
	return false
}

// load a lot of words for testing
func loadDict() []string {
	bytes, _ := os.ReadFile("british-english-insane.txt")
	words := strings.Split(string(bytes), "\n")
	if first1000 {
		words = words[:1000]
	}
	return words
}

func ttprint(t *Tok, prefix string, dep int) {
	if dep < 1 {
		return
	}
	full := prefix + "." + t.s
	fmt.Println(full)
	for _, c := range t.chl {
		ttprint(c, full, dep-1)
	}
}
