package main

import (
	"fmt"
	"sort"
	"sync"
)

// represents the queue of replies at each node.
// When your current request is at the head of the queue, one of the conditions for CS is fulfilled.
type Queue struct {
	q                    []Msg
	q_own_req_at_head_ch chan bool
	q_empty_ch           chan bool
	own_req              Msg
	parent_node_ch       chan Msg
	sync.Mutex
}

func NewQueue(parent_node_ch chan Msg) *Queue {
	// create and return a new queue
	return &Queue{
		q:                    make([]Msg, numNodes),
		q_own_req_at_head_ch: make(chan bool),
		q_empty_ch:           make(chan bool),
		parent_node_ch:       parent_node_ch}
}

func (self *Queue) push(msg Msg) {
	self.Lock()
	defer self.Unlock()
	fmt.Printf("n: pushed msg\n")
	self.q = append(self.q, msg)

	// check if own req at head
	if self.q[0] == self.own_req {
		self.parent_node_ch <- Msg{msgtype: own_req_at_q_head}
	}
}
func (self *Queue) pop() Msg {
	self.Lock()
	defer self.Unlock()
	if len(self.q) > 0 {
		val := self.q[0]
		self.q = self.q[1:]

		// re-prioritise requests in queue based on timestamp
		sort.Slice(self.q[:], func(i, j int) bool {
			return IsBefore(self.q[i].ts, self.q[j].ts)
		})

		// check if own req at head
		if self.q[0] == self.own_req {
			self.parent_node_ch <- Msg{msgtype: own_req_at_q_head}
		}
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

func atHead(){
	// run when queue is modified
	// each message in the queue is marked with a 
}
