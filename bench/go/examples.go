package bench

const COUNT = 1000000
const LOOPS = 1000

func f1() {
	var counts [LOOPS]uint64
	for i := 0; i < LOOPS; i++ {
		counts[i] += countevens(i) / COUNT
	}
	var cc uint64
	for i := 0; i < LOOPS; i++ {
		cc += counts[i] - LOOPS/2
	}
}
func countevens(x int) uint64 {
	var evens [COUNT]uint
	for i := 0; i < COUNT; i++ {
		evens[i] = uint(iseven(x + i))
	}
	var counter uint64 = 0
	for i := 0; i < COUNT; i++ {
		counter += uint64(evens[i])
	}
	return counter
}
func iseven(x int) int {
	var y int
	if x < 10 {
		y = x + 4
	} else {
		y = x
	}
	z := y % 2
	if z == 0 {
		return 1
	}
	return 0
}
