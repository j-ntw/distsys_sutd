package main

import "sync"

// represents the set of replies we need to receive for a given request.
// When the set is empty, one of the conditions for critical section is fulfilled.

type Set struct {
	s          map[int]bool
	s_empty_ch chan bool
	sync.Mutex
}

func NewSet() *Set {
	// create and return a new set
	return &Set{
		s:          make(map[int]bool),
		s_empty_ch: make(chan bool)}
}

func (self *Set) del(k int) {
	self.Lock()
	defer self.Unlock()
	delete(self.s, k)
	if len(self.s) == 0 {
		self.s_empty_ch <- true
	}
}

func (self *Set) add(k int) {
	self.Lock()
	defer self.Unlock()
	self.s[k] = true
}

func (self *Set) isEmpty() {
	self.Lock()
	defer self.Unlock()
	if len(self.s) == 0 {
		self.s_empty_ch <- true
	}
}

func (self *Set) init(ignore int) {
	self.Lock()
	defer self.Unlock()
	for i := 0; i < numNodes; i++ {
		if i != ignore {
			self.s[i] = true
		}

	}
}
