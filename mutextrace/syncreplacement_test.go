package sync

import (
	"fmt"
	"testing"
)

func TestLockUnlock(t *testing.T) {
	var m Mutex
	m.Lock()
	// m.Unlock()
}
func TestSprintFrame(t *testing.T) {
	c := getcallers()
	s := SprintFrame(&c.Frame, 0)
	fmt.Print(s)
}
func TestRprint(t *testing.T) {
	getcallers().Rprint()
}
