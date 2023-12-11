package main

import (
	"context"
	"fmt"
	"text/tabwriter"
	"time"
)

const timeout time.Duration = time.Millisecond * 500 //milliseconds

type CM struct {
	ch           chan Msg
	hb_ch        chan bool
	id           int
	cancelActive context.CancelFunc
	ctxActive    context.Context
	role         RoleType
	records      Records
}

func (cm *CM) print(w *tabwriter.Writer) {
	fmt.Printf("cm_%s:\n", cm.role.String())
	fmt.Fprintln(w, "Record\tOwner_ID\tcopy_set")
	for i, record := range cm.records.Get() {
		fmt.Fprintf(w, "%d\t%d\t%v\n", i, record.owner_id, record.copy_set)
	}
	w.Flush()
}

func (cm *CM) listen(ctxMain context.Context) {
	fmt.Printf("cm_%s: start listen\n", cm.role.String())
	defer fmt.Printf("cm_%s: exiting listen\n", cm.role)
	for {
		select {
		case <-cm.ctxActive.Done():
			return
		case <-ctxMain.Done():
			return
		case in_msg := <-cm.ch:
			// receive message

			mailbox.Append(in_msg)

			// fmt.Printf("cm_%s: receive %s\n", cm.role.String(), in_msg.String())
			switch msgtype := in_msg.msgtype; msgtype {
			case ReadRequest:
				cm.onReceiveReadRequest(in_msg)
			case WriteRequest:
				cm.onReceiveWriteRequest(in_msg)
			case ReadConfirmation:
				cm.onReceiveReadConfirmation(in_msg)
			case WriteConfirmation:
				cm.onReceiveWriteConfirmation(in_msg)
			case InvalidateConfirmation:
				cm.onReceiveInvalidateConfirmation(in_msg)
			}
			// if primary, sync state to backup
			if cm.role == Primary {
				copyState(Primary, Backup)
			}
		case <-cm.hb_ch:
			// if we receive a HB message as a Backup, we are no longer active
			go cm.onReceiveHeartBeatCM_listen()
			return

		}

		// // add a delay in processing so we can inject kill commands
		// time.Sleep(time.Second)
	}
}

func (cm *CM) run(ctxMain context.Context) {
	// we link listen and pulse with the same context.
	// if we run the cancel function, we exit listen and pulse

	// TODO: its running too many times!!
	if cm.role == Primary {

		// start listening
		go cm.listen(ctxMain)
		// start sending heartbeat
		go cm.pulse(ctxMain)
	} else if cm.role == Backup {
		// for backup, we can use the same context, since a CM is either primary or backup.
		// again, if necessary, we can cancel the listen ctx inside monitor()
		go cm.monitor(ctxMain)
	}
}

func (cm *CM) pulse(ctxMain context.Context) {
	// only primary needs to pulse
	fmt.Printf("cm_%s: start pulse\n", cm.role.String())
	defer fmt.Printf("cm_%s: exiting pulse\n", cm.role)
	for {
		select {
		case <-cm.ctxActive.Done():
			// The context is canceled, exit the loop
			return
		case <-ctxMain.Done():
			// The context is canceled, exit the loop
			return
		default:
			// send one hearbeat to partner

			// do not send to self
			go func() {
				cm_arr[Backup].hb_ch <- true
			}()

			// pulse every 450ms
			time.Sleep(time.Millisecond * 450)
		}
	}

}

func (cm *CM) monitor(ctxMain context.Context) {
	// only backup needs to monitor

	fmt.Printf("cm_%s: start monitor\n", cm.role.String())
	defer fmt.Printf("cm_%s: stop monitor\n", cm.role.String())
	for {
		select {
		case <-ctxMain.Done():
			// The context is canceled, exit the loop
			return
		// receive HB message
		case <-cm.hb_ch:
			go cm.onReceiveHeartBeatCM_monitor()

		case <-time.After(timeout):
			// if Backup dont get any HB message within timeout, Backup becomes active
			fmt.Printf("cm_%s: Primary failure detected\n", cm.role.String())

			// update the reference that the processes use
			cm_ref.SetRef(&cm_arr[Backup])
			cm.listen(ctxMain)
		}

	}
}

func (cm *CM) onReceiveHeartBeatCM_listen() {
	// As a Backup, if we receive a HB when we are running, we will sync the Primary
	// and go back to monitoring

	// sync
	fmt.Printf("cm_%s: primary is so back\n", cm.role.String())
	copyState(Backup, Primary)

	// and go back to monitoring

	cm_ref.SetRef(&cm_arr[Primary])
	cm.cancelActive()

}

func (cm *CM) onReceiveHeartBeatCM_monitor() {
	// no logic for normal HB
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

func (cm *CM) down() {

	fmt.Printf("cm_%s: set Down\n", cm.role.String())
	cm.cancelActive()
}

func newCM(id int, ch chan Msg, hb_ch chan bool) *CM {
	recordTable := make([]CM_Record, 0)
	// intialise each consecutive range of pages to each consecutive process
	for i := 0; i < numPages; i++ {
		record := newRecord(GetInitialOwner(i))
		recordTable = append(recordTable, *record)
	}
	ctxActive, cancelActive := context.WithCancel(context.Background())

	cm := CM{
		ch:           ch,
		hb_ch:        hb_ch,
		records:      *newRecords(recordTable),
		id:           id + numProcesses,
		role:         Unused,
		cancelActive: cancelActive,
		ctxActive:    ctxActive,
	}
	return &cm
}
func newCMArray() *[]CM {
	cm_arr := make([]CM, numCM)
	ch := make(chan Msg)
	hb_ch := make(chan bool)
	for i := range cm_arr {
		cm_arr[i] = *newCM(i, ch, hb_ch)
	}
	cm_arr[0].role = Primary
	cm_arr[1].role = Backup
	for i := range cm_arr {
		cm_arr[i].print(w)
	}
	return &cm_arr
}
