package main

import "fmt"

const (
	AddressSpaceSize = numPages / numProcesses
)

type Process struct {
	ch     chan Msg
	id     int
	ptable []Page // page table
	clock  VectorClock
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
		case ReadForward:
			go p.onReceiveReadForward(in_msg)
		case WriteForward:
			go p.onReceiveWriteForward(in_msg)
		case ReadPage:
			go p.onReceiveReadPage(in_msg)
		case WritePage:
			go p.onReceiveWritePage(in_msg)
		case Invalidate:
			go p.onReceiveInvalidate(in_msg)
		default:
			fmt.Printf("msgtype: %v", msgtype)
		}
	}
}

// Read
func (p *Process) SendReadRequest(page_no int) {
	out_msg := Msg{
		ReadRequest,
		p.id,
		-1, // CM
		page_no,
		p.id,
		p.clock.Get(),
	}
	send(p.id, cm.ch, out_msg)

	// lock page
	p.ptable[page_no].isLocked = true
}

func (p *Process) onReceiveReadForward(in_msg Msg) {
	// lock page, change access
	p.ptable[in_msg.page_no].isLocked = true
	p.ptable[in_msg.page_no].access = ReadOnly

	// send page to requester
	out_msg := Msg{ReadPage, p.id, in_msg.requester_id, in_msg.page_no, in_msg.requester_id, p.clock.Get()}
	send(p.id, p_arr[in_msg.requester_id].ch, out_msg)
	// we simulate the sending of pages with the sendpage typed message,
	// ideally the actual page will be included
}

func (p *Process) onReceiveReadPage(in_msg Msg) {
	// send read confirmation to CM
	out_msg := Msg{ReadConfirmation, p.id, cm.id, in_msg.page_no, in_msg.requester_id, p.clock.Get()}
	send(p.id, cm.ch, out_msg)
}

// Write
func (p *Process) SendWriteRequest(page_no int) {
	out_msg := Msg{
		WriteRequest,
		p.id,
		-1, // CM
		page_no,
		p.id,
		p.clock.Get(),
	}
	send(p.id, cm.ch, out_msg)
}

func (p *Process) onReceiveInvalidate(in_msg Msg) {
	// send back to CM InvalidateConfirmation
	out_msg := Msg{InvalidateConfirmation, p.id, cm.id, in_msg.page_no, in_msg.requester_id, cm.clock.Get()}
	send(p.id, cm.ch, out_msg)
}

func (p *Process) onReceiveWriteForward(in_msg Msg) {
	// invalidate own copy by setting access to nil, isOwner to false
	// sending data is simulated with the send page message type
	p.ptable[in_msg.page_no].isOwner = false
	p.ptable[in_msg.page_no].access = Nil

	// sendPage to requester
	out_msg := Msg{WritePage, p.id, in_msg.requester_id, in_msg.page_no, in_msg.requester_id, cm.clock.Get()}
	send(p.id, p_arr[in_msg.requester_id].ch, out_msg)
}

func (p *Process) onReceiveWritePage(in_msg Msg) {
	// send write confirmation to CM
	out_msg := Msg{WriteConfirmation, p.id, cm.id, in_msg.page_no, in_msg.requester_id, p.clock.Get()}
	send(p.id, cm.ch, out_msg)
}

func newProcess(id int) *Process {
	// assign page range
	ptable := make([]Page, numPages)
	for i := 0; i < numPages; i++ {

		isOwner := GetInitialOwner(i) == id
		page := newPage(isOwner)
		ptable = append(ptable, *page)
	}

	return &Process{
		ch:     make(chan Msg),
		id:     id,
		ptable: ptable,
	}
}

func GetInitialOwner(id int) int {
	return id / AddressSpaceSize
}
