package sync

import "testing"

func TestLockUnlock(t *testing.T) {
	var m Mutex
	m.Lock()
	// m.Unlock()
}
