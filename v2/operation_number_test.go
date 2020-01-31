package gubrak

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomInt(t *testing.T) {
	result := RandomInt(12, 13)
	assert.True(t, result >= 12 && result <= 13)
}
