package main

import "fmt"

type Msg struct {
	msgtype Msgtype
	from    int
	to      int
	ts      int
	data    int
}

// message types
type Msgtype int

const ( // iota is reset to 0
	election     Msgtype = iota //  0
	broadcast                   //  1
	coordination                //  2
	victory                     //  3
	ack                         //  4
)

func send(ch chan Msg, msg Msg) {
	ch <- msg
}

func (self *Node) SendElectionMsg() {
	// broadcast victory msg to all ids larger than itself
	fmt.Printf("n%d: SendElectionMsg\n", self.id)
	// do not send to self
	for i, other_ch := range self.ch_arr[self.id+1:] {
		out_msg := Msg{election, self.id, i, 0, 0}
		go send(other_ch, out_msg)

	}
	self.cmd <- awaiting_ack
}

func (self *Node) SendVictoryMsg() {
	// broadcast election message to all and wait
	fmt.Printf("n%d: SendVictoryMsg\n", self.id)
	for i, other_ch := range self.ch_arr {
		// do not send to self
		if i != self.id {
			out_msg := Msg{victory, self.id, i, 0, 0}
			go send(other_ch, out_msg)
		}
	}
}
