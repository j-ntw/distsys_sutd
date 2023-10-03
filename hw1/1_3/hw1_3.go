package main

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"sync"
	"text/tabwriter"
	"time"
)

const (
	numClients  = 10
	server_id   = 0
	numEntities = numClients + 1
)

type Msg struct {
	from int
	to   int
	ts   [numEntities]int
	data int
}
type Mailbox struct {
	msg_arr []Msg
	mu      sync.Mutex
}
type VectorClock struct {
	mu sync.Mutex
	ts [numEntities]int
}

func server(ch_arr [numEntities]chan Msg) {
	clock := VectorClock{}
	fmt.Println("start server")
	for {

		// recieve on public channel
		in_msg := <-ch_arr[0]
		// fmt.Printf("%d->%d @%d: %d\n", in_msg.from, in_msg.to, in_msg.ts, in_msg.data)

		// increment own clock
		clock.Inc(server_id)

		// adjust clock
		clock.AdjustClock(clock.ts, in_msg.ts)

		// flip a coin to send or drop
		if coinFlip() {
			// broadcast on private channels
			go Broadcast(in_msg, ch_arr)

			// increment own clock
			clock.Inc(server_id)

		} else {
			// drop msg
			// fmt.Printf("%d->%d @%d: %d\n", in_msg.from, in_msg.to, in_msg.ts, in_msg.data)
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
			SleepRand()
		}
	}()
	go func() {
		for {
			// recieve on private channel
			in_msg := <-ch_client
			// fmt.Printf("%d->%d @%v: %d [rB]\n", in_msg.from, in_msg.to, in_msg.ts, in_msg.data)

			// increment own clock
			clock.Inc(client_id)

			// adjust clock
			clock.AdjustClock(clock.ts, in_msg.ts)

			// save message
			mailbox.Append(in_msg)
		}
	}()

}

func coinFlip() bool {
	return rand.Intn(2) == 1
}

func SleepRand() {
	// sleep sporadically for [1,1000] ms
	randamt := rand.Intn(1000) + 1
	// fmt.Printf("sleeping: %d ms\n", randamt)
	amt := time.Duration(randamt)
	time.Sleep(time.Millisecond * amt)
}

func Broadcast(broadcast_msg Msg, ch_arr [numEntities]chan Msg) {
	// broadcast from server(ch0) to all channels except originator
	// fmt.Printf("%d->%d @%d: %d [sB]\n", broadcast_msg.from, broadcast_msg.to, broadcast_msg.ts, broadcast_msg.data)
	for i, ch_client := range ch_arr {

		if i != broadcast_msg.from && i != 0 { // dont forward to server ororiginator
			ch_client <- broadcast_msg
		}
	}
}

func (clock *VectorClock) AdjustClock(ts [numEntities]int, msg_ts [numEntities]int) {
	clock.mu.Lock()

	// element wise comparison/swap of ts
	for i := 0; i < numEntities; i++ {
		if msg_ts[i] > ts[i] {
			clock.ts[i] = msg_ts[i]

		} else {
			clock.ts[i] = ts[i]
		}
	}
	// fmt.Printf("\nadjust clock:\n%v\n%v\n\n", ts, clock.ts)
	clock.mu.Unlock()
}

func (clock *VectorClock) Inc(id int) {
	// increment ts for a particular id
	clock.mu.Lock()
	clock.ts[id]++
	clock.mu.Unlock()
}

func (mailbox *Mailbox) Append(msg Msg) {
	// append a message to the message array in mailbox
	mailbox.mu.Lock()
	mailbox.msg_arr = append(mailbox.msg_arr, msg)
	mailbox.mu.Unlock()
}
func (mailbox *Mailbox) PrintWhileLocked(w *tabwriter.Writer) {
	// assuming mailbox mutex is locked, print each item in the array
	// using the tabwriter formatting
	fmt.Fprintln(w, "From\tTo\tTimestamp\tData")
	for _, msg := range mailbox.msg_arr {
		fmt.Fprintf(w, "%d\t%d\t%v\t%d\n", msg.from, msg.to, msg.ts, msg.data)
	}
	w.Flush()
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
		// sort vector clock
		// A->B, A happens before B if every A_i <= B_i for all i \elem [0, len(A))
		// A-/>B if any A_i > B_i for all i \elem [0, len(A))
		for k := 0; k < numEntities; k++ {
			if mailbox.msg_arr[i].ts[k] > mailbox.msg_arr[j].ts[k] {
				return false
			}
		}
		return true
	})

	// print messages in table
	w := tabwriter.NewWriter(os.Stdout, 20, 0, 1, ' ', 0)
	mailbox.PrintWhileLocked(w)
	mailbox.mu.Unlock()
	fmt.Println("Done")
}
