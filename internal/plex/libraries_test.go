package plex

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_LibraryByPath(t *testing.T) {
	p := testClient()
	require.NotNil(t, p)

	lib, err := p.LibraryByPath("/mnt/media/anime")
	require.NoError(t, err)
	require.NotNil(t, lib)
	require.Equal(t, "Anime", lib.Title)
}

func TestClient_RefreshLibraryPath(t *testing.T) {
	p := testClient()
	require.NotNil(t, p)

	err := p.RefreshLibraryPath("/mnt/media/anime/solo leveling")
	require.NoError(t, err)
}

func TestClient_GetLibraries(t *testing.T) {
	p := testClient()
	require.NotNil(t, p)

	libs, err := p.GetLibraries()
	require.NoError(t, err)
	require.NotNil(t, libs)
	require.NotEmpty(t, libs)
}
