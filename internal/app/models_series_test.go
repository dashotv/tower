package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnector_SeriesAllUnwatched(t *testing.T) {
	c := testConnector()
	if c == nil {
		t.Skip("No test connector")
		return
	}

	series := &Series{}
	err := c.Series.Find("644d88003359bb748dd63096", series)
	assert.NoError(t, err, "Find Series")

	got, err := c.SeriesUserUnwatched(series)
	assert.NoError(t, err, "unwatched")
	assert.Greater(t, got, 0, "unwatched")
}

func TestConnector_SeriesBySearch(t *testing.T) {
	db := testConnector()
	if db == nil {
		t.Skip("No test connector")
		return
	}

	cases := []struct{ title, expected string }{
		{"the first order", "The First Order"},
		{"yishi-zhi-zun", "Ancient Lords"},
	}

	for _, c := range cases {
		series, err := db.SeriesBySearch(c.title)
		require.NoError(t, err)
		require.NotNil(t, series)
		require.Equal(t, c.expected, series.Title)
	}
}

func TestEvents_SeriesEpisodeByRunic(t *testing.T) {
	series, err := app.DB.SeriesBySearch("the age of cosmos exploration")
	require.NoError(t, err)
	require.NotNil(t, series)

	e, err := app.DB.SeriesEpisodeBy(series, 1, 9)
	require.NoError(t, err)
	require.NotNil(t, e)
}
