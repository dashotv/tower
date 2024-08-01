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
		{"6626fa42e7d0b66eef8a57e7", "/mnt/media/anime/unnamed memory/unnamed memory - 01x03 #003 - what the forest dreams of"},
		{"66245728a067cb89f8403a57", "/mnt/media/tv/last week tonight with john oliver/last week tonight with john oliver - 11x09 - april 21, 2024: ufos"},
		{"6623736da067cb89f84022c7", "/mnt/media/movies/red rocket/red rocket"},
		{"66275b4fafa4a0dc7092d327", "/mnt/media/donghua/blader soul/blader soul - 01x015"},
		{"664b7d74a65b5d1db94edb21", "/mnt/media/donghua/core sense/core sense - 01x002 - 记者暗访"},
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

func TestDestinator_File(t *testing.T) {
	err := startDestination(context.TODO(), app)
	require.NoError(t, err)

	cases := []struct {
		fileID      string
		destination string
	}{
		{"667504a777694a06672e05d7", "/mnt/media/donghua/mysterious treasures/mysterious treasures s01e01-06.mp4"},
	}

	destinator := app.Destinator

	for _, c := range cases {
		t.Run(c.fileID, func(tt *testing.T) {
			f := &File{}
			err := app.DB.File.Find(c.fileID, f)
			require.NoError(tt, err)
			require.NotNil(tt, f)

			dest, err := destinator.File(f)
			require.NoError(tt, err)
			require.Equal(tt, c.destination, dest)
		})
	}
}
