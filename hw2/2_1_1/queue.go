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
	parent_node          *Node
	watch_map            map[Msg]bool // contains foreign request messages we need to reply to
	sync.Mutex
}

func NewQueue() *Queue {
	// create and return a new queue
	return &Queue{
		q_own_req_at_head_ch: make(chan bool),
		watch_map:            make(map[Msg]bool),
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

	self.checkHeadWhileLocked()
}

func (self *Queue) pop() Msg {
	self.Lock()
	defer self.Unlock()
	if len(self.q) > 0 {
		val := self.q[0]
		self.q = self.q[1:]
		if len(self.q) > 0 {
			self.checkHeadWhileLocked()
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
func (self *Queue) watch(req_msg Msg) {
	// after a foreign request is pushed to priority queue
	// watch for a particular message to be head of queue and
	// replies to the sender's channel.
	self.Lock()
	defer self.Unlock()
	self.watch_map[req_msg] = true
}

func (self *Queue) checkHeadWhileLocked() {
	// run when queue is modified
	// unsafe
	fmt.Printf("n%d: queue%v\n", self.parent_node.id, self.q)
	fmt.Printf("n%d: set%v\n", self.parent_node.id, self.watch_map)
	// if head of queue is my own message
	if self.q[0].from == self.parent_node.id {
		// notify via ch
		if len(self.q_own_req_at_head_ch) == 0 {
			go func() {
				self.q_own_req_at_head_ch <- true // todo why block here?
				fmt.Printf("n%d done notify msg\n", self.parent_node.id)
			}()
		}

		// else if head of queue is in watch map (some foreign message)

	} else if self.watch_map[self.q[0]] {
		// check queue if we can reply anyone
		fmt.Printf("n%d can reply someone\n", self.parent_node.id)
		// create and send reply

		reply_msg := Msg{reply, self.parent_node.id, self.q[0].from, [numNodes]int(zeroVector)}
		fmt.Printf("\n\n\n%v\n\n\n", reply_msg)
		go self.parent_node.reply(reply_msg)

		// delete from watch map
		delete(self.watch_map, self.q[0])
	} else {
		fmt.Printf("n%d: head is some other msg\n", self.parent_node.id)
	}
}
