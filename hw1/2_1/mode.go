package main

import "sync"

type Mode struct {
	mode      int // current mode
	next_mode int
	mu        sync.Mutex
}

func (self *Node) isMode(mode int) bool {

	// check mode
	self.mode.mu.Lock()
	state := self.mode.mode == mode
	self.mode.mu.Unlock()
	return state

}
func (self *Node) getMode() int {

	// get mode
	self.mode.mu.Lock()
	mode := self.mode.mode
	self.mode.mu.Unlock()
	return mode
}
func (self *Node) getNextMode() int {

	// get mode
	self.mode.mu.Lock()
	next_mode := self.mode.next_mode
	self.mode.mu.Unlock()
	return next_mode
}
func (self *Node) setMode(mode int) {

	// set mode
	self.mode.mu.Lock()
	self.mode.mode = mode
	self.mode.mu.Unlock()
}

func (self *Node) setNextMode(next_mode int) {

	// get mode
	self.mode.mu.Lock()
	self.mode.next_mode = next_mode
	self.mode.mu.Unlock()
}
