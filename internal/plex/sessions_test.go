package plex

import (
	"fmt"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHistory(t *testing.T) {
	c := testClient()
	list, err := c.GetHistory()
	assert.NoError(t, err)
	assert.NotNil(t, list)
	slices.Reverse(list)
	for _, h := range list {
		fmt.Printf("%s %d %s - %s\n", h.RatingKey, h.AccountID, h.GrandparentTitle, h.Title)
	}
}
