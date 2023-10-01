package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	maxClients = 10
)

type Msg struct {
	id   int
	data int
}

func server(ch [maxClients]chan Msg) {
	for i := 0; ; i++ {

		// get value from channel
		in_msg := <-ch
		// flip a coin to send or drop
		if coinFlip() {
			// broadcast with server id, msg

			fmt.Printf("s_%d broadcast: %d\n", broadcast, broadcast[1])
		}

	}
}
func client(ch chan Msg, client_id int) {
	for i := 0; ; i++ {
		// create message
		out_msg := []int{client_id, rand.Intn(10000)}
		ch <- out_msg
		fmt.Printf("c_%d broadcast: %d\n", out_msg[0], out_msg[1])
		sleepRand() // do i need to sleep for nonzero time

		// check for message
		in_msg := <-ch
		if in_msg[0] != 0 {
			// put it back

			//TODO i think we need 2 channels, otherwise clients might take back their own msg?
		}

	}

}

func coinFlip() bool {
	return rand.Intn(2) == 1
}

func sleepRand() {
	//sleep sporadically
	amt := time.Duration(rand.Intn(1000))
	time.Sleep(time.Millisecond * amt)
}
func main() {
	var ch_arr [maxClients]chan Msg
	for i := 0; i < maxClients; i++ {
		// create a channel of type Msg
		ch := make(chan Msg)
		go client(ch, i) // TODO check race con
		ch_arr[i] = ch
	}
	go server(ch_arr)
	var input string
	fmt.Scanln(&input)
}
