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

func (m MessageType) String() string {
	return [...]string{"ReadRequest", "WriteRequest", "ReadForward", "WriteForward", "ReadPage", "WritePage",
		"ReadConfirmation", "WriteConfirmation", "Invalidate", "InvalidateConfirmation"}[m]
}

func (msg Msg) String() string {
	return fmt.Sprintf("Msg{msgtype: %s, from: %d, to: %d, page_no: %d, requester_id: %d}",
		msg.msgtype.String(), msg.from, msg.to, msg.page_no, msg.requester_id)
}

type Msg struct {
	msgtype      MessageType
	from         int
	to           int
	page_no      int
	requester_id int
}

func send(id int, ch chan Msg, msg Msg) {
	// use as goroutine
	ch <- msg
	// fmt.Printf("n%d: send %s %d->%d\n", id, msg.msgtype.String(), msg.from, msg.to)
	fmt.Printf("n%d: send %s\n", id, msg.String())
}
