package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	numClients  = 10
	timeDilator = 10
	serverClock = 11
	numEntities = numClients + 1
)

type Msg struct {
	id    int
	data  int
	clock [numEntities]int
}

func server(ch_arr [numClients]chan Msg, ch_server chan Msg) {
	var clock = [numEntities]int{}
	fmt.Println("start server")
	for {
		// increment own clock
		clock[numClients] += 1 * serverClock

		// check which channel needs attention
		in_msg := <-ch_server
		fmt.Printf("server recieve from c_%d: %d, clock_%d\n ", in_msg.id, in_msg.data, in_msg.clock)

		// adjust clock
		clock = adjustClock(numClients, clock, in_msg.clock)

		// flip a coin to send or drop
		if coinFlip() {
			go broadcast(in_msg, ch_arr)
		} else {
			fmt.Printf("server drop: c_%d: %d, clock_%d\n", in_msg.id, in_msg.data, in_msg.clock)
		}

	}
}

func client(ch_client chan Msg, client_id int, ch_server chan Msg) {
	var clock = [numEntities]int{}
	fmt.Printf("start c_%d\n", client_id)
	for {
		// increment own clock
		clock[client_id] += 1 * client_id

		// create message
		out_msg := Msg{client_id, rand.Intn(10000), clock}

		// send on public channel
		ch_server <- out_msg
		fmt.Printf("c_%d send to server: %d, clock_%d\n", out_msg.id, out_msg.data, out_msg.clock)

		// sleep for nonzero time
		sleepRand()
		go func() {
			// recieve on private channel
			in_msg := <-ch_client
			fmt.Printf("c_%d recieve from c_%d: %d, clock_%d\n", client_id, in_msg.id, in_msg.data, in_msg.clock)

			// adjust clock
			clock = adjustClock(client_id, clock, in_msg.clock)
		}()
	}
}

func coinFlip() bool {
	return rand.Intn(2) == 1
}

func sleepRand() {
	//sleep sporadically for [1,1000] * timeDilator ms
	randamt := rand.Intn(1000) + 1
	fmt.Printf("sleeping: %d ms\n", randamt)
	amt := time.Duration(randamt)
	time.Sleep(time.Millisecond * amt * timeDilator)
}

func broadcast(broadcast_msg Msg, ch_arr [numClients]chan Msg) {
	fmt.Printf("server broadcast msg: c_%d: %d\n", broadcast_msg.id, broadcast_msg.data)
	for i, ch_client := range ch_arr {
		if i != broadcast_msg.id {
			ch_client <- broadcast_msg
		}
	}
}

func adjustClock(id int, clock [numEntities]int, msg_clock [numEntities]int) [numEntities]int {
	if msg_clock[id] > clock[id] {
		fmt.Printf("adjust clock: %v->%v", clock, msg_clock)
		return msg_clock

	} else {
		fmt.Printf("adjust clock: %v->%v", clock, clock)
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
		// prevent race condition
		go func(mindex int) {
			go client(ch_arr[mindex], mindex, ch_server)
		}(i)
	}
	fmt.Println("create server")
	go server(ch_arr, ch_server)
	var input string
	fmt.Scanln(&input)
}
