package main

import "fmt"

type CM struct {
	ch     chan Msg
	id     int
	clock  VectorClock
	record []CM_Record
}

func (cm *CM) listen() {
	fmt.Printf("cm%d: start listen\n", cm.id)
	for {
		// receive message
		in_msg := <-cm.ch
		mailbox.Append(in_msg)
		fmt.Printf("cm%d: receive %v\n", cm.id, in_msg)

		// increment own vectorclock
		// TODO: not sure if i should include CM in process array
		// cm.clock.AdjustClock(in_msg.ts)
		// cm.clock.Inc(cm.id)

		switch msgtype := in_msg.msgtype; msgtype {
		case ReadRequest:
			go cm.onReceieveReadRequest(in_msg)
		case WriteRequest:
		// case ReadForward: CM sends these to P
		// case WriteForward:

		// case SendPage: only P sends to P
		case ReadConfirmation: // CM receives the confirmations

		case WriteConfirmation:
		// case Invalidate: // CM sends invalidate command to process

		case InvalidateConfirmation:

		default:
			fmt.Printf("msgtype: %v", msgtype)
		}
	}
}

func (cm *CM) onReceieveReadRequest(in_msg Msg) {
	// check page owner, sends read forward to owner
	owner_id := cm.record[in_msg.page_no].owner_id
	out_msg := Msg{ReadForward, cm.id, owner_id, in_msg.page_no, in_msg.requester_id, cm.clock.Get()}
	go send(cm.id, p_arr[owner_id].ch, out_msg)

	// add requester to copy set and lock this page
	cm.record[in_msg.page_no].copy_set[in_msg.from] = true
	cm.record[in_msg.page_no].isLocked = true
}

func (cm *CM) Run() {}
func newCM(id int) *CM {
	return &CM{
		ch: make(chan Msg),
	}
}
