package main

import "sync"

var (
	zeroVector = make([]int, numProcesses)
)

type VectorClock struct {
	sync.Mutex
	ts [numProcesses]int
}

func (clock *VectorClock) Inc(id int) {
	// increment ts for a particular id
	clock.Lock()
	defer clock.Unlock()
	clock.ts[id]++
}
func (clock *VectorClock) Get() [numProcesses]int {
	// increment ts for a particular id
	clock.Lock()
	defer clock.Unlock()
	return clock.ts

}

func (clock *VectorClock) AdjustClock(msg_ts [numProcesses]int) {
	clock.Lock()
	defer clock.Unlock()
	// element wise comparison/swap of ts
	for i := 0; i < numProcesses; i++ {
		if msg_ts[i] > clock.ts[i] {
			clock.ts[i] = msg_ts[i]

		}
	}
}

func IsBefore(msgA Msg, msgB Msg) bool {
	// compare 2 messages by vector clocks
	// credit: Sean Yap
	// A->B, A happens before B if every A_i <= B_i for all i \elem [0, len(A))
	// A-/>B if any A_i > B_i for all i \elem [0, len(A))
	sumA := 0
	sumB := 0
	// A_bigger_count := 0
	for k := 0; k < numProcesses; k++ {
		sumA += msgA.ts[k]
		sumB += msgB.ts[k]
	}
	if sumA == sumB {
		// equal vector clocks, tie break via node id
		return msgA.from < msgB.from
	} else {
		// return A_bigger_count == len(msgA.ts)
		return sumA < sumB
	}
}
