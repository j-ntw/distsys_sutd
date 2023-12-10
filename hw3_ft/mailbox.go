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
func (mailbox *Mailbox) print(w *tabwriter.Writer) {
	mailbox.Lock()
	defer mailbox.Unlock()
	// print each item in the array
	// using the tabwriter formatting
	fmt.Fprintln(w, "Type\tFrom\tTo\tTimestamp")
	for _, msg := range mailbox.msg_arr {
		fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%d\n", msg.msgtype.String(), msg.from, msg.to, msg.page_no, msg.page_no)
	}
	w.Flush()
}
