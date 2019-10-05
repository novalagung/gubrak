package gubrak

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChainnableMap(t *testing.T) {
	type Sample struct {
		EbookName      string
		DailyDownloads int
	}

	data := []Sample{
		{EbookName: "clean code", DailyDownloads: 10000},
		{EbookName: "rework", DailyDownloads: 12000},
		{EbookName: "detective comics", DailyDownloads: 11500},
	}

	op := New().
		From(data).
		Filter(func(each Sample) bool {
			return each.DailyDownloads >= 10000
		}).
		Map(func(each Sample) string {
			return each.EbookName
		})

	assert.False(t, op.IsError())
	assert.Nil(t, op.Error())

	assert.Nil(t, err)
	assert.EqualValues(t, []string{"clean code", "rework", "detective comics"}, newData)
}
