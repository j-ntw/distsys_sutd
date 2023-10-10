package main

import (
	"fmt"
	"time"
)

const (
	down = iota
	coordinating
	following
	electing
)

type Node struct {
	ch     chan Msg
	ch_arr [numNodes]chan Msg
	id     int
	data   Data
	mode   Mode
}

func NewNode(id int) *Node {
	return &Node{}
}

func (self *Node) Bully() {

	panic("unimplemented")
}

func (self *Node) Elect() {
	if self.id == (len(self.ch_arr) - 1) {
		self.setNextMode(coordinating) //todo other mode like victory
		// send victory msg
	} else {
		// send election message and wait
		self.SendElectionMsg()
	}
}

func (self *Node) SendElectionMsg() {

}

func (self *Node) SingleBroadcast() {
	data := self.getData()
	for i, other_ch := range self.ch_arr {
		// todo how to detect liveness w timeout when sending?
		out_msg := Msg{normal, self.id, i, 0, data}
		other_ch <- out_msg
	}

	// sleep periodically
	time.Sleep(period * time.Millisecond)

}
func (self *Node) listen() {
	for {
		select {
		case in_msg := <-self.ch:

			switch msg_mode := in_msg.Msgtype; msg_mode {
			case normal:
				if self.isMode(following) {
					self.setData(in_msg.data)
				} else {
					fmt.Printf("n%d: mode%d, why rx msg%d\n", self.id, coordinating, normal)
					panic(self.id)
				}

			case election:
				// check id, if lower then Bully else transmit victory message
				self.Bully() // >:)
			// todo other msgtypes

			default:
				fmt.Printf("msg_mode: %v", msg_mode)
			}

		case <-time.After(timeout * time.Millisecond):
			// start election
			self.Bully()
			// no default
		}
	}
}
func (self *Node) Run() {
	// start listener
	go self.listen()

	// main node loop
	for {
		switch next_mode := self.getNextMode(); next_mode {
		case down:
			self.setMode(down)
			return
		case coordinating:
			if self.getMode() != next_mode {
				self.setMode(next_mode)
			}
			self.SingleBroadcast() // blocking
		case following:
			if self.getMode() != next_mode {
				self.setMode(next_mode)
				//start following
				// self.sync()
			}
			// was already following
		case electing:
			self.Elect()
		default:
			fmt.Printf("n%d: Unknown mode, shutting down...\n", self.id)
			self.setMode(down)
		}

	}

}

func (self *Node) Boot() {
	fmt.Printf("n%d: Booting...\n", self.id)
	self.setNextMode(electing)
	go self.Run()
}
