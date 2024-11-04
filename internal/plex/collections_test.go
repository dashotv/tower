package plex

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCollection(t *testing.T) {
	c := testClient()
	ratingKey := "246975"
	list, err := c.GetCollection(ratingKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, list)
}

func TestGetCollectionChildren(t *testing.T) {
	c := testClient()
	ratingKey := "246975"
	list, err := c.GetCollectionChildren(ratingKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, list)

	for _, item := range list {
		fmt.Printf("Title: %s\n", item.Title)
	}
}
