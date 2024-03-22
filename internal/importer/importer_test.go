package importer

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

var testImporter *Importer

func init() {
	godotenv.Load("../../.env")
	opts := &Options{
		Language:  "eng",
		TvdbKey:   os.Getenv("TVDB_KEY"),
		TmdbToken: os.Getenv("TMDB_TOKEN"),
		FanartURL: os.Getenv("FANART_API_URL"),
		FanartKey: os.Getenv("FANART_API_KEY"),
	}

	i, err := New(opts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	testImporter = i
}

func Test_New(t *testing.T) {
	opts := &Options{
		Language:  "eng",
		TvdbKey:   os.Getenv("TVDB_KEY"),
		TmdbToken: os.Getenv("TMDB_TOKEN"),
		FanartURL: os.Getenv("FANART_URL"),
		FanartKey: os.Getenv("FANART_KEY"),
	}

	i, err := New(opts)
	require.NoError(t, err)
	require.NotNil(t, i)

	token := strings.SplitN(i.Tvdb.Token, ".", 2)
	fmt.Printf("token: %+v\n", token)
	fmt.Println("tvdb:", token[1])
}

func Test_Series(t *testing.T) {
	s, err := testImporter.Series(83602) // Lie to Me
	require.NoError(t, err)
	require.NotNil(t, s)
	require.Equal(t, "Lie to Me", s.Title)
}

func Test_SeriesEpisodes_Absolute(t *testing.T) {
	episodes, err := testImporter.SeriesEpisodes(392226, EpisodeOrderAbsolute) // Lie to Me
	require.NoError(t, err)
	require.NotNil(t, episodes)
	require.Greater(t, len(episodes), 0)
}

func Test_Series_Images(t *testing.T) {
	covers, backgrounds, err := testImporter.SeriesImages(83602) // Lie to Me
	require.NoError(t, err)
	require.NotNil(t, covers)
	require.NotNil(t, backgrounds)
	require.Greater(t, len(covers), 0)
	require.Greater(t, len(backgrounds), 0)

	spew.Dump(covers)
}
