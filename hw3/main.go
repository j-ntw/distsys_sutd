package main

type Msg struct {
	// Define your Msg structure fields here
}
type Node struct {
	ch     chan Msg
	id     int
	page   int
	access bool
}
type CM struct {
	ch       chan Msg
	id       int
	page     int
	copy_set Set
	owner    int
}

var (
	numNodes = 3
	// node_arr node_arr [numNodes]Node
	// ch_arr [numNodes]chan Msg

)

func main() {

	// single central manager (with one backup)
	// a few other nodes/processes

	go CM.Run()
	for n := range node_arr {
		n.Run()
	}
}
