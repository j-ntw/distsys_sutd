package main

type AccessType int

const (
	ReadOnly AccessType = iota
	WriteOnly
	ReadWrite
)

type MessageType int

const (
	ReadRequest MessageType = iota
	WriteRequest
	ReadForward
	WriteForward
	SendPage
	ReadConfirmation
	WriteConfirmation
	Invalidate
	InvalidateConfirmation
)

type Msg struct {
	// Define your Msg structure fields here
	messageType MessageType
	from        int
	to          int
	page        Page
	accessType  AccessType //AccessType is for sendPage
}
type Page struct {
	id   int // owner process id
	data int
}
type Process struct {
	ch     chan Msg
	id     int
	page   int
	access bool
}
type CM struct {
	ch       chan Msg
	id       int
	page     int
	copy_set Set
	owner    int
	pages    []Page
}

// CM maintains record for all pages
// the owner and all its copies
// TODO: does each process own multiple pages?
// only one process can own a page at a time (no co owner, ownership can be passed around)
// each page can have at most one owner. there can be page that is orphaned

var (
	numprocesss = 3
	// process_arr process_arr [numprocesss]process
	// ch_arr [numprocesss]chan Msg

)

func main() {

	// single central manager (with one backup)
	// a few other processes
	// P3 wants to read page x1 (page fault at P3)
	// P3 sends read req to CM (X1. P1)
	// CM checks page owner and sends read forward to owner P1, adds P3 in copy set

	go CM.Run()
	for n := range process_arr {
		n.Run()
	}
}
