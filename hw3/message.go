package main

import "fmt"

// message types

type MessageType int

const (
	ReadRequest MessageType = iota
	WriteRequest
	ReadForward
	WriteForward
	ReadPage
	WritePage
	ReadConfirmation
	WriteConfirmation
	Invalidate
	InvalidateConfirmation
)

type Msg struct {
	msgtype      MessageType
	from         int
	to           int
	page_no      int
	requester_id int
	ts           [numProcesses]int
}

func send(id int, ch chan Msg, msg Msg) {
	// use as goroutine
	ch <- msg
	fmt.Printf("n%d: send %d %d->%d\n", id, msg.msgtype, msg.from, msg.to)
}
