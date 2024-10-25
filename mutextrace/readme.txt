The purpose of this package is to create a Mutex type that is able to track callers and identify conflicts.
The current approach is to create a new "Mutex" object in this package which wraps a regular sync.Mutex.
this lets us:
-get trace for functions calling m.Lock(), m.Unlock() on a given Mutex
-build graph of frames which call these functions, and their ancestors
-analyze graph for odd behaviour (such as frame trying to Lock mutex already locked by parent frame)

This may be enough to find most Mutex problems

downsides to this approach:
- Unable to create a graph of possible Mutex callers, and therefore also:
-can’t fully predict if an apparent stuck lock condition is really stuck, may be execution paths to unlock
-currently requires package user to modify the code they want to debug/analyze
-can’t do anything useful until runtime

Goal of this step: simple checks:
    1. write to stdin how many calls are waiting, who last set lock
    2. write a warning if frame sets a lock then later frame in the same stack trace tries to set the same lock
    3. create a tree of execution flow


future plans: 
-add main package which substitutes sync libraries automatically (easy)
-look into gathering more info from go execution traces 
-look into getting gc refs to mu
-goal: eventually do analysis before runtime to find paths to unlock, then use at runtime to see when no paths remain
-produce a graph of possible lock and unlock calls by frame and goroutine