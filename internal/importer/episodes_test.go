package importer

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_LoadEpisodeMap(t *testing.T) {
	resp, err := testImporter.loadEpisodesMap(337284, EpisodeOrderDefault)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp)

	for id, e := range resp {
		fmt.Printf("%d: %s\n", id, e.Title)
	}
}
