package sleeprand

import (
	"math/rand"
	"time"
)

func SleepRand() {
	// sleep sporadically for [1,1000] ms
	randamt := rand.Intn(1000) + 1
	amt := time.Duration(randamt)
	time.Sleep(time.Millisecond * amt)
}
