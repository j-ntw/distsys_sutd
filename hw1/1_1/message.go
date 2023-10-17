package main

import (
	"fmt"
)

type Msg struct {
	id   int
	data int
}

func Broadcast(broadcast_msg Msg, ch_arr [NumClients]chan Msg) {
	fmt.Printf("s broadcast: %d: %d\n", broadcast_msg.id, broadcast_msg.data)
	for i, ch_client := range ch_arr {
		if i != broadcast_msg.id {
			ch_client <- broadcast_msg
		}
	}
}
