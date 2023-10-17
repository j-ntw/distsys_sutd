package main

import "fmt"

const (
	numNodes        = 10
	timeout         = 1000 // timeout to wait for messages from other nodes
	broadcast_delay = 1100 // broadcast delay period
)

func main() {
	// create nodes, channels
	var ch_arr [numNodes]chan Msg
	var node_arr [numNodes]Node

	for i := range ch_arr {
		ch_arr[i] = make(chan Msg)
		node_arr[i] = *NewNode(i)
	}

	// boot nodes, try to elect self, election done
	for i := range node_arr {
		node_arr[i].ch_arr = ch_arr
		go node_arr[i].Boot()
	}

	// run while waiting for input
	var input string

	for {
		fmt.Scanln(&input)
		switch command := input; command {
		case "exit":
			return
		}
	}

}
