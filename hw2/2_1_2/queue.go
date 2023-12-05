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

func (self *Queue) pushIfPending(msg Msg) bool {
	// returns intended recepient id if replying now, else return -1
	self.Lock()
	defer self.Unlock()
	// if no other reply is pending, reply now
	if len(self.q) == 0 {
		return true
	}
	// if the head of queue is prior to msg aka pending
	if IsBefore(self.q[0], msg) {
		// sort/re-prioritise requests in queue based on timestamp
		self.q = append(self.q, msg)

		// sort/re-prioritise requests in queue based on timestamp
		sort.SliceStable(self.q, func(i, j int) bool {
			return IsBefore(self.q[i], self.q[j])
		})
		return false
	} else {
		// not pending, reply neow
		return true
	}
}

func (self *Queue) walk() {
	self.Lock()
	defer self.Unlock()
	for i := range self.q {
		// not pending, reply neow
		reply_msg := Msg{reply, self.parent_node.id, self.q[i].from, [numNodes]int(zeroVector)}
		self.parent_node.reply(reply_msg)

	}
	self.q = self.q[:0] // reslice to empty
}
