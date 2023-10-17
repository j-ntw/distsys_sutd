package main

import "sync"

type Mode struct {
	mode Command // current mode
	mu   sync.Mutex
}

func (self *Node) isMode(mode Command) bool {

	// check mode
	self.mode.mu.Lock()
	defer self.mode.mu.Unlock()
	return self.mode.mode == mode

}

func (self *Node) getMode() Command {

	// get mode
	self.mode.mu.Lock()
	defer self.mode.mu.Unlock()
	return self.mode.mode
}

func (self *Node) setMode(mode Command) {

	// set mode
	self.mode.mu.Lock()
	defer self.mode.mu.Unlock()
	self.mode.mode = mode

}
