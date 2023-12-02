package main

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"
)

const numProcesses = 3

type AccessType int

const (
	ReadOnly AccessType = iota
	WriteOnly
	ReadWrite
)

type Page struct {
	id   int // owner process id
	data int
}

// CM maintains record for all pages
// the owner and all its copies
// TODO: does each process own multiple pages?
// only one process can own a page at a time (no co owner, ownership can be passed around)
// each page can have at most one owner. there can be page that is orphaned

var (
	mailbox Mailbox
	// ch_arr  [numProcesses]chan Msg
	cm_arr [2]CM
	p_arr  [numProcesses]Process
)

func main() {

	// create CM, processes
	for i := range cm_arr {
		cm_arr[i] = *newCM(i)
	}

	for i := range p_arr {
		p_arr[i] = *newProcess(i)
	}

	// single central manager (with one backup)
	// a few other processes
	// Read Protocol
	// 1. P3 wants to read page x1 (page fault at P3)
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
	go cm_arr[0].listen()
	for i := range p_arr {
		p_arr[i].Run()
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
		return IsBefore(mailbox.msg_arr[i], mailbox.msg_arr[j])

	})

	// print messages in table
	w := tabwriter.NewWriter(os.Stdout, 20, 0, 1, ' ', 0)
	mailbox.PrintWhileLocked(w)
	fmt.Println("Done")
}
