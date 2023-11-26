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
	// Read Protocol
	// 1. P3 wants to read page x1 (page fault at P3)
	// 2. P3 sends read req to CM (X1. P1)
	// 3. CM checks page owner and sends read forward to owner P1, adds P3 in copy set
	// 4. P3 sends Read confirmation to central manager

	// Write Protocol
	// P2 encounters page fault while writing page X1
	// 1. P2 sends a write request to the CM
	// 2. P3 sends confirmation of invalidation msg to CM
	// 3. CM receives invalidate cfm and clears copy set
	// 4. CM sends write forward for page X1 to P1.
	// 5. P1 sends X1 to P2, invalidates own copy of X1.
	// 6. P2 sends write confirmation for X1 to CM.
	go CM.Run()
	for n := range process_arr {
		n.Run()
	}
}
