package main

import "fmt"

// message types

type MessageType int

const (
	ReadRequest MessageType = iota
	WriteRequest
	ReadForward
	WriteForward
	SendPage
	ReadConfirmation
	WriteConfirmation
	Invalidate
	InvalidateConfirmation
)

type Msg struct {
	msgtype    MessageType
	from       int
	to         int
	page       Page
	accessType AccessType //AccessType is for sendPage
	ts         [numProcesses]int
}

func send(id int, ch chan Msg, msg Msg) {
	// use as goroutine
	ch <- msg
	fmt.Printf("n%d: send %d %d->%d\n", id, msg.msgtype, msg.from, msg.to)
}
