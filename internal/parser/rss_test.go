package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRssParser_Parse(t *testing.T) {
	p := New("rss", "https://nyaa.si/?page=rss&c=1_2&f=0")
	assert.NotNil(t, p, "instantiate parser")

	err := p.Parse()
	assert.NoError(t, err, "parse feed")

	items, err := p.Items()
	assert.NoError(t, err, "get items")
	assert.Len(t, items, 75, "items length")
}

func TestRssParser_Process(t *testing.T) {
	p := New("rss", "https://nyaa.si/?page=rss&c=1_2&f=0")
	assert.NotNil(t, p, "instantiate parser")

	err := p.Process()
	assert.NoError(t, err, "parse feed")
}
