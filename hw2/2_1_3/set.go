package main

import "sync"

// represents the set of replies we need to receive for a given request.
// When the set is empty, one of the conditions for critical section is fulfilled.

type Set struct {
	s           map[int]bool
	majority_ch chan bool
	_isCritical bool
	sync.Mutex
}

func NewSet() *Set {
	// create and return a new set
	return &Set{
		_isCritical: false,
		s:           make(map[int]bool),
		majority_ch: make(chan bool)}
}

func (self *Set) del(k int) {
	self.Lock()
	defer self.Unlock()
	delete(self.s, k)
	if len(self.s) < majority {
		self._isCritical = false
	}
}

func (self *Set) add(k int) {
	self.Lock()
	defer self.Unlock()
	self.s[k] = true
	if len(self.s) >= majority {
		self.majority_ch <- true
		self._isCritical = true
	}
}
func (self *Set) isCritical() bool {
	self.Lock()
	self.Unlock()
	return self._isCritical
}

func (self *Set) isEmpty() {
	self.Lock()
	defer self.Unlock()
}

func (self *Set) init(ignore int) {
	self.Lock()
	defer self.Unlock()
	for key := range self.s {
		delete(self.s, key)
	}
}
