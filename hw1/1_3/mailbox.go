package main

import (
	"fmt"
	"sync"
	"text/tabwriter"
)

type Mailbox struct {
	msg_arr []Msg
	mu      sync.Mutex
}

func (mailbox *Mailbox) Append(msg Msg) {
	// append a message to the message array in mailbox
	mailbox.mu.Lock()
	mailbox.msg_arr = append(mailbox.msg_arr, msg)
	mailbox.mu.Unlock()
}
func (mailbox *Mailbox) PrintWhileLocked(w *tabwriter.Writer) {
	// assuming mailbox mutex is locked, print each item in the array
	// using the tabwriter formatting
	fmt.Fprintln(w, "From\tTo\tTimestamp\tData")
	for _, msg := range mailbox.msg_arr {
		fmt.Fprintf(w, "%d\t%d\t%v\t%d\n", msg.from, msg.to, msg.ts, msg.data)
	}
	w.Flush()
}
