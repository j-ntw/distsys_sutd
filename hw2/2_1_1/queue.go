package main

import (
	"fmt"
	"sync"
)

// represents the queue of replies at each node.
// When your current request is at the head of the queue, one of the conditions for CS is fulfilled.
type Queue struct {
	q                 []Msg
	q_own_req_at_head chan bool
	q_empty_ch        chan bool
	own_req           Msg
	mu                sync.Mutex
}

func NewQueue() *Queue {
	// create and return a new queue
	return &Queue{
		q:                 make([]Msg, numNodes),
		q_own_req_at_head: make(chan bool),
		q_empty_ch:        make(chan bool)}
}

func (self *Queue) push(msg Msg) {
	self.mu.Lock()
	defer self.mu.Unlock()
	fmt.Printf("n: pushed msg\n")
	self.q = append(self.q, msg)

	// check if own req at head
	if self.q[0] == self.own_req {
		//TODO maybe use counters instead, channels might be blocking on send
		self.q_own_req_at_head <- true
	}
}
func (self *Queue) pop() Msg {
	self.mu.Lock()
	defer self.mu.Unlock()
	if len(self.q) > 0 {
		val := self.q[0]
		self.q = self.q[1:]

		// check if q empty
		if len(self.q) == 0 {
			self.q_empty_ch <- true
		}
		// check if own req at head
		if self.q[0] == self.own_req {
			self.q_own_req_at_head <- true
		}
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
