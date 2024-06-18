package app

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	runic "github.com/dashotv/runic/client"
)

func TestWant_Movie(t *testing.T) {
	want := NewWant(nil, nil)
	m := &Medium{ReleaseDate: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), SearchParams: &SearchParams{Resolution: 1080}}
	m.ID = primitive.NewObjectID()
	want.movies["title"] = m
	assert.Equal(t, m, want.releaseMovie(&runic.Release{Title: "title", Resolution: "1080", Year: 2020}))
	assert.Equal(t, m, want.releaseMovie(&runic.Release{Title: "TITLE", Resolution: "1080", Year: 2020}))
	assert.NotEqual(t, m, want.releaseMovie(&runic.Release{Title: "BLARG"}))
}

func TestWant_Series(t *testing.T) {
	m := &Medium{SeasonNumber: 1, EpisodeNumber: 1}
	want := NewWant(nil, nil)
	want.series_titles["title"] = "id"
	want.series_episodes["id"] = []*Medium{m}
	assert.Equal(t, m, want.releaseEpisode(&runic.Release{Title: "title", Resolution: "1080", Season: 1, Episode: 1}))
	assert.Equal(t, m, want.releaseEpisode(&runic.Release{Title: "TITLE", Resolution: "1080", Season: 1, Episode: 1}))
	assert.NotEqual(t, m, want.releaseEpisode(&runic.Release{Title: "title", Resolution: "1080", Season: 1, Episode: 2}))
}
