package gubrak

import (
	randMath "math/rand"
)

// RandomInt function generates random numeric data between specified min and max
func RandomInt(min, max int) int {
	return randMath.Intn((max-min)+1) + min
}
