package main

import (
	"SleepRand"
	"fmt"
	"math/rand"
)

const (
	numNodes = 10
)

func client(ch_arr [numNodes]chan Msg, id int, isCoordinator bool) {
	self_ch := ch_arr[id]
	replica := id
	for {
		// sync
		if isCoordinator {
			//broadcast replica periodically
			for i, other_ch := range ch_arr {
				// todo timeout
				out_msg := Msg{normal, id, i, 0, replica}
				other_ch <- out_msg
			}
			SleepRand.SleepRand()
		} else {
			// receive replica
			// todo timeout
			in_msg := <-self_ch
			replica = in_msg.data
		}

	}
}
func main() {
	var ch_arr [numNodes]chan Msg

	fmt.Println("create nodes")
	// select a random client as coordinator
	startingCoordinatorId := rand.Intn(numNodes)
	for i := range ch_arr {
		// make a channel of type Msg
		// add ch to array
		ch_arr[i] = make(chan Msg)
		isCoordinator := (i == startingCoordinatorId) // false unless they are actually the selected coordinator, then true
		go client(ch_arr, i, isCoordinator)

	}

	var input string
	fmt.Scanln(&input)

}
