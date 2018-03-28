package test

import (
	"errors"
	. "github.com/novalagung/gubrak"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDate(t *testing.T) {
	var result interface{} = Now()

	if _, ok := result.(time.Time); !ok {
		assert.Error(t, errors.New("Now() function is not working"))
	}
}
