package gubrak

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
