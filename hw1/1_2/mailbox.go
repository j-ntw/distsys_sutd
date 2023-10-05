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
	mailbox.mu.Lock()
	mailbox.msg_arr = append(mailbox.msg_arr, msg)
	mailbox.mu.Unlock()
}
func (mailbox *Mailbox) PrintWhileLocked(w *tabwriter.Writer) {
	fmt.Fprintln(w, "From\tTo\tTimestamp\tData")
	for _, msg := range mailbox.msg_arr {
		fmt.Fprintf(w, "%d\t%d\t%d\t%d\n", msg.from, msg.to, msg.ts, msg.data)
	}
	w.Flush()
}
