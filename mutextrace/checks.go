package sync

import "fmt"

// look for a stack trace where call is waiting to lock mutex already locked by earlier frame
func checkForDoubleLock(c *caller) bool {
	var set, trying *caller
	if c.flags&WAITING != 0 {
		trying = c
	}
	for c.prev != nil {
		c = c.prev[0]
		if c.flags&LOCKSET != 0 {
			set = c
		}
	}
	if set != nil && trying != nil {
		fmt.Printf("WARNING: possible Mutex error (unusual behaviour): %v waiting for Lock(), lock set by previous frame %v\n", sprintFrame(&trying.Frame, 1), sprintFrame(&set.Frame, 1))
		return true
	}
	return false
}
