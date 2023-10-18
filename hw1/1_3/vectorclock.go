package main

import "sync"

type VectorClock struct {
	mu sync.Mutex
	ts [numEntities]int
}

func (clock *VectorClock) isCV(msg_ts [numEntities]int) bool {

	clock.mu.Lock()
	defer clock.mu.Unlock()
	// check if before
	return !IsBefore(clock.ts, msg_ts)
}

func (clock *VectorClock) AdjustClock(ts [numEntities]int, msg_ts [numEntities]int) {
	clock.mu.Lock()
	defer clock.mu.Unlock()
	// element wise comparison/swap of ts
	for i := 0; i < numEntities; i++ {
		if msg_ts[i] > ts[i] {
			clock.ts[i] = msg_ts[i]

		} else {
			clock.ts[i] = ts[i]
		}
	}
	// fmt.Printf("\nadjust clock:\n%v\n%v\n\n", ts, clock.ts)

}

func (clock *VectorClock) Inc(id int) {
	// increment ts for a particular id
	clock.mu.Lock()
	defer clock.mu.Unlock()
	clock.ts[id]++

}
