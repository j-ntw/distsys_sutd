package main

import (
	"fmt"
	"os"
	"text/tabwriter"
)

const (
	numProcesses = 3
	numPages     = 5
	numCM        = 1
)

// CM maintains record for all pages
// the owner and all its copies
// TODO: does each process own multiple pages?
// only one process can own a page at a time (no co owner, ownership can be passed around)
// each page can have at most one owner. there can be page that is orphaned

var (
	mailbox Mailbox
	cm      = *newCM(numProcesses)
	p_arr   [numProcesses]Process
	w       = tabwriter.NewWriter(os.Stdout, 10, 0, 1, ' ', 0)
)

func main() {
	// instantiate/print
	cm.print(w)
	p_arr := make([]Process, numProcesses)
	for i := range p_arr {
		p_arr[i] = *newProcess(i)
		p_arr[i].print(w)
	}

	// start listeners
	go cm.listen()
	for i := range p_arr {
		go p_arr[i].listen()
	}
	// Read Protocol
	// 1. P3 wants to read page 1 (page fault at P3)
	// 2. P3 sends read req to CM (X1. P1)
	// 3. CM checks page owner and sends read forward to owner P1, adds P3 in copy set
	// 4. P3 sends Read confirmation to central manager

	// Write Protocol
	// P2 encounters page fault while writing page X1
	// 1. P2 sends a write request to the CM
	// 2. P3 sends confirmation of invalidation msg to CM
	// 3. CM receives invalidate cfm and clears copy set
	// 4. CM sends write forward for page X1 to P1.
	// 5. P1 sends X1 to P2, invalidates own copy of X1.
	// 6. P2 sends write confirmation for X1 to CM.

	// P3 wants to read page x1 (send request)
	p_arr[2].SendReadRequest(1)

	// run while waiting for input
	var input string
	fmt.Scanln(&input)
	// this block runs when user enters any input (final button is Enter key)
	// stops goroutines from adding to mailbox and processes its contents
	// sort messages in mailbox by timestamp
	mailbox.Lock()
	defer mailbox.Unlock()

	// print messages in table

	mailbox.PrintWhileLocked(w)
	fmt.Println("Done")
}
