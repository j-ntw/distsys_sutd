
package main

import (
	"fmt"
	"math/rand"
	"time"
)



func sleepRand() {
	//sleep sporadically
	randamt := rand.Intn(1000)
	fmt.Printf("sleeping: %d ms\n", randamt)
	amt := time.Duration(randamt)
	time.Sleep(time.Millisecond * amt)
}
func main() {
	for i := 0; i < 5; i++ {

		sleepRand()
	}
}
