package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const file = "../british-english-insane.txt"

func main() {
	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	s := string(data)
	s = strings.ToLower(s)
	cc := make(map[rune]int)
	for _, v := range s {
		cc[v]++
	}
	for c, n := range cc {
		if n == 0 {
			continue
		}
		fmt.Printf("%#U:%d\n", c, n)
	}
	// fmt.Println(lookup)
	// <-time.After(5 * time.Second)

}
