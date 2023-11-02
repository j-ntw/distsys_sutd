package main

import "sync"
// represents the set of replies we need to receive for a given reply.
// When the set is empty, one of the conditions for critical section is fulfilled.

type Set struct {
	s  map[int]bool
	mu sync.Mutex
}

func NewSet() *Set {
	// create and return a new set
	return &Set{
		s: make(map[int]bool)}
}

func (self *Set) del(k int) {
	self.mu.Lock()
	defer self.mu.Unlock()
	delete(self.s, k)
}

func (self *Set) add(k int) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.s[k] = true
}

func (self *Set) isEmpty() bool {
	self.mu.Lock()
	defer self.mu.Unlock()
	return len(self.s) == 0
}

//
func (self *Set) init(ignore int) {
	self.mu.Lock()
	defer self.mu.Unlock()
	for i := 0; i < numNodes; i++ {
		if i != ignore {
			self.s[i] = true
		}

	}
}
