package gubrak

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomString(t *testing.T) {
	assert.NotEqual(t,
		RandomString(32),
		RandomString(32),
	)
}

func TestRandomStringLength(t *testing.T) {
	assert.Len(t, RandomString(32), 32)
}

func TestReplaceCaseInsensitive(t *testing.T) {
	result := ReplaceCaseInsensitive("lOrEm IPsUm DoLor Sit AMEt", "ipsum", "batman")
	assert.Equal(t, result, "lOrEm batman DoLor Sit AMEt")
}
