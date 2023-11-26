package main

import "fmt"

func (cm *CM) onReceiveReadConfirmation(msg Message) {
	//	TODO: not sure what to do here
	fmt.Printf("[CM] Received Read Confirmation \n")
}

func (cm *CM) onInvalidateConfirmation(msg Message) {
	fmt.Printf("[CM] Received Read Confirmation \n")
}
