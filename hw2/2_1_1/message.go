package main

import "fmt"

type Msg struct {
	msgtype Msgtype
	from    int
	to      int
	ts      [numNodes]int
}

// message types
type Msgtype int

const ( // iota is reset to 0
	req  Msgtype = iota //  0
	resp                //  1
	release
)

func send(ch chan Msg, msg Msg) {
	fmt.Printf("send %d->%d\n", msg.from, msg.to)
	ch <- msg
}
