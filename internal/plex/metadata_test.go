package plex

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func TestGetMetadataByKey(t *testing.T) {
	resp, err := testClient().GetMetadataByKey("233509")
	require.NoError(t, err)
	spew.Dump(resp)
}
func TestPutMetadataPrefs(t *testing.T) {
	key := "234212" // Everlasting God of Sword
	err := testClient().PutMetadataPrefs(key, map[string]string{"showOrdering": "tvdbAbsolute"})
	require.NoError(t, err)
}
