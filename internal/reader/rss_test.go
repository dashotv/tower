package reader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRssReader_Parse(t *testing.T) {
	p := New("rss", "https://nyaa.si/?page=rss&c=1_2&f=0")
	assert.NotNil(t, p, "instantiate reader")

	err := p.Parse()
	assert.NoError(t, err, "parse feed")

	items, err := p.Items()
	assert.NoError(t, err, "get items")
	assert.Len(t, items, 75, "items length")
}

func TestRssReader_Process(t *testing.T) {
	p := New("rss", "https://nyaa.si/?page=rss&c=1_2&f=0")
	assert.NotNil(t, p, "instantiate reader")

	err := p.Process()
	assert.NoError(t, err, "parse feed")
}
