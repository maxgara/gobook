package main

// import (
// 	"fmt"
// 	"os"
// 	"testing"
// )

// func TestWebBase(t *testing.T) {
// 	f, err := os.OpenFile("temp.html", os.O_TRUNC|os.O_WRONLY, 0)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	w := newDoc("wbstyle.css")
// 	g := w.startGrid(2)
// 	g.add("hello")
// 	g.add("world")
// 	g.add("what's")
// 	g.add("up")
// 	w.render(f)
// 	fmt.Println(*w)
// 	w.render(os.Stdout)
// }
