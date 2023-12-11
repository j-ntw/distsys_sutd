package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"text/tabwriter"
	"time"
)

// global also used in other files
const (
	numProcesses = 10
	numPages     = 100
	numCM        = 1
	numReads     = 100
	numWrites    = 100
)

var (
	mailbox Mailbox
	wg      sync.WaitGroup
	// instantiate/print
	cm         = *newCM(0)
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

	ctx, cancel := context.WithCancel(context.Background())

	// start listeners
	startTime := time.Now()
	go cm.listen(ctx)
	for i := range p_arr {
		go p_arr[i].listen()
	}

	if test_read {
		wg.Add(numReads)
		go func() {
			for i := 0; i < numReads; i++ {
				// random Process wants to read random page (send request)
				randomPage := rand.Intn(numPages)
				randomProcess := rand.Intn(numProcesses)
				p_arr[randomProcess].SendReadRequest(randomPage)

			}
		}()
	}
	if test_write {
		wg.Add(numWrites)
		go func() {
			for i := 0; i < numWrites; i++ {
				// random Process wants to write random page (send request)
				randomPage := rand.Intn(numPages)
				randomProcess := rand.Intn(numProcesses)
				p_arr[randomProcess].SendWriteRequest(randomPage)

			}
		}()
	}

	// wait for all requests to finish
	wg.Wait()
	elapsedTime := time.Since(startTime)

	cancel()
	mailbox.print(w)
	fmt.Printf("Done in %d ms\n", elapsedTime.Milliseconds())
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
