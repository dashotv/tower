package importer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSeriesUpdated(t *testing.T) {
	t.Log(testImporter.Tvdb.Token)
	list, err := testImporter.SeriesUpdated(time.Now().Add(-24 * time.Hour).Unix())
	assert.NoError(t, err)
	assert.NotEmpty(t, list)
	t.Logf("updated: %v", list)
}
