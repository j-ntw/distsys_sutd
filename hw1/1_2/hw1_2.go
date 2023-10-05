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
	serverClock = numClients + 2
	server_id   = -1
)

func server(ch_arr [numClients]chan Msg, ch_server chan Msg) {
	clock := LamportClock{}
	fmt.Println("start server")
	for {

		// recieve on public channel
		in_msg := <-ch_server
		fmt.Printf("%d->%d @%d: %d\n", in_msg.from, in_msg.to, in_msg.ts, in_msg.data)

		// adjust clock
		clock.AdjustClock(server_id, clock.ts, in_msg.ts)

		// increment own clock
		clock.Inc()

		// flip a coin to send or drop
		if CoinFlip.CoinFlip() {
			// broadcast on private channels
			broadcast_msg := Msg{in_msg.from, in_msg.to, clock.ts, in_msg.data}
			go Broadcast(broadcast_msg, ch_arr)

			// increment own clock
			clock.Inc() //* serverClock

		} else {
			// todo change drop
			fmt.Printf("%d->%d @%d: %d\n", in_msg.from, in_msg.to, in_msg.ts, in_msg.data)
		}

	}
}

func client(ch_client chan Msg, client_id int, ch_server chan Msg, mailbox *Mailbox) {

	clock := LamportClock{}
	fmt.Printf("start c_%d\n", client_id)
	go func() {
		for {
			// increment own clock
			clock.Inc()

			// create message
			out_msg := Msg{client_id, server_id, clock.ts, rand.Intn(10000)}

			// send on public channel
			ch_server <- out_msg
			fmt.Printf("%d->%d @%d: %d\n", out_msg.from, out_msg.to, out_msg.ts, out_msg.data)

			// sleep for nonzero time
			SleepRand.SleepRand()
		}
	}()
	go func() {
		for {
			// receive on private channel
			in_msg := <-ch_client
			fmt.Printf("%d->%d @%d: %d [rB]\n", in_msg.from, in_msg.to, in_msg.ts, in_msg.data)

			// adjust clock
			clock.AdjustClock(client_id, clock.ts, in_msg.ts)

			// increment own clock
			clock.Inc()

			// save message
			mailbox.Append(in_msg)
		}
	}()

}

func main() {
	var ch_arr [numClients]chan Msg
	var ch_server chan Msg = make(chan Msg)
	var mailbox Mailbox

	fmt.Println("create clients")
	for i := range ch_arr {
		// make a channel of type Msg
		// add ch to array
		ch_arr[i] = make(chan Msg)

		go client(ch_arr[i], i, ch_server, &mailbox)

	}
	fmt.Println("create server")
	go server(ch_arr, ch_server)
	var input string
	fmt.Scanln(&input)

	// sort messages in mailbox by timestamp
	mailbox.mu.Lock()
	sort.Slice(mailbox.msg_arr, func(i, j int) bool {
		return mailbox.msg_arr[i].ts < mailbox.msg_arr[j].ts
	})

	// print messages in table
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	mailbox.PrintWhileLocked(w)
	mailbox.mu.Unlock()
	fmt.Println("Done")
}
