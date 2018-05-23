package gubrak

import (
	randMath "math/rand"
	"time"
)

func init() {
	randMath.Seed(time.Now().UnixNano())
}

// RandomStringAlphabet function generate random alphabet string in defined length
func RandomString(length int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, length)
	for i := range b {
		b[i] = letters[randMath.Intn(len(letters))]
	}

	return string(b)
}
