package main

import (
	"SleepRand"
	"fmt"
)

type Node struct {
	ch    chan Msg
	id    int
	queue Queue // request queue
	set   Set   // replies set for own request
	clock VectorClock
	voted bool
}

func NewNode(id int) *Node {
	// create and return a new node
	// other details like coordinator, data and mode are left as default
	// ch and ch_arr are assigned in the main program for loop
	x := &Node{
		id:    id,
		queue: *NewQueue(),
		set:   *NewSet()}
	x.queue.parent_node = x
	return x
}

func (self *Node) Critical() {
	// normal function
	fmt.Printf("n%d: critical\n", self.id)
}

func (self *Node) Broadcast(out_msg Msg) {
	// one off broadcast while coordinating
	fmt.Printf("n%d: do broadcast\n", self.id)
	for i, other_ch := range ch_arr {
		// do not send to self
		if i != self.id {
			out_msg.to = i
			go send(self.id, other_ch, out_msg)
		}
	}
}

func (self *Node) reply(reply_msg Msg) {
	// use as goroutine
	// reply immediately
	self.clock.Inc(self.id)
	reply_msg.ts = self.clock.Get()
	to_ch := ch_arr[reply_msg.to]
	go send(self.id, to_ch, reply_msg)
}

func (self *Node) listen() {
	// listens for msg from other nodes
	// use as goroutine
	fmt.Printf("n%d: start listen\n", self.id)
	for {
		// receive message
		in_msg := <-self.ch
		mailbox.Append(in_msg)
		fmt.Printf("n%d: receive %d %d->%d\n", self.id, in_msg.msgtype, in_msg.from, in_msg.to)

		// increment own vectorclock

		self.clock.AdjustClock(in_msg.ts)
		self.clock.Inc(self.id)

		switch msgtype := in_msg.msgtype; msgtype {
		case req:
			if self.voted {
				// if voted
				self.queue.push(in_msg)
			} else {
				reply_msg := Msg{vote, self.id, in_msg.from, [numNodes]int(zeroVector)}
				self.reply(reply_msg)
				self.voted = true
			}
		case vote:
			self.set.add(in_msg.from)
		case release:
			// TODO: see who else is pending my vote based on queue
			self.queue.pop()
			if self.queue.isEmpty() {
				self.voted = false
			} else {
				reply_msg := Msg{vote, self.id, self.queue.peek().from, [numNodes]int(zeroVector)}
				self.reply(reply_msg)
				self.voted = true
			}
		default:
			fmt.Printf("msgtype: %v", msgtype)
		}
	}
}

func (self *Node) Run() {
	// Run is non blocking
	self.ch = ch_arr[self.id]
	// start listener
	go self.listen()

	// periodically request to enter critical section
	// stamp request
	self.clock.Inc(self.id)
	req_msg := Msg{req, self.id, 0, self.clock.Get()} // placeholder to_id

	// reset reply_set
	self.set.init(self.id)
	self.set.add(self.id)

	// broadcast request message
	self.Broadcast(req_msg)

	// wait for majority of votes
	<-self.set.majority_ch

	fmt.Printf("n%d: execute\n", self.id)
	// execute critical section
	self.Critical()

	// send release message
	release_msg := Msg{release, self.id, 0, self.clock.Get()}
	self.Broadcast(release_msg)

	// sleep before repeating
	SleepRand.SleepRand()

}
