package main

import "fmt"

func main() {
	var s = []int{1, 2, 5}
	s = addInt(s, 1, 6)
	fmt.Println(s)
}

func addInt(s []int, i ...int) []int {
	fmt.Printf("mystery:%v\n", i)
	inlen := len(s)
	out := make([]int, inlen+len(i))
	copy(out[:inlen], s)
	copy(out[inlen:], i)
	return out
}

// make()
// //return new int slice with length l and capacity c
// func makeslice(l int, c int){
// 	var arr = []int
// }
