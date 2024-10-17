package main

import (
	"fmt"
	"runtime/debug"
	"sync"
)

//find Mutex problems in running program causing deadlock.

func main() {
	var m sync.Mutex
	// m = test2()
	m.Lock()
	m.Lock()

	// go test(&m)
	// time.Sleep(time.Second)
	// m.Lock()
	// fmt.Printf("auth called. Stack:\n%s\n", debug.Stack())
}

//	func Watch(mu *sync.Mutex) sync.Mutex {
//		var m sync.Mutex
//	}
func test(m *sync.Mutex) {

	m.Lock()
	for {
		fmt.Printf("Stack:\n%s\n", debug.Stack())
	}
}

func (m M) Lock() {
	fmt.Printf("Fake Lock() called!!!\n")
}

type Sync struct {
	Mutex interface {
		Lock()
	}
}

// implements interface "Mutex"
type M struct {
	b bool
}

func test2() {
	var sync Sync
	sync.Mutex = M
	var m sync.Mutex
	s.Lock()
	{
	}
	return sync.Mutex(m.Mutex)
}

//*possibilities:
//- create alternative "sync.Mutex" package. **Current approach**
//-
