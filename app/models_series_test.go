package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnector_SeriesAllUnwatched(t *testing.T) {
	c := testConnector()
	if c == nil {
		t.Skip("No test connector")
		return
	}

	series := &Series{}
	err := c.Series.Find("644d88003359bb748dd63096", series)
	assert.NoError(t, err, "Find Series")

	got, err := c.SeriesUserUnwatched(series)
	assert.NoError(t, err, "unwatched")
	assert.Greater(t, got, 0, "unwatched")
}
