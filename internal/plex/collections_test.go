package plex

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCollectionChildren(t *testing.T) {
	c := testClient()
	ratingKey := "236191"
	list, err := c.GetCollectionChildren(ratingKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, list)

	for _, item := range list {
		fmt.Printf("Title: %s\n", item.Title)
	}
}
