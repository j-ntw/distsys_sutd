package CoinFlip

import "math/rand"

func CoinFlip() bool {
	return rand.Intn(2) == 1
}
