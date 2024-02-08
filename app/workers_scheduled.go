package app

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

func init() {
	initializers = append(initializers, func(app *Application) error {
		app.Workers.ScheduleFunc("0 */15 * * * *", "PopularReleases", PopularReleases)
		return nil
	})
}

func PopularReleases() error {
	limit := 25
	intervals := map[string]int{
		"daily":   1,
		"weekly":  7,
		"monthly": 30,
	}

	for f, i := range intervals {
		for _, t := range releaseTypes {
			date := time.Now().AddDate(0, 0, -i)

			results, err := app.DB.ReleasesPopularQuery(t, date, limit)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("popular releases %s %s", f, t))
			}

			app.Cache.Set(fmt.Sprintf("releases_popular_%s_%s", f, t), results)
		}
	}

	return nil
}
