package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWant_Movie(t *testing.T) {
	want := NewWant(nil, nil)
	want.movies["title"] = "id"
	assert.Equal(t, "id", want.Movie("title"))
	assert.Equal(t, "id", want.Movie("TITLE"))
	assert.NotEqual(t, "id", want.Movie("Blarg"))
}

func TestWant_Series(t *testing.T) {
	want := NewWant(nil, nil)
	want.series_titles["title"] = "id"
	want.series_episodes["id"] = []*Episode{
		{
			SeasonNumber:  1,
			EpisodeNumber: 1,
		},
	}
	assert.Equal(t, "000000000000000000000000", want.Episode("title", 1, 1))
	assert.Equal(t, "000000000000000000000000", want.Episode("TITLE", 1, 1))
	assert.NotEqual(t, "000000000000000000000000", want.Episode("title", 1, 2))
}
