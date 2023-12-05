package main

import (
	"fmt"
	"sort"
	"sync"
)

// represents the queue of replies at each node.
// When your current request is at the head of the queue, one of the conditions for CS is fulfilled.
type Queue struct {
	q           []Msg
	parent_node *Node
	watch_map   map[Msg]bool // contains foreign request messages we need to reply to
	sync.Mutex
}

func NewQueue() *Queue {
	// create and return a new queue
	return &Queue{
		watch_map: make(map[Msg]bool),
	}

}

func (self *Queue) push(msg Msg) {
	// add any message to the back of the queue (atomic?)
	self.Lock()
	defer self.Unlock()
	fmt.Printf("n%d push msg\n", self.parent_node.id)
	self.q = append(self.q, msg)

	// sort/re-prioritise requests in queue based on timestamp
	sort.SliceStable(self.q, func(i, j int) bool {
		return IsBefore(self.q[i], self.q[j])
	})
}

func (self *Queue) pop() Msg {
	self.Lock()
	defer self.Unlock()
	if len(self.q) > 0 {
		val := self.q[0]
		self.q = self.q[1:]
		return val
	}
	return Msg{}

}

func (self *Queue) peek() Msg {
	self.Lock()
	defer self.Unlock()
	if len(self.q) > 0 {
		return self.q[0]
	}
	return Msg{}
}

func (self *Queue) isEmpty() bool {
	self.Lock()
	defer self.Unlock()
	return len(self.q) == 0

}
