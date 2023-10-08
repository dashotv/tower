package app

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRelease(t *testing.T) {
	c := testConnector()
	if c == nil {
		t.Skip("No test connector")
		return
	}

	list, err := c.Release.Query().Where("type", "tv").Limit(1).Run()
	assert.NoError(t, err, "query")
	assert.Len(t, list, 1, "query")
}

func TestReleasesPopular(t *testing.T) {
	c := testConnector()
	if c == nil {
		t.Skip("No test connector")
		return
	}

	limit := 25
	date := time.Now().AddDate(0, 0, -1)

	rels, err := c.Release.Query().Where("type", "tv").Desc("created_at").Limit(1).Run()
	assert.NoError(t, err, "query")

	if rels[0].CreatedAt.Before(date) {
		t.Skip("no releases yesterday")
	}

	start := time.Now()
	list, err := ReleasesPopularQuery(c.Release.Collection, "movies", date, limit)
	end := time.Now()
	assert.NoError(t, err, "popular releases today")
	assert.Len(t, list, limit, "popular releases today")

	fmt.Printf("time: %s\n", end.Sub(start))
	for _, r := range list {
		// fmt.Printf("%35s %5s %s\n", r.PublishedAt, r.Type, r.Name)
		fmt.Printf("%-35.35s (%d) %d\n", r.Name, r.Year, r.Count)
	}
}
