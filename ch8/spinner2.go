package main

import (
	"fmt"
	"time"
)

const fibDepthConcurrent = 7 //goroutines = 1 + 2 ^ (fibDepthConcurrent)
var goroutines int

const N = 45 //fibonacci target number

func main() {
	start := time.Now()
	go spinner(time.Second / 4)
	fibN := cfib(N)
	fmt.Printf("elapsed: %v\n", time.Now().Sub(start))
	fmt.Printf("\rfibonacci %v = %v\n", N, fibN)
	fmt.Printf("fibonacci goroutines: %v\n", goroutines)

}
func fib(n int) int {
	if n < 3 {
		return 1
	}
	return fib(n-1) + fib(n-2)
}

func cfib(n int) int {
	c := make(chan int)
	go cfibchan(n, c)
	return <-c
}

func cfibchan(n int, c chan int) {
	// fmt.Printf("getting fib %v\n", n)
	if n < 3 {
		c <- 1
	}
	if N-n >= fibDepthConcurrent {
		//drop down to normal fib once concurrency allowance is gone
		c <- fib(n)
		return
	}
	//create 2 more goroutines
	cc := make(chan int)
	go cfibchan(n-1, cc)
	go cfibchan(n-2, cc)
	goroutines += 2
	var sum int
	sum += <-cc
	sum += <-cc
	c <- sum
}

func spinner(t time.Duration) {
	for {
		for _, c := range `-\|/` {
			fmt.Printf("\r\r%c", c)
			time.Sleep(t)
		}
	}
}
