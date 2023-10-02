package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	NumClients  = 10
	TimeDilator = 10
)

type Msg struct {
	id   int
	data int
}

func CoinFlip() bool {
	return rand.Intn(2) == 1
}

func SleepRand() {
	//sleep sporadically for [1,1000] * TimeDilator ms
	randamt := rand.Intn(1000) + 1
	fmt.Printf("sleeping: %d ms\n", randamt)
	amt := time.Duration(randamt)
	time.Sleep(time.Millisecond * amt * TimeDilator)
}

func Broadcast(broadcast_msg Msg, ch_arr [NumClients]chan Msg) {
	fmt.Printf("server broadcast msg: c_%d: %d\n", broadcast_msg.id, broadcast_msg.data)
	for i, ch_client := range ch_arr {
		if i != broadcast_msg.id {
			ch_client <- broadcast_msg
		}
	}
}
