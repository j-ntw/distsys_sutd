package main

type AccessType int

const (
	ReadOnly AccessType = iota
	WriteOnly
	ReadWrite
	Nil
)

// the page struct represents metadata about the page that a process knows
// a page can have one owner
// each process can own multiple pages
// an array of Page structs form a pagetable, with the index of each page as the page number.
type Page struct {
	isOwner  bool
	isLocked bool
	access   AccessType
}

func newPage() *Page {
	return &Page{
		isOwner:  false,
		isLocked: false,
		access:   Nil,
	}
}

type CM_Record struct {
	owner_id int
	copy_set map[int]bool // set of process_ids who have read only copies
	isLocked bool
}
