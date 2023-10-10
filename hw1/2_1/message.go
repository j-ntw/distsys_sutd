package main

type Msg struct {
	Msgtype int
	from    int
	to      int
	ts      int
	data    int
}

// message types
const ( // iota is reset to 0
	election     = iota // election == 0
	normal       = iota // normal == 1
	coordination = iota // coordination == 2
	victory      = iota // victory == 3
)
