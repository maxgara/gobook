package main

import (
	"fmt"
	"time"
)

const (
	MILLION = 1000000 //1,000,000
	HUNDRED = 100
)

func main() {
	start := time.Now()
	var primes = []int{2, 3}
	var div bool
	for i := 3; i < MILLION; i++ {
		div = false
		for _, p := range primes {
			if i%p == 0 {
				div = true
				break
			}
		}
		if !div {
			primes = append(primes, i)
		}
	}
	secs := time.Since(start).Seconds()
	fmt.Println(primes)
	fmt.Printf("secs: %.2f\n", secs)
}
