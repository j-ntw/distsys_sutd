package main

import (
	"context"
	"fmt"
	"sync"
	"text/tabwriter"
)

type CM struct {
	ch      chan Msg
	id      int
	records Records
	sync.Mutex
}

func (cm *CM) print(w *tabwriter.Writer) {
	fmt.Printf("cm:\n")
	fmt.Fprintln(w, "Record\tOwner_ID\tcopy_set")
	for i, record := range cm.records.Get() {
		fmt.Fprintf(w, "%d\t%d\t%v\n", i, record.owner_id, record.copy_set)
	}
	w.Flush()
}

func (cm *CM) listen(ctx context.Context) {
	fmt.Printf("cm%d: start listen\n", cm.id)
	for {
		select {
		case <-ctx.Done():
			return

		case in_msg := <-cm.ch:
			// receive message
			mailbox.Append(in_msg)

			fmt.Printf("cm%d: receive %s\n", cm.id, in_msg.String())
			switch msgtype := in_msg.msgtype; msgtype {
			case ReadRequest:
				go cm.onReceiveReadRequest(in_msg)
			case WriteRequest:
				go cm.onReceiveWriteRequest(in_msg)
			case ReadConfirmation:
				go cm.onReceiveReadConfirmation(in_msg)
			case WriteConfirmation:
				go cm.onReceiveWriteConfirmation(in_msg)
			case InvalidateConfirmation:
				go cm.onReceiveInvalidateConfirmation(in_msg)
			}
		}
	}
}

// Read
func (cm *CM) onReceiveReadRequest(in_msg Msg) {
	// check page owner, sends read forward to owner
	owner_id := cm.records.GetOwner(in_msg.page_no)
	out_msg := Msg{ReadForward, cm.id, owner_id, in_msg.page_no, in_msg.requester_id}
	go send(p_arr[owner_id].ch, out_msg)

	// add requester to copy set
	cm.records.SetRequester(in_msg.page_no, in_msg.from)
}
func (cm *CM) onReceiveReadConfirmation(in_msg Msg) {
	wg.Done()
}

// Write
func (cm *CM) onReceiveWriteRequest(in_msg Msg) {
	if cm.records.IsCopySetEmpty(in_msg.page_no) {
		// directly invalidateConfirm with self
		out_msg := Msg{InvalidateConfirmation, cm.id, cm.id, in_msg.page_no, in_msg.requester_id}
		go send(cm.ch, out_msg)
	} else {
		// send invalidate to copy set
		for copy_holder_id := range cm.records.GetCopySet(in_msg.page_no) {
			// send invalidate to each copy_holder
			out_msg := Msg{Invalidate, cm.id, copy_holder_id, in_msg.page_no, in_msg.requester_id}
			go send(p_arr[copy_holder_id].ch, out_msg)
		}
	}
}

func (cm *CM) onReceiveInvalidateConfirmation(in_msg Msg) {
	// remove copy_holder from copy_set
	cm.records.DeleteCopyHolder(in_msg.page_no, in_msg.from)
	// send write forward to page owner
	owner_id := cm.records.GetOwner(in_msg.page_no)
	out_msg := Msg{WriteForward, cm.id, owner_id, in_msg.page_no, in_msg.requester_id}
	go send(p_arr[owner_id].ch, out_msg)
}

func (cm *CM) onReceiveWriteConfirmation(in_msg Msg) {
	wg.Done()
}

func newCM(id int) *CM {
	recordTable := make([]CM_Record, 0)
	// intialise each consecutive range of pages to each consecutive process
	for i := 0; i < numPages; i++ {
		record := newRecord(GetInitialOwner(i))
		recordTable = append(recordTable, *record)
	}

	cm := CM{
		ch:      make(chan Msg),
		records: *newRecords(recordTable),
		id:      id + numProcesses,
	}
	cm.print(w)
	return &cm
}
