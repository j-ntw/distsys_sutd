package main

import (
	"fmt"
	"slices"
	"time"
)

type Command int

const (
	down Command = iota
	coordinating
	following
	electing
	awaiting_ack
	awaiting_victory
	declaring_victory
)

var (
	bullystates = []Command{electing,
		awaiting_ack,
		awaiting_victory,
		declaring_victory}
)

type Node struct {
	ch          chan Msg
	ch_arr      [numNodes]chan Msg
	id          int
	coordinator int
	data        Data
	mode        Mode
	cmd         chan Command
	trigger_ch  chan bool
}

func NewNode(id int) *Node {
	return &Node{
		id:         id,
		trigger_ch: make(chan bool),
		cmd:        make(chan Command)}
}

func (self *Node) Bully() {
	if self.id == (len(self.ch_arr) - 1) {
		self.cmd <- coordinating
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
		// if change mode, then break this loop
		if !self.isMode(coordinating) {
			fmt.Printf("n%d: stop broadcast\n", self.id)
			return
		}
		// else coordinating
		data := self.getData()
		fmt.Printf("n%d: do broadcast\n", self.id)
		for i, other_ch := range self.ch_arr {
			// do not send to self
			if i != self.id {
				out_msg := Msg{broadcast, self.id, i, 0, data}
				go send(other_ch, out_msg)
			}
		}
		// sleep periodically
		time.Sleep(period * time.Millisecond)
	}
}

func (self *Node) listen() {
	// listens for msg from other nodes
	fmt.Printf("n%d: start listen\n", self.id)
	for {
		select {
		case in_msg := <-self.ch:

			switch msgtype := in_msg.msgtype; msgtype {
			case broadcast:
				if self.isMode(following) {
					// fmt.Printf("n%d: mode%d, why rx msg%d\n", self.id, coordinating, broadcast)
					self.setData(in_msg.data)
				} else {
					fmt.Printf("n%d: mode%d, why rx msg%d from %d\n", self.id, coordinating, broadcast, in_msg.from)
					panic(self.id)
				}

			case election:
				if self.id > in_msg.from {
					// if i am bigger send ack to sender
					// start bully if not already bullying
					ack_msg := Msg{ack, self.id, in_msg.from, 0, 0}
					self.ch_arr[in_msg.from] <- ack_msg
					if !slices.Contains(bullystates, self.getMode()) {
						self.cmd <- electing
					}
				}

			case ack:
				// stop election process, dont proceed to declaring_victory
				// subsequent ack messages are consumed but no op
				if self.isMode(awaiting_ack) {
					self.trigger_ch <- true
				}

			case victory:
				// trigger early exit of timeout
				if self.isMode(awaiting_victory) {
					self.trigger_ch <- false
				}
				self.coordinator = in_msg.from
				self.cmd <- following

			default:
				fmt.Printf("msgtype: %v", msgtype)
			}

		case <-time.After(timeout * time.Millisecond):
			// start election
			self.cmd <- electing
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
		fmt.Printf("n%d: %d->%d \n", self.id, self.getMode(), next_mode)
		self.setMode(next_mode)

		switch next_mode {

		case down:
			return
		case coordinating:
			self.coordinator = self.id
			go self.Broadcast()

		case following:
			// no op, handled in self.listen()
		case electing:
			go self.Bully()
		case awaiting_ack:
			// wait for timeouts
			select {
			case value := <-self.trigger_ch:
				// todo check if a victory message came in while i tried to bully someone
				if value {
					self.cmd <- awaiting_victory
				} else {
					// no op, victory case in self.listen() has sent "following" command
					// triggered to avoid timeout
				}
			case <-time.After(timeout * time.Millisecond):
				// if timeout, declare self as victor
				self.cmd <- declaring_victory
			}

		case awaiting_victory:
			// wait for timeouts
			select {
			case <-self.trigger_ch:
				// no op, handled in self.listen()
				// triggered to avoid timeout

			case <-time.After(timeout * time.Millisecond):
				// if timeout, restart election process
				self.cmd <- electing
			}

		case declaring_victory:
			self.SendVictoryMsg()

		default:
			fmt.Printf("n%d: Unknown mode, shutting down...\n", self.id)
			self.cmd <- down
		}
	}
}

func (self *Node) Boot() {
	fmt.Printf("n%d: Booting...\n", self.id)
	self.ch = self.ch_arr[self.id]
	// clear channel
	// L:
	// 	for {
	// 		select {
	// 		case _, ok := <-self.ch:
	// 			if !ok { //ch is closed //immediately return err
	// 				break L
	// 			}
	// 		default: //all other case not-ready: means nothing in ch for now
	// 			break L
	// 		}
	// 	}
	fmt.Printf("n%d: electing...\n", self.id)
	go func() {
		self.cmd <- electing
	}()

	fmt.Printf("n%d: running...\n", self.id)
	go self.Run()
}
