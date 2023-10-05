package main

import "fmt"

type Msg struct {
	from int
	to   int
	ts   int
	data int
}

func Broadcast(broadcast_msg Msg, ch_arr [numClients]chan Msg) {
	fmt.Printf("%d->%d @%d: %d\n", broadcast_msg.from, broadcast_msg.to, broadcast_msg.ts, broadcast_msg.data)
	for i, ch_client := range ch_arr {
		if i != broadcast_msg.from {
			ch_client <- broadcast_msg
		}
	}
}
