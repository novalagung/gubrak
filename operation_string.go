package gubrak

import (
	"crypto/rand"
	"math/big"
	"regexp"
)

// RandomString function generate random alphabet string in defined length
func RandomString(length int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, length)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		b[i] = letters[n.Int64()]
	}

	return string(b)
}

// ReplaceCaseInsensitive function replace all string that match with `find` without caring about it's case
func ReplaceCaseInsensitive(text, find, replacement string) string {
	re := regexp.MustCompile(`(?i)` + find)
	return re.ReplaceAllString(text, replacement)
}
