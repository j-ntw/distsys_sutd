package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
)

// global also used in other files
const (
	numProcesses = 3
	numPages     = 5
	numCM        = 2
)

var (
	mailbox Mailbox
	// instantiate/print
	cm_ref     *CM_REF // a reference to either primary or backup CM that processes use with a mutex
	cm_arr     = *newCMArray()
	p_arr      = *newProcessArray()
	w          = tabwriter.NewWriter(os.Stdout, 10, 0, 1, ' ', 0)
	test_read  bool
	test_write bool
)

func main() {
	// get test case
	flag.BoolVar(&test_read, "r", false, "Testing read once.")
	flag.BoolVar(&test_write, "w", false, "Testing write once.")
	flag.Parse()

	// set main up
	cm_ref = newCM_REF(&cm_arr[0])
	ctx, cancel := context.WithCancel(context.Background())
	// start listeners
	go cm_arr[int(Primary)].run(ctx) // Primary starts in active running/listening
	go cm_arr[int(Backup)].run(ctx)  // Backup starts in passive monitoring
	for i := range p_arr {
		go p_arr[i].listen()
	}
	go func() {
		if test_read {
			// P3 wants to read page x1 (send request)
			p_arr[2].SendReadRequest(1)
			//down Primary, Backup should detect loss of HB and start itself
			cm_arr[Primary].down()

			// // P3 wants to read page x0 (send request)
			p_arr[2].SendReadRequest(0)
			// // // start Primary again
			cm_arr[Primary].run(ctx)
		} else if test_write {
			// optional: make some copies
			// p_arr[0].SendReadRequest(1)
			cm_arr[Primary].down()
			cm_arr[Primary].run(ctx)
			// P3 wants to read page x1 (send request)
			p_arr[2].SendWriteRequest(1)
		}

		// this block runs when user enters any input (final button is Enter key)
		// stops goroutines from adding to mailbox and processes its contents

	}()
	// run while waiting for input

	var input string
	fmt.Scanln(&input)
	cancel()
	mailbox.print(w)
	fmt.Println("Done")
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

// CM maintains record for all pages
// the owner and all its copies

// only one process can own a page at a time (no co owner, ownership can be passed around)
// each page can have only one owner.
