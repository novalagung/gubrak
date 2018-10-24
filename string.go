package gubrak

import (
	randMath "math/rand"
	"regexp"
	"time"
)

func init() {
	randMath.Seed(time.Now().UnixNano())
}

// RandomString function generate random alphabet string in defined length
func RandomString(length int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, length)
	for i := range b {
		b[i] = letters[randMath.Intn(len(letters))]
	}

	return string(b)
}

// ReplaceCaseInsensitive function replace all string that match with `find` without caring about it's case
func ReplaceCaseInsensitive(text, find, replacement string) string {
	re := regexp.MustCompile(`(?i)` + find)
	return re.ReplaceAllString(text, replacement)
}
