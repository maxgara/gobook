package main

import "fmt"

func main() {

	strcount("teststr")
}

func strcount(s string) {
	// var a rune = 'a'
	for i := 0; i < 26; i++ {
		fmt.Printf("%c%c", 'a'+rune(i), 'A'+rune(i))
	}
}
