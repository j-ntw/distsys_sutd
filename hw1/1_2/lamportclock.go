package main

import (
	"fmt"
	"sync"
)

type LamportClock struct {
	mu sync.Mutex
	ts int
}

func (clock *LamportClock) AdjustClock(id int, ts int, msg_ts int) {
	clock.mu.Lock()
	if msg_ts > ts {
		fmt.Printf("adjust clock_%d: %d->%d\n", id, ts, msg_ts)
		clock.ts = msg_ts

	} else {
		fmt.Printf("adjust clock_%d: %d->%d\n", id, ts, ts)
		clock.ts = ts
	}
	clock.mu.Unlock()

}

func (clock *LamportClock) Inc() {
	clock.mu.Lock()
	clock.ts++
	clock.mu.Unlock()
}
