package main

import "sync"

// represents the queue of replies at each node.
// When your current request is at the head of the queue, one of the conditions for CS is fulfilled.
type Queue struct {
	q  []Msg
	mu sync.Mutex
}

func NewQueue() *Queue {
	// create and return a new queue
	return &Queue{
		q: make([]Msg, numNodes)}
}

func (self *Queue) push(msg Msg) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.q = append(self.q, msg)
}
func (self *Queue) pop() Msg {
	self.mu.Lock()
	defer self.mu.Unlock()
	if len(self.q) > 0 {
		val := self.q[0]
		self.q = self.q[1:]
		return val
	}
	return Msg{}

}
func (self *Queue) peek() Msg {
	self.mu.Lock()
	defer self.mu.Unlock()
	if len(self.q) > 0 {
		return self.q[0]
	}
	return Msg{}
}
