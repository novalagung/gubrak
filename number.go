package gubrak

import (
	"math/rand"
	"time"
)

// RandomInt function generates random numeric data between specified min and max
func RandomInt(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn((max-min)+1) + min
}
