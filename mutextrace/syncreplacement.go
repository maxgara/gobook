package sync

//replace standard sync package to hook Mutex Open() and Close() calls
import (
	"fmt"
	"runtime"
	"strings"
	stdsync "sync"
)

// keep track of what Mutexes are locked
type Mutex struct {
	sMutex  stdsync.Mutex
	callers []*caller
}
type caller struct {
	runtime.Frame
	Pack  string
	Short string
	next  *caller
	prev  *caller
}

func (m *Mutex) Lock() {
	fmt.Println("hook Lock()")
	c := getcallers()
	if len(m.callers) == 0 {
		m.callers = []*caller{}
	}
	m.callers = append(m.callers, c)
	c.Rprint()
	m.sMutex.Lock()
	// printFrames(f, 1)
}

func (m *Mutex) Unlock() {
	fmt.Println("Hook Unlock()")
	m.sMutex.Unlock()
}

// return direct Lock() *caller with linked prev callers
func getcallers() (c *caller) {
	c = new(caller)
	var pcs = make([]uintptr, 20) //program counters for calling funcs
	n := runtime.Callers(1, pcs)
	pcs = pcs[:n] //remove extra buffer space
	frames := runtime.CallersFrames(pcs)
	var cc *caller
	cc = c
	for {
		f, more := frames.Next()
		cc.Function = f.Function //includes package path, we will break it up
		cc.Line = f.Line
		pathend := strings.LastIndex(c.Function, "/")
		if pathend == -1 {
			pathend = 0
		}
		dot := strings.Index(c.Function[pathend:], ".")
		if dot != -1 {
			cc.Pack = cc.Function[:dot]
			cc.Short = cc.Function[dot:]
		}

		if !more {
			break
		}
		//step back one frame
		current := cc
		prev := &caller{}
		cc.prev = prev
		cc = prev
		cc.next = current
	}
	return c
}

const PRINTFRAMESFULL = 0
const PRINTSHORTFRAMES = 1

// print frames to stdout in a readable format
func printFrames(frames *runtime.Frames, mode int) {
	const sep = "****************************"
	// var pcs = make([]uintptr, 20)
	// n := runtime.Callers(0, pcs)
	// pcs = pcs[:n]
	// frames := runtime.CallersFrames(pcs)
	//print frames
	for {
		f, more := frames.Next()
		fmt.Printf(SprintFrame(&f, mode))
		if !more {
			break
		}
	}
	fmt.Println(sep)
}

// print one frame in a nice format
func SprintFrame(f *runtime.Frame, mode int) string {
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

// recursively print callers
func (c *caller) Rprint() {
	str := SprintFrame(&c.Frame, 0)
	fmt.Print(str)
	if c.prev == nil {
		return
	}
	c.prev.Rprint()
}
