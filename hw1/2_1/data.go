package main

import "sync"

type Data struct {
	data int
	mu   sync.Mutex
}

func (self *Node) getData() int {

	// get data
	self.data.mu.Lock()
	defer self.data.mu.Unlock()
	return self.data.data
}

func (self *Node) setData(data int) {

	// set data
	self.data.mu.Lock()
	defer self.data.mu.Unlock()
	self.data.data = data
}
