package main

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"
)

// implement lamports shared priority queue with vector clock

const (
	numNodes = 3
)

var ch_arr [numNodes]chan Msg
var node_arr [numNodes]Node
var mailbox Mailbox

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
	// this block runs when user enters any input (final button is Enter key)
	// stops goroutines from adding to mailbox and processes its contents
	// sort messages in mailbox by timestamp
	mailbox.Lock()
	defer mailbox.Unlock()
	sort.SliceStable(mailbox.msg_arr, func(i, j int) bool {
		return IsBefore(mailbox.msg_arr[i].ts, mailbox.msg_arr[j].ts)

	})

	// print messages in table
	w := tabwriter.NewWriter(os.Stdout, 20, 0, 1, ' ', 0)
	mailbox.PrintWhileLocked(w)
	fmt.Println("Done")
}
