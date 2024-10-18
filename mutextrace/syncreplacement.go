package sync

//replace standard sync package to hook Mutex Open() and Close() calls
import (
	"fmt"
	"runtime"
	"strings"
	stdsync "sync"
)

const PRINTFRAMESFULL = 0
const PRINTSHORTFRAMES = 1

// keep track of what Mutexes are locked
type Mutex struct {
	sMutex  stdsync.Mutex
	callers []caller // where methods on the mutex have been called from
}

// represent a frame as a node in call graph
type caller struct {
	runtime.Frame
	next  []*caller //child nodes
	prev  []*caller //parent nodes
	flags int       //flags representing the kind of caller (direct, indirect, active, inactive, etc.)
}

func (m *Mutex) Lock() {
	fmt.Println("hook Lock()")
	frames := getcFrames()
	c := ftoc(frames)
	fmt.Printf("direct caller:%v\n", *c)
	if m.callers == nil {
		m.callers = []caller{}
	}
	m.callers = append(m.callers, *c)
	ac := common()
	// for len(c.prev) != 0 {
	// 	c = c.prev[0]
	// 	fmt.Print(sprintFrame(&c.Frame, 0))
	// }

	m.sMutex.Lock()
	// printFrames(f, 1)
}

func (m *Mutex) Unlock() {
	fmt.Println("Hook Unlock()")
	m.sMutex.Unlock()
}

// Walk Caller linked list backwards ***may not be needed
func forEachParent(c *caller, f func(*caller) bool) error {
	return nil
}

// convert runtime.Frames into *caller linked list
func ftoc(frames *runtime.Frames) *caller {
	c := new(caller)
	first := c // keep direct caller to return
	f, more := frames.Next()
	for {
		c.Function = f.Function
		c.Line = f.Line
		c.File = f.File
		f, more = frames.Next()
		// fmt.Println(sprintFrame())
		//stop if this was the oldest frame in the stack trace
		if !more {
			break
		}
		//attach new parent to child
		p := caller{next: []*caller{c}}
		c.prev = []*caller{&p}
		//move c pointer to parent
		c = &p
	}
	return first
}

// walk caller tree forwards *** may not be needed
func forEachChild(c *caller, f func(*caller) bool) error {
	return nil
}

// get common ancestor function. ancestor1 is in c1's trace while ancestor 2 is in c2's trace. n is depth
func common(c1 *caller, c2 *caller, n int) (a1 *caller, a2 *caller, nr int) {
	type apair struct {
		a1 *caller
		a2 *caller
		n  int
	}
	var pairs = []apair{}
	if a1.prev != nil {
		for _, v := range a1.prev {
			a1Anc, a2Anc, n1 := common(v, a2, n+1)
			pairs = append(pairs, apair{a1Anc, a2Anc, n1})
		}
	}
	if a2.prev != nil {
		for _, v := range a2.prev {
			a1Anc, a2Anc, n2 := common(a1, v, n+1)
			pairs = append(pairs, apair{a1Anc, a2Anc, n2})
		}
	}
	var minp *apair
	var minn = 10000
	for _, v := range pairs {
		if v.n < minn {
			minp = v
		}
	}
	return minp.a1, minp.a2, minn
}

// get frames for current call stack
func getcFrames() (f *runtime.Frames) {
	var pcs = make([]uintptr, 20) //program counters for calling funcs
	n := runtime.Callers(3, pcs)
	pcs = pcs[:n] //remove extra buffer space
	return runtime.CallersFrames(pcs)
}

// print a frame for debugging
func sprintFrame(f *runtime.Frame, mode int) string {
	const sep = "****************************"
	file := f.File
	line := f.Line
	function := f.Function
	pathend := strings.LastIndex(f.Function, "/")
	if pathend == -1 {
		pathend = 0
	}
	dot := strings.Index(function[pathend:], ".") + pathend
	var pack, short string
	if dot != -1 {
		pack = function[:dot]
		short = function[dot+1:]
	}
	switch mode {
	case PRINTFRAMESFULL:
		return fmt.Sprintf("%s\nFILE:%v\nLINE:%v\nFUNC:%v\nPACKAGE:%v\nSHORT:%v\n", sep, file, line, function, pack, short)
	case PRINTSHORTFRAMES:
		return fmt.Sprintf("%s\nPACKAGE:%v\nSHORT:%v\n", sep, pack, short)
	}
	return ""
}

// break a function call like github.com/x/mutextrace.getlockcallers into package path and func name values
func parseFunc(fun string) (p string, f string) {
	pathend := strings.LastIndex(fun, "/")
	if pathend == -1 {
		pathend = 0
	}
	dot := strings.Index(fun[pathend:], ".") + pathend
	if dot != -1 {
		p = fun[:dot]
		f = fun[dot+1:]
	}
	return p, f
}
