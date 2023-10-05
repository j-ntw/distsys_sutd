package main

import "sync"

type VectorClock struct {
	mu sync.Mutex
	ts [numEntities]int
}

func (clock *VectorClock) isCV(msg_ts [numEntities]int) bool {
	// locks clock mutex
	clock.mu.Lock()
	// check if before
	notCV := IsBefore(clock.ts, msg_ts)
	clock.mu.Unlock()
	return !notCV

}

func (clock *VectorClock) AdjustClock(ts [numEntities]int, msg_ts [numEntities]int) {
	clock.mu.Lock()

	// element wise comparison/swap of ts
	for i := 0; i < numEntities; i++ {
		if msg_ts[i] > ts[i] {
			clock.ts[i] = msg_ts[i]

		} else {
			clock.ts[i] = ts[i]
		}
	}
	// fmt.Printf("\nadjust clock:\n%v\n%v\n\n", ts, clock.ts)
	clock.mu.Unlock()
}

func (clock *VectorClock) Inc(id int) {
	// increment ts for a particular id
	clock.mu.Lock()
	clock.ts[id]++
	clock.mu.Unlock()
}
