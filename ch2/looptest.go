package main

import "fmt"

func main() {
	var arr = []int{6, 2, 6, 7, 4, 6, 3}
	for i, v := range arr {
		if v == 6 {
			i = i + 2
		}
		fmt.Printf("%d:%v\n", i, v)
	}
}
