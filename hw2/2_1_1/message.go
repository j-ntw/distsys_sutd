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
	req     Msgtype = iota //  0
	reply                  //  1
	release                //  2
)

func send(id int, ch chan Msg, msg Msg) {
	// use as goroutine
	fmt.Printf("n%d: send %d->%d\n", id, msg.from, msg.to)
	ch <- msg
}
