package main

type Msg struct {
	from int
	to   int
	ts   [numEntities]int
	data int
}

func Broadcast(broadcast_msg Msg, ch_arr [numEntities]chan Msg) {
	// broadcast from server(ch0) to all channels except originator
	// fmt.Printf("%d->%d @%d: %d [sB]\n", broadcast_msg.from, broadcast_msg.to, broadcast_msg.ts, broadcast_msg.data)
	for i, ch_client := range ch_arr {

		if i != broadcast_msg.from && i != 0 { // dont forward to server or originator
			ch_client <- broadcast_msg
		}
	}
}
