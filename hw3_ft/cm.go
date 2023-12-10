package main

import (
	"fmt"
	"sync"
	"text/tabwriter"
	"time"
)

const timeout time.Duration = time.Millisecond * 500 //milliseconds

type CM struct {
	ch           chan Msg
	id           int
	mode_changed chan bool
	mode         Mode
	role         RoleType
	records      []CM_Record
	sync.Mutex
}

func (cm *CM) print(w *tabwriter.Writer) {
	fmt.Printf("cm%d: %s\n", cm.id, cm.role.String())
	fmt.Fprintln(w, "Record\tOwner_ID\tisLocked\tcopy_set")
	for i, record := range cm.records {
		fmt.Fprintf(w, "%d\t%d\t%v\t%v\n", i, record.owner_id, record.isLocked, record.copy_set)
	}
	w.Flush()
}

func (cm *CM) listen() {
	fmt.Printf("cm%d: start listen\n", cm.id)
	for {
		if !cm.GetIsRunning() {
			break
		}
		// receive message
		in_msg := <-cm.ch
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
		case HeartBeatCM:
			go cm.onReceiveHeartBeatCM(in_msg)
		case Down:
			return
		default:
			fmt.Printf("msgtype: %v", msgtype)
			panic(in_msg)

		}

		// if primary, sync state to backup
		if cm.role == Primary {
			copyState(Primary, Backup)
		}
		// add a delay in processing so we can inject kill commands
		time.Sleep(time.Second)
	}
}

func (cm *CM) ChangeMode(isRunning bool) {
	cm.SetIsRunning(isRunning)
	cm.mode_changed <- isRunning
	if isRunning {
		fmt.Printf("cm%d: mode running\n", cm.id)
	} else {
		fmt.Printf("cm%d: mode stopped\n", cm.id)
	}

}

func (cm *CM) GetIsRunning() bool { // helper method
	return cm.mode.GetIsRunning()
}

func (cm *CM) SetIsRunning(newBool bool) { // helper method
	cm.mode.SetIsRunning(newBool)
}

func (cm *CM) run(change bool) {
	go cm.ChangeMode(change) // use goroutine because receiving channel is only called later
	for {
		<-cm.mode_changed

		if cm.role == Primary {
			if cm.GetIsRunning() {
				// start listening
				go cm.listen()
				// start sending heartbeat
				go cm.throb()
			} else {
				return
			}
		}
		if cm.role == Backup {
			if cm.GetIsRunning() {
				// start listening
				go cm.listen()
			} else {
				// backups typically start off monitoring
				go cm.monitor()
			}
		}
	}
}

func scaleDuration(duration time.Duration, factor float64) time.Duration {
	// Convert duration to float64, multiply, and convert back to time.Duration
	scaledDuration := time.Duration(float64(duration) * factor)
	return scaledDuration
}

func (cm *CM) throb() {
	// only primary needs to throb
	fmt.Printf("cm%d: start throb\n", cm.id)
	for {
		if !cm.GetIsRunning() {
			return
		}
		// send one hearbeat to partner
		out_msg := Msg{msgtype: HeartBeatCM}

		// do not send to self
		send(cm.id, cm_arr[1].ch, out_msg)

		// throb every 450ms
		time.Sleep(scaleDuration(timeout, 0.9))
	}

}

func (cm *CM) monitor() {
	// only backup needs to throb

	fmt.Printf("cm%d: start monitor\n", cm.id)
	for {
		select {
		// receive message
		case in_msg := <-cm.ch:
			mailbox.Append(in_msg) // TODO: may not want to record HB messages
			fmt.Printf("cm%d: receive %s\n", cm.id, in_msg.String())
			switch msgtype := in_msg.msgtype; msgtype {
			case HeartBeatCM:
				go cm.onReceiveHeartBeatCM(in_msg)
			case Down:
				return
			default:
				fmt.Printf("msgtype: %v", msgtype)
				panic(in_msg)
			}
		case <-time.After(timeout):
			// if you dont get any message within timeout, Backup becomes active
			cm.ChangeMode(true)
			// update the reference that the processes use
			cm_ref.SetRef(&cm_arr[1])
			return
		}

	}
}
func (cm *CM) onReceiveHeartBeatCM(in_msg Msg) {
	// no logic for normal HB, but as a Backup, if we receive a HB when we are running, we will sync the Primary
	// and go back to monitoring

	if cm.GetIsRunning() { // check if Backup is running
		// sync
		copyState(Backup, Primary)

		// and go back to monitoring
		cm.ChangeMode(false)
		cm_ref.SetRef(&cm_arr[0])
	}
}

