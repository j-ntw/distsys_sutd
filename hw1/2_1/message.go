package main

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
	election     Msgtype = iota // election == 0
	normal                      // normal == 1
	coordination                // coordination == 2
	victory                     // victory == 3
	ack
)

func send(ch chan Msg, msg Msg) {
	ch <- msg
}

func (self *Node) SendElectionMsg() {
	// broadcast victory msg to all ids larger than itself
	for i, other_ch := range self.ch_arr[self.id+1:] {
		out_msg := Msg{election, self.id, i, 0, 0}
		go send(other_ch, out_msg)
	}
	self.cmd <- awaiting_ack
}

func (self *Node) SendVictoryMsg() {
	// broadcast election message to all and wait
	for i, other_ch := range self.ch_arr {
		out_msg := Msg{victory, self.id, i, 0, 0}
		go send(other_ch, out_msg)
	}
}
