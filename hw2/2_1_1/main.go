package main

import "fmt"

// implement lamports shared priority queue with vector clock

const (
	numNodes = 3
)

var ch_arr [numNodes]chan Msg
var node_arr [numNodes]Node

// all machines are connected to all other machines (use channels)
func main() {
	// create nodes, channels

	for i := range ch_arr {
		ch_arr[i] = make(chan Msg)
		node_arr[i] = *NewNode(i)
	}

	// run nodes
	for i := range node_arr {
		node_arr[i].ch_arr = ch_arr
		go node_arr[i].Run()
	}

	// run while waiting for input
	var input string
	fmt.Scanln(&input)
}
