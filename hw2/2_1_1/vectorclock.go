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

func IsBefore(tsA [numNodes]int, tsB [numNodes]int) bool {
	// A->B, A happens before B if every A_i <= B_i for all i \elem [0, len(A))
	// A-/>B if any A_i > B_i for all i \elem [0, len(A))
	for k := 0; k < numNodes; k++ {
		if tsA[k] > tsB[k] {
			return false
		}
	}
	return true
}
