package sync

import (
	"runtime"
	"testing"
)

func TestLockUnlock(t *testing.T) {
	var m Mutex
	m.Lock()
	// m.Unlock()
}
func R(n int) *runtime.Frames {
	if n < 10 {
		return R(n + 1)
	}
	return getcFrames()
}

// func TestGetcFrames(t *testing.T) {
// 	frames := R(0)
// 	for f, m := frames.Next(); m; f, m = frames.Next() {
// 		// fmt.Print(sprintFrame(&f, 0))
// 	}
// }

// func TestSprintFrame(t *testing.T) {
// 	c := getcallers()
// 	s := SprintFrame(&c.Frame, 0)
// 	fmt.Print(s)
// }
// func TestRprint(t *testing.T) {
// 	getcallers().Rprint()
// }
