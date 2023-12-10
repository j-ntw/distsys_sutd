package main

type AccessType int

const (
	ReadOnly AccessType = iota
	ReadWrite
	Nil
)

// the page struct represents metadata about the page that a process knows
// a page can have one owner
// each process can own multiple pages
// an array of Page structs form a pagetable, with the index of each page as the page number.
type Page struct {
	isOwner bool

	access AccessType
}

func (a AccessType) String() string {
	return [...]string{"ReadOnly", "ReadWrite", "Nil"}[a]
}

func newPage(isOwner bool) *Page {
	// if you are the owner, you know that you have readwrite access
	// if you are not the owner, upon initialising you dont know what access you have == you dont have access
	return &Page{
		isOwner: isOwner,

		access: func() AccessType {
			if isOwner {
				return ReadWrite
			} else {
				return Nil
			}
		}(),
	}
}

