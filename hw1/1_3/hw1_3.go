package main

import (
	"CoinFlip"
	"SleepRand"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"text/tabwriter"
)

const (
	numClients  = 10
	server_id   = 0
	numEntities = numClients + 1
)

func IsBefore(tsA [numEntities]int, tsB [numEntities]int) bool {
	// A->B, A happens before B if every A_i <= B_i for all i \elem [0, len(A))
	// A-/>B if any A_i > B_i for all i \elem [0, len(A))
	for k := 0; k < numEntities; k++ {
		if tsA[k] > tsB[k] {
			return false
		}
	}
	return true
}

func server(ch_arr [numEntities]chan Msg) {
	clock := VectorClock{}
	fmt.Println("start server")
	for {

		// recieve on public channel
		in_msg := <-ch_arr[0]

		// check for causality violation
		if clock.isCV(in_msg.ts) {
			fmt.Printf("Causality Violation: %v\n", in_msg)
		}
		// increment own clock
		clock.Inc(server_id)

		// adjust clock
		clock.AdjustClock(clock.ts, in_msg.ts)

		// flip a coin to send or drop
		if CoinFlip.CoinFlip() {
			// broadcast on private channels
			go Broadcast(in_msg, ch_arr)

			// increment own clock
			clock.Inc(server_id)

		}

	}
}

func client(ch_client chan Msg, client_id int, ch_server chan Msg, mailbox *Mailbox) {
	clock := VectorClock{}
	fmt.Printf("start c_%d\n", client_id)
	go func() {
		for {
			// increment own clock
			clock.Inc(client_id)

			// create message
			out_msg := Msg{client_id, server_id, clock.ts, rand.Intn(10000)}

			// send on public channel
			ch_server <- out_msg
			// fmt.Printf("%d->%d @%d: %d\n", out_msg.from, out_msg.to, out_msg.ts, out_msg.data)

			// sleep for nonzero time
			SleepRand.SleepRand()
		}
	}()
	go func() {
		for {
			// recieve on private channel
			in_msg := <-ch_client

			// check CV
			if clock.isCV(in_msg.ts) {
				fmt.Printf("Causality Violation: %v\n", in_msg)
			}
			// increment own clock
			clock.Inc(client_id)

			// adjust clock
			clock.AdjustClock(clock.ts, in_msg.ts)

			// save message
			mailbox.Append(in_msg)
		}
	}()

}

func main() {

	// initialize
	var ch_arr [numEntities]chan Msg
	var mailbox Mailbox

	fmt.Println("create clients")
	for i := range ch_arr {
		// make a channel of type Msg
		// add ch to array
		ch_arr[i] = make(chan Msg)
		if i != 0 {
			// create client ids 1-10
			go client(ch_arr[i], i, ch_arr[0], &mailbox)
		}
	}
	fmt.Println("create server")

	go server(ch_arr)

	// run while waiting for input
	var input string
	fmt.Scanln(&input)

	// this block runs when user enters any input (final button is Enter key)
	// stops goroutines from adding to mailbox and processes its contents

	// sort messages in mailbox by timestamp
	mailbox.mu.Lock()

	sort.SliceStable(mailbox.msg_arr, func(i, j int) bool {
		return IsBefore(mailbox.msg_arr[i].ts, mailbox.msg_arr[j].ts)

	})

	// print messages in table
	w := tabwriter.NewWriter(os.Stdout, 20, 0, 1, ' ', 0)
	mailbox.PrintWhileLocked(w)
	mailbox.mu.Unlock()
	fmt.Println("Done")
}
