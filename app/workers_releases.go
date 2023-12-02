package app

import (
	"context"
	"fmt"
	"time"

	"github.com/dashotv/minion"
	"github.com/pkg/errors"
)

// PopularReleases updates the popular releases cache
type PopularReleases struct{}

func (j *PopularReleases) Kind() string { return "PopularReleases" }
func (j *PopularReleases) Work(ctx context.Context, job *minion.Job[*PopularReleases]) error {
	log.Named("popular_releases").Debug("popular releases")
	limit := 25
	intervals := map[string]int{
		"daily":   1,
		"weekly":  7,
		"monthly": 30,
	}

	for f, i := range intervals {
		for _, t := range releaseTypes {
			date := time.Now().AddDate(0, 0, -i)

			results, err := db.ReleasesPopularQuery(t, date, limit)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("popular releases %s %s", f, t))
			}

			cache.Set(fmt.Sprintf("releases_popular_%s_%s", f, t), results)
		}
	}

	return nil
}
