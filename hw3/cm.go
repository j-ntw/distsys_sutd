package main

import (
	"fmt"
	"text/tabwriter"
)

type CM struct {
	ch      chan Msg
	id      int
	records []CM_Record
}

func (cm *CM) print(w *tabwriter.Writer) {
	fmt.Printf("cm%d:\n", cm.id)
	fmt.Fprintln(w, "Record\tOwner_ID\tisLocked\tcopy_set")
	for i, record := range cm.records {
		fmt.Fprintf(w, "%d\t%d\t%v\t%v\n", i, record.owner_id, record.isLocked, record.copy_set)
	}
	w.Flush()
}

func (cm *CM) listen() {
	fmt.Printf("cm%d: start listen\n", cm.id)
	for {
		// receive message
		in_msg := <-cm.ch
		mailbox.Append(in_msg)
		// fmt.Printf("cm%d: receive %v\n", cm.id, in_msg)
		fmt.Printf("cm%d: receive %s\n", cm.id, in_msg.String())
		switch msgtype := in_msg.msgtype; msgtype {
		case ReadRequest:
			go cm.onReceiveReadRequest(in_msg)
		case WriteRequest:
			go cm.onReceiveWriteRequest(in_msg)
		case ReadConfirmation: // CM receives the confirmations
			go cm.onReceiveReadConfirmation(in_msg)
		case WriteConfirmation:
			go cm.onReceiveWriteConfirmation(in_msg)
		case InvalidateConfirmation:
			go cm.onReceiveInvalidateConfirmation(in_msg)
		default:
			fmt.Printf("msgtype: %v", msgtype)
		}
	}
}

// Read
func (cm *CM) onReceiveReadRequest(in_msg Msg) {
	// check page owner, sends read forward to owner
	owner_id := cm.records[in_msg.page_no].owner_id
	out_msg := Msg{ReadForward, cm.id, owner_id, in_msg.page_no, in_msg.requester_id}
	send(cm.id, p_arr[owner_id].ch, out_msg)

	// add requester to copy set and lock this page
	cm.records[in_msg.page_no].copy_set[in_msg.from] = true
	cm.records[in_msg.page_no].isLocked = true
}
func (cm *CM) onReceiveReadConfirmation(in_msg Msg) {
	cm.records[in_msg.page_no].isLocked = false
}

// Write
func (cm *CM) onReceiveWriteRequest(in_msg Msg) {

	// send invalidate to copy set
	for copy_holder_id := range cm.records[in_msg.page_no].copy_set {
		// send invalidate to each copy_holder
		out_msg := Msg{Invalidate, cm.id, copy_holder_id, in_msg.page_no, in_msg.requester_id}
		send(cm.id, p_arr[copy_holder_id].ch, out_msg)
	}
}

func (cm *CM) onReceiveInvalidateConfirmation(in_msg Msg) {
	// remove copy_holder from copy_set
	delete(cm.records[in_msg.page_no].copy_set, in_msg.from)
	// send write forward to page owner
	owner_id := cm.records[in_msg.page_no].owner_id
	out_msg := Msg{WriteForward, cm.id, owner_id, in_msg.page_no, in_msg.requester_id}
	send(cm.id, p_arr[owner_id].ch, out_msg)
}

func (cm *CM) onReceiveWriteConfirmation(in_msg Msg) {
	cm.records[in_msg.page_no].isLocked = false
}

func newCM(id int) *CM {
	recordTable := make([]CM_Record, 0)
	// intialise each consecutive range of pages to each consecutive process
	for i := 0; i < numPages; i++ {
		record := newRecord(GetInitialOwner(i))
		recordTable = append(recordTable, *record)
	}
	return &CM{
		ch:      make(chan Msg),
		records: recordTable,
	}
}
