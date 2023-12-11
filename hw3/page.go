package main

import "sync"

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
	access  AccessType
}

type Pages struct {
	pages []Page
	sync.RWMutex
}

func (pt *Pages) Get() []Page {
	pt.Lock()
	defer pt.Unlock()
	return pt.pages
}

func (pt *Pages) SetOwner(page_no int, isOwner bool) {
	pt.Lock()
	defer pt.Unlock()
	pt.pages[page_no].isOwner = isOwner
}

func (pt *Pages) SetAccess(page_no int, access AccessType) {
	pt.Lock()
	defer pt.Unlock()
	pt.pages[page_no].access = access
}

func newPageTable(pages []Page) *Pages {
	return &Pages{pages: pages}
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
