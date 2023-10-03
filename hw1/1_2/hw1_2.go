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
	timeDilator = 1
	serverClock = numClients + 2
	server_id   = -1
)

type Msg struct {
	from int
	to   int
	ts   int
	data int
}

type Mailbox struct {
	msg_arr []Msg
	mu      sync.Mutex
}

type LamportClock struct {
	mu sync.Mutex
	ts int
}

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
		clock.Inc() // * serverClock

		// flip a coin to send or drop
		if CoinFlip() {
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
			SleepRand()
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

func CoinFlip() bool {
	return rand.Intn(2) == 1
}

func SleepRand() {
	//sleep sporadically for [1,1000] * timeDilator ms
	randamt := rand.Intn(1000) + 1
	fmt.Printf("sleeping: %d ms\n", randamt)
	amt := time.Duration(randamt)
	time.Sleep(time.Millisecond * amt * timeDilator)
}

func Broadcast(broadcast_msg Msg, ch_arr [numClients]chan Msg) {
	fmt.Printf("%d->%d @%d: %d\n", broadcast_msg.from, broadcast_msg.to, broadcast_msg.ts, broadcast_msg.data)
	for i, ch_client := range ch_arr {
		if i != broadcast_msg.from {
			ch_client <- broadcast_msg
		}
	}
}

func (clock *LamportClock) AdjustClock(id int, ts int, msg_ts int) {
	clock.mu.Lock()
	if msg_ts > ts {
		fmt.Printf("adjust clock_%d: %d->%d\n", id, ts, msg_ts)
		clock.ts = msg_ts

	} else {
		fmt.Printf("adjust clock_%d: %d->%d\n", id, ts, ts)
		clock.ts = ts
	}
	clock.mu.Unlock()

}

func (clock *LamportClock) Inc() {
	clock.mu.Lock()
	clock.ts += 1
	clock.mu.Unlock()
}

func (mailbox *Mailbox) Append(msg Msg) {
	mailbox.mu.Lock()
	mailbox.msg_arr = append(mailbox.msg_arr, msg)
	mailbox.mu.Unlock()
}
func (mailbox *Mailbox) PrintWhileLocked(w *tabwriter.Writer) {
	fmt.Fprintln(w, "From\tTo\tTimestamp\tData")
	for _, msg := range mailbox.msg_arr {
		fmt.Fprintf(w, "%d\t%d\t%d\t%d\n", msg.from, msg.to, msg.ts, msg.data)
	}
	w.Flush()
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
