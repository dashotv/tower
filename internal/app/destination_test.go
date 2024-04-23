package app

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDestinator_Destination(t *testing.T) {
	err := startDestination(context.TODO(), app)
	require.NoError(t, err)

	cases := []struct {
		downloadID  string
		destination string
	}{
		{"6626fa42e7d0b66eef8a57e7", "/mnt/media/anime/unnamed memory/unnamed memory - 01x003"},
		{"66245728a067cb89f8403a57", "/mnt/media/tv/last week tonight with john oliver/last week tonight with john oliver - 11x09"},
		{"6623736da067cb89f84022c7", "/mnt/media/movies/red rocket/red rocket"},
		{"66275b4fafa4a0dc7092d327", "/mnt/media/donghua/blader soul/blader soul - 01x015"},
	}

	destinator := app.Destinator

	for _, c := range cases {
		t.Run(c.downloadID, func(tt *testing.T) {
			d := &Download{}
			err := app.DB.Download.Find(c.downloadID, d)
			require.NoError(tt, err)
			require.NotNil(tt, d)
			app.DB.processDownloads([]*Download{d})
			require.NotNil(tt, d.Medium)

			dest, err := destinator.Destination(d.Kind, d.Medium)
			require.NoError(tt, err)
			require.Equal(tt, c.destination, dest)
		})
	}
}
