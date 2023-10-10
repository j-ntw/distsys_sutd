package main

import (
	"fmt"
	"time"
)

type Command int

const (
	down Command = iota
	coordinating
	following
	electing
)

type Node struct {
	ch          chan Msg
	ch_arr      [numNodes]chan Msg
	id          int
	coordinator int
	data        Data
	mode        Mode
	cmd         chan Command
}

func NewNode(id int) *Node {
	return &Node{}
}

func (self *Node) Bully() {
	if self.id == (len(self.ch_arr) - 1) {
		self.cmd <- coordinating // todo other mode like victory
		// broadcast victory msg to all
		self.SendVictoryMsg()
	} else {
		// broadcast election message to all and wait
		self.SendElectionMsg()
	}
}

func (self *Node) Broadcast() {
	// continuous broadcast while coordinating
	for {
		if !self.isMode(coordinating) {
			return
		}
		data := self.getData()
		for i, other_ch := range self.ch_arr {
			out_msg := Msg{normal, self.id, i, 0, data}
			other_ch <- out_msg
		}
		// sleep periodically
		time.Sleep(period * time.Millisecond)
	}
}

func (self *Node) listen() {
	for {
		select {
		case in_msg := <-self.ch:

			switch msg_mode := in_msg.msgtype; msg_mode {
			case normal:
				if self.isMode(following) {
					self.setData(in_msg.data)
				} else {
					fmt.Printf("n%d: mode%d, why rx msg%d\n", self.id, coordinating, normal)
					panic(self.id)
				}
			case election:
				if self.id > in_msg.from {
					// if i am bigger send ack to sender
					// start bully
					ack_msg := Msg{ack, self.id, in_msg.from, 0, 0}
					self.ch_arr[in_msg.from] <- ack_msg
					self.Bully()
				} else {
					// if i am smaller, update coordinator to sender.id
					self.coordinator = in_msg.from
				}
			case ack:
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
		// only iterates on next command
		next_mode := <-self.cmd
		switch next_mode {
		case down:
			return
		case coordinating:
			self.coordinator = self.id
			go self.Broadcast()

		case following:
			// no op, handled in self.listen()
		case electing:
			self.Bully()
		default:
			fmt.Printf("n%d: Unknown mode, shutting down...\n", self.id)
			self.cmd <- down
		}
	}
}

func (self *Node) Boot() {
	fmt.Printf("n%d: Booting...\n", self.id)
	self.ch = self.ch_arr[self.id]
	self.cmd <- electing
	go self.Run()
}
