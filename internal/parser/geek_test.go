package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeekParser_Parse(t *testing.T) {
	p := New("geek", "https://api.nzbgeek.info/api?t=tvsearch&cat=5020,5030,5040,5045,5050")
	assert.NotNil(t, p, "instantiate parser")

	err := p.Parse()
	assert.NoError(t, err, "parse feed")

	items, err := p.Items()
	assert.NoError(t, err, "get items")
	assert.Len(t, items, 100, "items length")
}
