package main

import (
	"fmt"
	"time"
)

const (
	numNodes = 10
	timeout  = 1000
	period   = 500
)

func node(ch_arr [numNodes]chan Msg, id int, isCoordinator bool) {

	// initialise
	self_ch := ch_arr[id]
	replica := id
	mode := Mode{}
	mode.mu.Lock()
	mode.mode = normal
	mode.mu.Unlock()
	// start election procedure on boot

	// sync
	go func() {
		for {
			if isCoordinator {
				//broadcast replica periodically
				for i, other_ch := range ch_arr {
					// todo how to detect liveness w timeout when sending?
					out_msg := Msg{normal, id, i, 0, replica}
					other_ch <- out_msg
				}

				// sleep periodically
				time.Sleep(period * time.Millisecond)
			} else {
				// receive replica
				// todo timeout
				select {
				case in_msg := <-self_ch:

					switch msg_mode := in_msg.Msgtype; msg_mode {
					case normal:
						replica = in_msg.data
					case election:
						Bully() // >:)
					// todo other msgtypes

					default:
						fmt.Printf("msg_mode: %v", msg_mode)
					}

				case <-time.After(timeout * time.Millisecond):
					// start election
					Bully()
					// no default
				}
			}
		}
	}()

}

func main() {
	// creates nodes
	var ch_arr [numNodes]chan Msg
	var node_arr [numNodes]Node
	fmt.Println("create nodes")

	// create nodes
	// boot nodes, try to elect self, election done
	// start syncing

	for i := range ch_arr {
		ch_arr[i] = make(chan Msg)
		node_arr[i].self_ch = ch_arr[i]
		// go node(ch_arr, i, isCoordinator)
	}
	for i := range node_arr {
		node_arr[i].ch_arr = ch_arr
		node_arr[i].Elect()
	}
	// launch a monitor
	// let run for awhile
	// todo kill coordinator

	var input string
	go func() {
		for {
			fmt.Scanln(&input)
			switch command := input; command {
			case "exit":
				return
			}
		}
	}()

}
