package main

import (
	"fmt"
	"time"
)

const (
	maxGoroutines = 10
)

// without race condition
func main() {
	for i := 0; i < maxGoroutines; i++ {
		go func(mindex int) {
			fmt.Println(mindex)
		}(i)
	}
	time.Sleep(time.Second)
}
