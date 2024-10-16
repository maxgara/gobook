package main

import "fmt"

func main() {
	var s = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Printf("s=%v, &s=%p, &(s[0])=%p\n", s, &s, &(s[0]))
	s = append(s[:1], s[2:]...)
	fmt.Println("append")
	fmt.Printf("s=%v, &s=%p, &(s[0])=%p\n", s, &s, &(s[0]))
}
