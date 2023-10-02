package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	numClients  = 10
	timeDilator = 1
	serverClock = numClients + 2
)

type Msg struct {
	id    int
	data  int
	clock int
}

type LamportClock struct {
	mu sync.Mutex
	ts int
}

func server(ch_arr [numClients]chan Msg, ch_server chan Msg) {
	clock := 0
	fmt.Println("start server")
	for {

		// recieve on public channel
		in_msg := <-ch_server
		fmt.Printf("server recieve from c_%d: %d, clock_%d\n", in_msg.id, in_msg.data, in_msg.clock)

		// adjust clock
		clock = adjustClock(-1, clock, in_msg.clock)

		// increment own clock
		clock += 1 // * serverClock

		// flip a coin to send or drop
		if CoinFlip() {
			// broadcast on private channels
			broadcast_msg := Msg{in_msg.id, in_msg.data, clock}
			go Broadcast(broadcast_msg, ch_arr)

			// increment own clock
			clock += 1 //* serverClock

		} else {
			fmt.Printf("server drop: c_%d: %d, clock_%d\n", in_msg.id, in_msg.data, in_msg.clock)
		}

	}
}

func client(ch_client chan Msg, client_id int, ch_server chan Msg) {

	clock := LamportClock{}
	fmt.Printf("start c_%d\n", client_id)
	go func() {
		for {
			// increment own clock
			clock.ts += 1 //* (client_id + 1)

			// create message
			out_msg := Msg{client_id, rand.Intn(10000), clock.ts}

			// send on public channel
			ch_server <- out_msg
			fmt.Printf("c_%d send to server: %d, clock_%d\n", out_msg.id, out_msg.data, out_msg.clock)

			// sleep for nonzero time
			SleepRand()
		}
	}()
	go func() {
		for {
			// recieve on private channel
			in_msg := <-ch_client
			fmt.Printf("c_%d recieve from c_%d: %d, clock_%d\n", client_id, in_msg.id, in_msg.data, in_msg.clock)

			// adjust clock
			clock.ts = adjustClock(client_id, clock.ts, in_msg.clock)

			// increment own clock
			clock.ts += 1 //* (client_id + 1)
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
	fmt.Printf("server broadcast msg: c_%d: %d\n", broadcast_msg.id, broadcast_msg.data)
	for i, ch_client := range ch_arr {
		if i != broadcast_msg.id {
			ch_client <- broadcast_msg
		}
	}
}

func adjustClock(id int, clock int, msg_clock int) int {
	if msg_clock > clock {
		fmt.Printf("adjust clock_%d: %d->%d\n", id, clock, msg_clock)
		return msg_clock

	} else {
		fmt.Printf("adjust clock_%d: %d->%d\n", id, clock, clock)
		return clock
	}

}

func main() {
	var ch_arr [numClients]chan Msg
	var ch_server chan Msg = make(chan Msg)
	fmt.Println("create clients")
	for i := range ch_arr {
		// make a channel of type Msg
		// add ch to array
		ch_arr[i] = make(chan Msg)

		go client(ch_arr[i], i, ch_server)

	}
	fmt.Println("create server")
	go server(ch_arr, ch_server)
	var input string
	fmt.Scanln(&input)
}
