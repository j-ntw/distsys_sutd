package main

import (
	"fmt"
)

type Node struct {
	ch     chan Msg
	ch_arr [numNodes]chan Msg
	id     int
	queue  Queue // request queue
	set    Set
	clock  VectorClock
}

func NewNode(id int, own_ch chan Msg) *Node {
	// create and return a new node
	// other details like coordinator, data and mode are left as default
	// ch and ch_arr are assigned in the main program for loop
	return &Node{
		id:    id,
		queue: *NewQueue(own_ch),
		set:   *NewSet()}
}
func (self *Node) Critical() {
	// normal function
	fmt.Printf("n%d: critical", self.id)
}

func (self *Node) reply(resp_msg Msg) {
	// use as goroutine

	// wait until head of queue is the request message we just pushed
	self.queue.peek()
	<-self.queue.q_empty_ch

	// reply immediately
	to_ch := self.ch_arr[resp_msg.to]
	go send(to_ch, resp_msg)

}

func (self *Node) Broadcast(out_msg Msg) {
	// one off broadcast while coordinating
	fmt.Printf("n%d: do broadcast\n", self.id)
	for i, other_ch := range self.ch_arr {
		// do not send to self
		if i != self.id {
			out_msg.to = i
			go send(other_ch, out_msg)
		}
	}

}

func (self *Node) listen() {
	// listens for msg from other nodes
	// use as goroutine
	fmt.Printf("n%d: start listen\n", self.id)
	for {
		// receive message
		in_msg := <-self.ch
		fmt.Printf("receive %d->%d\n", in_msg.from, in_msg.to)
		// increment own vectorclock
		self.clock.Inc(self.id)
		fmt.Printf("n%d: execute\n", self.id)
		switch msgtype := in_msg.msgtype; msgtype {
		case req:
			// add req to own queue
			self.queue.push(in_msg)
			resp_msg := Msg{resp, self.id, in_msg.from, [numNodes]int(zeroVector)}
			// if pending replies, hold reply
			go self.try_reply(resp_msg)
		case resp:
			self.set.del(in_msg.from)
		case release:
			self.queue.pop()
		case own_req_at_q_head:
			// internal messaging to self
			if len(self.queue.q_own_req_at_head_ch) == 0 {
				// its ok to check length then push because listen is only called once and runs in a linear loop
				self.queue.q_own_req_at_head_ch <- true
			}

		default:
			fmt.Printf("msgtype: %v", msgtype)
		}

	}
}

func (self *Node) Run() {
	// Run is non blocking
	self.ch = self.ch_arr[self.id]
	// start listener
	go self.listen()

	// main node loop
	for {
		// periodically request to enter critical section
		// stamp request
		self.clock.Inc(self.id)
		req_msg := Msg{req, self.id, 0, self.clock.Get()}
		// add to own queue
		self.queue.own_req = req_msg
		self.queue.push(req_msg)

		fmt.Printf("n%d: add to own q\n", self.id)

		// reset reply_set
		self.set.init(self.id)

		// broadcast request message
		self.Broadcast((req_msg))
		fmt.Printf("n%d: broadcast\n", self.id)

		// block while waiting for replies, waiting for own request to pop
		<-self.set.s_empty_ch
		<-self.queue.q_own_req_at_head_ch

		fmt.Printf("n%d: execute\n", self.id)
		// execute critical section
		self.Critical()

		// exit critical section
		self.queue.pop()

		// send release message
		release_msg := Msg{release, self.id, 0, self.clock.Get()}
		self.Broadcast(release_msg)

		// sleep before repeating
		// SleepRand.SleepRand()
		// fmt.Printf("n%d: bye\n", self.id)
		// break

	}
}
