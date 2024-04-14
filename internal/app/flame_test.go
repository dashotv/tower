package app

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFlame_TorrentAdd(t *testing.T) {
	err := setupFlame(app)
	require.NoError(t, err)

	d := &Download{URL: "https://webtorrent.io/torrents/big-buck-bunny.torrent", Status: "loading"}

	resp, err := app.FlameTorrentAdd(d)
	require.NoError(t, err)
	require.NotNil(t, resp)
	fmt.Printf("%+v\n", resp)
}

func TestFlame_Torrent(t *testing.T) {
	err := setupFlame(app)
	require.NoError(t, err)

	resp, err := app.FlameTorrent("dd8255ecdc7ca55fb0bbf81323d87062db1f6d1c")
	require.NoError(t, err)
	require.NotNil(t, resp)
	fmt.Printf("%+v\n", resp)
}

func TestFlame_TorrentRemove(t *testing.T) {
	err := setupFlame(app)
	require.NoError(t, err)

	err = app.FlameTorrentRemove("dd8255ecdc7ca55fb0bbf81323d87062db1f6d1c")
	require.NoError(t, err)
}
