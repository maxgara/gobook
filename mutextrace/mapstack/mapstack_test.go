package mapstack

import "testing"

func TestRun(t *testing.T) {

}
func Test1(t *testing.T) {
	vlookup := make(map[string]*Vertex) //collection of all frames and their addresses
	f0(vlookup)
	if len(vlookup) == 0 {
		t.Fail()
	}
	printmap(vlookup)
}
func f0(l map[string]*Vertex) {
	addstack(l)
}
func f1(l map[string]*Vertex, n int) {
	addstack(l)
	if n < 3 {
		f1(l, n+1)
	}
}
