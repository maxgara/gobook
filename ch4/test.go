package main

import "fmt"

type State int

const (
	NY State = iota
	CA
	NV
	CO
)

func main() {

	var arr = [3]string{"hey", "hi", "hello"}
	// var friends = [4]string{NY: "john", NV: "Becs", CA: "Kev", CO: "PArker"} // order fixed by enum
	fmt.Printf("%v %p\n", arr, &arr)
	for i := range arr {
		fmt.Printf("arr[%d]:%s\t%p\n", i, arr[i], &(arr[i]))
	}
	arr[1] = "whaddup"
	fmt.Printf("%v %p\n", arr, &arr)
	for i := range arr {
		fmt.Printf("arr[%d]:%s\t%p\n", i, arr[i], &(arr[i]))
	} // fmt.Println(friends)

	_ = NY
}
