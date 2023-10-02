package main

import (
	"fmt"
	"math/rand"
)

func server(ch_arr [NumClients]chan Msg, ch_server chan Msg) {
	fmt.Println("start server")
	for {

		// receive message
		in_msg := <-ch_server

		// flip a coin to send or drop
		fmt.Printf("server recieve from c_%d: %d\n", in_msg.id, in_msg.data)
		if CoinFlip() {
			go Broadcast(in_msg, ch_arr)
		} else {
			fmt.Printf("server drop: c_%d: %d\n", in_msg.id, in_msg.data)
		}

	}
}

func client(ch_client chan Msg, client_id int, ch_server chan Msg) {
	fmt.Printf("start c_%d\n", client_id)
	go func() {
		for {

			// create message
			out_msg := Msg{client_id, rand.Intn(10000)}

			// send on public channel
			ch_server <- out_msg

			fmt.Printf("c_%d send to server: %d\n", out_msg.id, out_msg.data)
			SleepRand() // do i need to sleep for nonzero time

		}
	}()
	go func() {
		for {
			// recieve on private channel
			in_msg := <-ch_client
			fmt.Printf("c_%d recieve from c_%d: %d\n", client_id, in_msg.id, in_msg.data)
		}
	}()
}

func main() {
	var ch_arr [NumClients]chan Msg
	var ch_server chan Msg = make(chan Msg)
	// w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, '.', tabwriter.AlignRight|tabwriter.Debug)
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
