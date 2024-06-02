package plex

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetWatchlist(t *testing.T) {
	c := testClient()
	token := os.Getenv("PLEX_TOKEN")

	list, err := c.GetWatchlist(token)
	assert.NoError(t, err)
	assert.NotNil(t, list)
}

func TestWatchlists_Detail(t *testing.T) {
	c := testClient()
	token := os.Getenv("PLEX_TOKEN")

	list, err := c.GetWatchlist(token)
	assert.NoError(t, err)
	assert.NotNil(t, list)

	details, err := c.GetWatchlistDetail(token, list)
	assert.NoError(t, err)
	assert.NotNil(t, details)
}
