package main

import (
	"fmt"
	"testing"
)

var testdata = []string{"aab", "aba", "baa", "bba", "bab"}

// func TestComp(t *testing.T) {
// 	dl := len(testdata)
// 	var groups [][]string
// 	for i, s1 := range testdata[:dl-1] {
// 		fmt.Printf("i=%v\n", i)
// 		var group []string
// 		for _, s2 := range testdata[i+1:] {

// 			if comp(s1, s2) {
// 				group = append(group, s2)
// 			}
// 		}
// 		if len(group) != 0 {
// 			group = append(group, s1)
// 			groups = append(groups, group)
// 		}
// 	}
// 	fmt.Println(groups)
// }

//	func TestRun(t *testing.T) {
//		run(testdata)
//	}

func TestGetAnagrams(t *testing.T) {
	anagrams := getAnagrams(testdata)
	fmt.Println(anagrams)
}
