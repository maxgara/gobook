package main

import "fmt"

func main() {
	s := "hey! こんに"
	fmt.Printf("str_rep\t\t%x\n", s)
	fmt.Printf("rune_rep\t%x\n\n", []rune(s))
	fmt.Printf("rune_rep_v2\t")
	r := []rune(s)
	for _, ru := range r {
		fmt.Printf("%.4x", ru)
	}
	fmt.Println()
	fmt.Println("s_idx\trune\tcharacter")
	for i, r := range s {
		fmt.Printf("s[%d]\t%x\t%c\n", i, r, r)
	}
	fmt.Printf("%s\n", s)
}
