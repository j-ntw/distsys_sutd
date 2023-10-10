package main

import "sync"

type Mode struct {
	mode Command // current mode
	mu   sync.Mutex
}

func (self *Node) isMode(mode Command) bool {

	// check mode
	self.mode.mu.Lock()
	state := self.mode.mode == mode
	self.mode.mu.Unlock()
	return state

}

func (self *Node) getMode() Command {

	// get mode
	self.mode.mu.Lock()
	mode := self.mode.mode
	self.mode.mu.Unlock()
	return mode
}

func (self *Node) setMode(mode Command) {

	// set mode
	self.mode.mu.Lock()
	self.mode.mode = mode
	self.mode.mu.Unlock()
}
