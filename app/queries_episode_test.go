package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnector_Upcoming(t *testing.T) {
	c := testConnector()

	got, err := c.Upcoming()
	assert.NoError(t, err, "Upcoming")
	assert.Greater(t, len(got), 0, "Upcoming")
}
