package main

import (
	"fmt"
	"sync"
	"text/tabwriter"
	"time"
)

type RoleType int

const (
	Primary RoleType = iota
	Backup
)
const timeout = 500 //milliseconds
type CM struct {
	ch           chan Msg
	id           int
	mode_changed chan bool
	isRunning    bool
	role         RoleType
	records      []CM_Record
	sync.Mutex
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
		if !cm.isRunning {
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
		case ReadConfirmation: // CM receives the confirmations
			go cm.onReceiveReadConfirmation(in_msg)
		case WriteConfirmation:
			go cm.onReceiveWriteConfirmation(in_msg)
		case InvalidateConfirmation:
			go cm.onReceiveInvalidateConfirmation(in_msg)
		case HeartBeatCM:
			go cm.onReceiveHeartBeatCM(in_msg)
		case Down:
			break
		default:
			fmt.Printf("msgtype: %v", msgtype)
			panic(in_msg)

		}

		// if primary, sync state to backup
		if cm.role == Primary {
			cm_arr[0].Lock()
			cm_arr[1].Lock()
			defer cm_arr[1].Unlock()
			defer cm_arr[0].Unlock()
			cm_arr[1].records = cm_arr[0].records
		}
		// add a delay in processing so we can inject kill commands
		time.Sleep(time.Second)
	}
}

func (cm *CM) ChangeMode(change bool) {
	cm.isRunning = change
	cm.mode_changed <- change
}

func (cm *CM) run(change bool) {
	go cm.ChangeMode(change)
	for {
		<-cm.mode_changed
		if cm.isRunning {
			// start listening
			go cm.listen()
			// start sending heartbeat
			go cm.throb()
		} else {
			go cm.monitor()
		}
	}

}

func (cm *CM) throb() {
	// only primary needs to throb
	if cm.role == Primary {
		for {
			if !cm.isRunning {
				break
			}
			// send one hearbeat to each partner
			out_msg := Msg{msgtype: HeartBeatCM}
			for partner := range cm_arr {
				send(cm.id, cm_arr[partner].ch, out_msg)
			}

			// sleep
			time.Sleep(time.Millisecond * timeout)
		}
	}

}

func (cm *CM) monitor() {
	// TODO: listen for heartbeat passively. if not actively listening then monitoring

	// if Primary, you do not receive heartbeats normally
	// if Backup, you always receive heartbeats.
	// when heartbeat stops, start running
	fmt.Printf("cm%d: start monitor\n", cm.id)
	noResponse := false
	go func() {
		// try to run
		for {
			if noResponse {
				break
			} else {
				//sleep
				time.Sleep(time.Millisecond * timeout)
			}
		}

	}()
	for {
		// receive message
		in_msg := <-cm.ch
		mailbox.Append(in_msg)
		fmt.Printf("cm%d: receive %s\n", cm.id, in_msg.String())

		switch msgtype := in_msg.msgtype; msgtype {
		case HeartBeatCM:
			go cm.onReceiveHeartBeatCM(in_msg)
		case Down:
			break
		default:
			fmt.Printf("msgtype: %v", msgtype)
			panic(in_msg)
		}
	}
	// leave monitor to run.
}
func (cm *CM) onReceiveHeartBeatCM(in_msg Msg) {
	// TODO: reset counter
	cm.noResponse = false
	// if Primary, you do not receive heartbeats normally
	// if Backup, you always receive heartbeats.
	// when heartbeat stops,
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
			fmt.Printf("here!!!\n")
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
	if cm.isRunning {
		cm.isRunning = false
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
		isRunning:    false,
		mode_changed: make(chan bool),
	}
	cm.print(w)
	return &cm
}
func newCMArray() *[]CM {
	cm_arr := make([]CM, numCM)
	for i := range cm_arr {
		cm_arr[i] = *newCM(i)
		cm_arr[i].print(w)
	}
	return &cm_arr
}
