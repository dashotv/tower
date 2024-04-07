package app

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFlame_Torrent(t *testing.T) {
	err := setupFlame(app)
	require.NoError(t, err)

	resp, err := app.Flame.Torrent("9f7cea6ea0d09ca0855c66026fe1c7ea2e274b0e")
	require.NoError(t, err)
	require.NotNil(t, resp)
	fmt.Printf("%+v\n", resp)
}

func TestFlame_RemoveTorrent(t *testing.T) {
	err := appSetup()
	require.NoError(t, err)
	err = setupFlame(app)
	require.NoError(t, err)

	err = app.Flame.RemoveTorrent("9f7cea6ea0d09ca0855c66026fe1c7ea2e274b0e")
	require.NoError(t, err)
}
