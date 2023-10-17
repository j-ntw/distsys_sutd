package main

import (
	"fmt"
)

const (
	numNodes = 10
	timeout  = 1000
	period   = 500
)

func main() {
	// creates nodes
	var ch_arr [numNodes]chan Msg
	var node_arr [numNodes]Node

	// make channels
	
	for i := range ch_arr {
		ch_arr[i] = make(chan Msg)
		node_arr[i] = *NewNode(i)
		node_arr[i].cmd = make(chan Command)
		node_arr[i].trigger_ch = make(chan bool)
	}

	// boot nodes, try to elect self, election done
	for i := range node_arr {
		node_arr[i].ch_arr = ch_arr
		node_arr[i].Boot()
	}
	// start syncing

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
