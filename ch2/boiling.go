//print boiling point of water

package main

import "fmt"

const boilingF = 212.0

func main() {
	var f = boilingF
	var c = (f - 32) * 5 / 9
	fmt.Printf("F:%f C:%f\n", f, c)
}
