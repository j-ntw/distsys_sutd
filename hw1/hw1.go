package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	numClients = 10
)

type Msg struct {
	id   int
	data int
}

func server(ch_arr [numClients]chan Msg, ch_server chan int) {
	for {
		// check which channel needs attention
		ch_id := <-ch_server
		in_msg := <-ch_arr[ch_id]
		// flip a coin to send or drop
		if coinFlip() {
			// broadcast with server id, msg
			broadcast(in_msg, ch_arr, ch_id)
		} else {
			fmt.Println("dropped msg")
		}

	}
}

func client(ch_client chan Msg, client_id int, ch_server chan int) {
	for {
		// create message
		out_msg := Msg{client_id, rand.Intn(10000)}
		ch_client <- out_msg
		// notify server
		ch_server <- client_id
		fmt.Printf("c_%d send: %d\n", out_msg.id, out_msg.data)
		sleepRand() // do i need to sleep for nonzero time
		in_msg := <-ch_client
		fmt.Printf("c_%d recieve from c_%d: %d\n", client_id, in_msg.id, in_msg.data)
	}
}

func coinFlip() bool {
	return rand.Intn(2) == 1
}

func sleepRand() {
	//sleep sporadically
	randamt := rand.Intn(1000)
	fmt.Printf("sleeping: %d ms\n", randamt)
	amt := time.Duration(randamt)
	time.Sleep(time.Millisecond * amt)
}

func broadcast(broadcast_msg Msg, ch_arr [numClients]chan Msg, id int) {
	fmt.Println("broadcast msg")
	for i, ch_client := range ch_arr {
		if i != id {
			ch_client <- broadcast_msg
		}
	}
}
func main() {
	var ch_arr [numClients]chan Msg
	var ch_server chan int
	fmt.Println("create clients")
	for i := range ch_arr {
		// make a channel of type Msg
		// add ch to array
		ch_arr[i] = make(chan Msg)
		// prevent race condition
		go func(mindex int) {
			client(ch_arr[mindex], mindex, ch_server)
		}(i)
	}
	fmt.Println("create server")
	go server(ch_arr, ch_server)
	var input string
	fmt.Scanln(&input)
}
