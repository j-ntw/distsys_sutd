package main

import (
	"fmt"
	"math/rand"
)

func coinFlip() bool {
	return rand.Intn(2) == 1
}
func main() {
	for i := 0; i < 5; i++ {
		fmt.Println("Bool: ", coinFlip())
	}
}
