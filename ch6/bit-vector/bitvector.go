package main

import "fmt"

type IntSet struct {
	words []uint64
}

func (s IntSet) Has(x int) bool {
	word := s.words[x/64]
	bit := (word >> (x % 64)) & 0b1
	return bool(bit == 0b1)
}
func getidx(i int) (int, int) {
	return i / 64, i % 64
}

func (s IntSet) String() string {
	max := len(s.words) * 64
	slc := make([]int, 0, max)
	for i := 0; i < max; i++ {
		if s.Has(i) {
			slc = append(slc, i)
		}
	}
	return fmt.Sprintf("%v\n", slc)
}
func (s *IntSet) Add(x int) {
	i, j := getidx(x)
	for i > len(s.words)-1 {
		s.words = append(s.words, 0)
	}
	mask := uint64(0b1 << j)
	s.words[i] |= mask //set correct bit to 1
}

// set s to union of s with v
func (s *IntSet) UnionWith(v *IntSet) {
	startlen := len(s.words)
	for i := 0; i < startlen; i++ {
		s.words[i] |= v.words[i]
	}
	if startlen < len(v.words) {
		s.words = append(s.words, v.words[startlen:]...)
	}
}

func main() {
	var words = []uint64{0, 16, 123, 5, 1, 1}
	var s = IntSet{words}
	fmt.Println(s)

}
