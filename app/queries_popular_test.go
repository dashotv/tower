package app

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/dashotv/grimoire"
	"github.com/stretchr/testify/assert"
)

func TestRelease(t *testing.T) {
	url := os.Getenv("TEST_MONGODB_URL")
	if url == "" {
		t.Skip("TEST_MONGODB_URL not set")
	}

	g, err := grimoire.New[*Release](url, "torch_development", "torrents")
	assert.NoError(t, err, "grimoire.New")

	list, err := g.Query().Where("type", "tv").Limit(1).Run()
	assert.NoError(t, err, "query")
	assert.Len(t, list, 1, "query")
}

func TestReleasesPopular(t *testing.T) {
	url := os.Getenv("TEST_MONGODB_URL")
	if url == "" {
		t.Skip("TEST_MONGODB_URL not set")
	}

	limit := 25
	date := time.Now().AddDate(0, 0, -1)

	g, err := grimoire.New[*Release](url, "torch_development", "torrents")
	assert.NoError(t, err, "grimoire.New")

	rels, err := g.Query().Where("type", "tv").Desc("created_at").Limit(1).Run()
	assert.NoError(t, err, "query")

	if rels[0].CreatedAt.Before(date) {
		t.Skip("no releases yesterday")
	}

	start := time.Now()
	list, err := ReleasesPopularQuery(g.Collection, "movies", date, limit)
	end := time.Now()
	assert.NoError(t, err, "popular releases today")
	assert.Len(t, list, limit, "popular releases today")

	fmt.Printf("time: %s\n", end.Sub(start))
	for _, r := range list {
		// fmt.Printf("%35s %5s %s\n", r.PublishedAt, r.Type, r.Name)
		fmt.Printf("%-35.35s %d\n", r.Name, r.Count)
	}
}
