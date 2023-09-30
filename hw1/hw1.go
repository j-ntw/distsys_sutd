package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	maxClients = 10
)

func server(c chan []int) {
	for i := 0; ; i++ {

		// get value from channel
		v := <-c
		// flip a coin to send or drop
		if coinFlip() {
			// broadcast with server id, msg
			c <- []int{0, v[1]}
			fmt.Println("Sending ", i)
		}

	}
}
func client(c chan []int, client_id int) {
	for i := 0; ; i++ {
		// create message
		out_msg := []int{client_id, rand.Intn(10000)}
		c <- out_msg

		sleepRand() // do i need to sleep for nonzero time

		// check for message
		in_msg := <-c
		if in_msg[0] != 0{
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
	// create a channel of type integer
	var c chan []int = make(chan []int)

	// launch go routines "server" and "client"

	for i := 1; i == maxClients; i++ {
		go client(c, i)
	}
	go server(c)
	var input string
	fmt.Scanln(&input)
}
