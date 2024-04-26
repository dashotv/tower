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

func TestClient_GetLibrarySection(t *testing.T) {
	p := testClient()
	require.NotNil(t, p)

	list, total, err := p.GetLibrarySection("2", "all", "4", 0, 10)
	require.NoError(t, err)
	require.NotNil(t, list)
	require.NotEmpty(t, list)
	require.Len(t, list, 10)
	require.Greater(t, total, int64(0))

	list, total, err = p.GetLibrarySection("2", "all", "4", 0, 1)
	require.NoError(t, err)
	require.Len(t, list, 1)
	require.Greater(t, total, int64(0))
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
