package main

import "fmt"

type CM struct {
	ch       chan Msg
	id       int
	clock    VectorClock
	page     int
	copy_set Set
	owner    int
	pages    []Page
}

func (cm *CM) listen() {
	fmt.Printf("n%d: start listen\n", cm.id)
	for {
		// receive message
		in_msg := <-cm.ch
		mailbox.Append(in_msg)
		fmt.Printf("n%d: receive %d %d->%d\n", cm.id, in_msg.msgtype, in_msg.from, in_msg.to)

		// increment own vectorclock

		cm.clock.AdjustClock(in_msg.ts)
		cm.clock.Inc(cm.id)

		switch msgtype := in_msg.msgtype; msgtype {
		case ReadRequest:
		case WriteRequest:

		// case ReadForward: CM sends these to P
		// case WriteForward:

		// case SendPage: only P sends to P
		case ReadConfirmation: // CM receives the confirmations
			cm.onReceiveReadConfirmation(in_msg)
		case WriteConfirmation:
		// case Invalidate: // CM sends invalidate command to process

		case InvalidateConfirmation:

		default:
			fmt.Printf("msgtype: %v", msgtype)
		}
	}
}
func (cm *CM) onReceiveReadConfirmation(msg Msg) {
	//	TODO: not sure what to do here
	fmt.Printf("[CM] Received Read Confirmation \n")
}

func (cm *CM) onInvalidateConfirmation(msg Msg) {
	fmt.Printf("[CM] Received Invalidate Confirmation \n")
}

func (cm *CM) Run() {}
func newCM(id int) *CM {
	return &CM{
		ch: make(chan Msg),
		id: id,
	}
}
