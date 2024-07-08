package importer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeriesUpdated(t *testing.T) {
	t.Log(testImporter.Tvdb.Token)
	list, err := testImporter.SeriesUpdated(1720371986)
	assert.NoError(t, err)
	assert.NotEmpty(t, list)
	t.Logf("updated: %v", list)
}
