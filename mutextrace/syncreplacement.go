package sync

//replace standard sync package to hook Mutex Open() and Close() calls
import (
	"fmt"
	"runtime"
	stdsync "sync"
)

type Mutex struct {
	sMutex stdsync.Mutex
}

// type sMutex stdsync.Mutex

func (m *Mutex) Lock() {
	fmt.Println("hook Lock()")
	m.sMutex.Lock()
	var pcs = []uintptr{}
	n := runtime.Callers(0, pcs)
	fmt.Printf("Callers: %v\n", n)
	for i, v := range pcs {
		fmt.Printf("[%v]:%v\n", i, v)
	}
}

func (m *Mutex) Unlock() {
	fmt.Println("Hook Unlock()")
	m.sMutex.Unlock()
}
