package main

import "fmt"

func main() {
	s := "123456789253"
	fmt.Println(commas(s))
}

func commas(s string) string {
	var out string = ""
	var count int = 1
	for i := len(s) - 1; i >= 0; i-- {
		if count == 3 {
			out = string(s[i]) + "," + out
			count = 0
		} else {
			count++
			out = string(s[i]) + out
		}
	}

	return out
}
