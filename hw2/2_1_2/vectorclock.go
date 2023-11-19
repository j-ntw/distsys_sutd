package main

import "sync"

var (
	zeroVector = make([]int, numNodes)
)

type VectorClock struct {
	sync.Mutex
	ts [numNodes]int
}

func (clock *VectorClock) Inc(id int) {
	// increment ts for a particular id
	clock.Lock()
	defer clock.Unlock()
	clock.ts[id]++
}
func (clock *VectorClock) Get() [numNodes]int {
	// increment ts for a particular id
	clock.Lock()
	defer clock.Unlock()
	return clock.ts

}

func (clock *VectorClock) AdjustClock(ts [numNodes]int, msg_ts [numNodes]int) {
	clock.Lock()
	defer clock.Unlock()
	// element wise comparison/swap of ts
	for i := 0; i < numNodes; i++ {
		if msg_ts[i] > ts[i] {
			clock.ts[i] = msg_ts[i]

		} else {
			clock.ts[i] = ts[i]
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
	for k := 0; k < numNodes; k++ {
		// if msgA.ts[k] > msgB.ts[k] {
		// 	A_bigger_count++
		// }
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

// A [1,1,1,0]
// B [0,2,0,0]
