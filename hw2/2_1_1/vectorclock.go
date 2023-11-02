package main

import "sync"

var (
	zeroVector = make([]int, numNodes)
)

type VectorClock struct {
	mu sync.Mutex
	ts [numNodes]int
}

func (clock *VectorClock) Inc(id int) {
	// increment ts for a particular id
	clock.mu.Lock()
	defer clock.mu.Unlock()
	clock.ts[id]++
}
func (clock *VectorClock) Get() [numNodes]int {
	// increment ts for a particular id
	clock.mu.Lock()
	defer clock.mu.Unlock()
	return clock.ts

}