// Read
func (cm *CM) onReceiveReadRequest(in_msg Msg) {
	cm.Lock()
	defer cm.Unlock()
	// check page owner, sends read forward to owner
	owner_id := cm.records[in_msg.page_no].owner_id
	out_msg := Msg{ReadForward, cm.id, owner_id, in_msg.page_no, in_msg.requester_id}
	send(cm.id, p_arr[owner_id].ch, out_msg)

	// add requester to copy set and lock this page
	cm.records[in_msg.page_no].copy_set[in_msg.from] = true
	cm.records[in_msg.page_no].isLocked = true
}

func (cm *CM) onReceiveReadConfirmation(in_msg Msg) {
	cm.Lock()
	defer cm.Unlock()
	cm.records[in_msg.page_no].isLocked = false
}

// Write
func (cm *CM) onReceiveWriteRequest(in_msg Msg) {
	cm.Lock()
	defer cm.Unlock()
	if len(cm.records[in_msg.page_no].copy_set) == 0 {
		// directly invalidateConfirm with self
		out_msg := Msg{InvalidateConfirmation, cm.id, cm.id, in_msg.page_no, in_msg.requester_id}
		send(cm.id, cm.ch, out_msg)
	} else {
		// send invalidate to copy set
		for copy_holder_id := range cm.records[in_msg.page_no].copy_set {
			// send invalidate to each copy_holder
			out_msg := Msg{Invalidate, cm.id, copy_holder_id, in_msg.page_no, in_msg.requester_id}
			send(cm.id, p_arr[copy_holder_id].ch, out_msg)
		}
	}
}

func (cm *CM) onReceiveInvalidateConfirmation(in_msg Msg) {
	cm.Lock()
	defer cm.Unlock()
	// remove copy_holder from copy_set
	delete(cm.records[in_msg.page_no].copy_set, in_msg.from)
	// send write forward to page owner
	owner_id := cm.records[in_msg.page_no].owner_id
	out_msg := Msg{WriteForward, cm.id, owner_id, in_msg.page_no, in_msg.requester_id}
	send(cm.id, p_arr[owner_id].ch, out_msg)
}

func (cm *CM) onReceiveWriteConfirmation(in_msg Msg) {
	cm.Lock()
	defer cm.Unlock()
	cm.records[in_msg.page_no].isLocked = false
}

func (cm *CM) down() {
	if cm.GetIsRunning() {
		cm.SetIsRunning(false)
		out_msg := Msg{msgtype: Down}
		go send(cm.id, cm.ch, out_msg)
	}
	// idempotent
}

func newCM(id int) *CM {
	recordTable := make([]CM_Record, 0)
	// intialise each consecutive range of pages to each consecutive process
	for i := 0; i < numPages; i++ {
		record := newRecord(GetInitialOwner(i))
		recordTable = append(recordTable, *record)
	}

	cm := CM{
		ch:           make(chan Msg),
		records:      recordTable,
		id:           id + numProcesses,
		mode:         *newMode(),
		role:         Unused,
		mode_changed: make(chan bool),
	}
	return &cm
}
func newCMArray() *[]CM {
	cm_arr := make([]CM, numCM)
	for i := range cm_arr {
		cm_arr[i] = *newCM(i)
	}
	cm_arr[0].role = Primary
	cm_arr[1].role = Backup
	for i := range cm_arr {
		cm_arr[i].print(w)
	}
	return &cm_arr
}
