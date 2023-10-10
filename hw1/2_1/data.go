package main

import "sync"

type Data struct {
	data int
	mu   sync.Mutex
}

func (self *Node) getData() int {

	// get data
	self.data.mu.Lock()
	data := self.data.data
	self.data.mu.Unlock()
	return data
}

func (self *Node) setData(data int) {

	// set data
	self.mode.mu.Lock()
	self.data.data = data
	self.mode.mu.Unlock()
}
