package main

import (
	"fmt"
	"sync"
	"text/tabwriter"
)

type Mailbox struct {
	msg_arr []Msg
	sync.Mutex
}

func (mailbox *Mailbox) Append(msg Msg) {
	// append a message to the message array in mailbox
	mailbox.Lock()
	defer mailbox.Unlock()
	mailbox.msg_arr = append(mailbox.msg_arr, msg)

}
func (mailbox *Mailbox) PrintWhileLocked(w *tabwriter.Writer) {
	// assuming mailbox mutex is locked, print each item in the array
	// using the tabwriter formatting
	fmt.Fprintln(w, "Type\tFrom\tTo\tTimestamp")
	for _, msg := range mailbox.msg_arr {
		fmt.Fprintf(w, "%d\t%d\t%d\t%v\n", msg.msgtype, msg.from, msg.to, msg.ts)
	}
	w.Flush()
}
