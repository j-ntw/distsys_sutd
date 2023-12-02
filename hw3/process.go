package main

import "fmt"

type Process struct {
	ch     chan Msg
	id     int
	page   int
	access bool
	clock  VectorClock
}

func (process *Process) read() {

}

func (p *Process) listen() {
	fmt.Printf("n%d: start listen\n", p.id)
	for {
		// receive message
		in_msg := <-p.ch
		mailbox.Append(in_msg)
		fmt.Printf("n%d: receive %d %d->%d\n", p.id, in_msg.msgtype, in_msg.from, in_msg.to)

		// increment own vectorclock

		p.clock.AdjustClock(in_msg.ts)
		p.clock.Inc(p.id)

		switch msgtype := in_msg.msgtype; msgtype {
		// case ReadRequest: requests are received by CM
		// case WriteRequest:

		case ReadForward:
		case WriteForward:

		case SendPage:
		// case ReadConfirmation: CM receives the confirmations
		// case WriteConfirmation:
		case Invalidate: // CM sends invalidate command to process

		// case InvalidateConfirmation:

		default:
			fmt.Printf("msgtype: %v", msgtype)
		}
	}
}
func (p *Process) Run() {}
func newProcess(id int) *Process {
	return &Process{
		ch: make(chan Msg),
		id: id,
	}
}
