package app

import (
	"fmt"
	"time"

	"github.com/dashotv/fae"
)

func init() {
	initializers = append(initializers, func(app *Application) error {
		app.Workers.ScheduleFunc("* * * * * *", "plex_session_updates", PlexSessionUpdates)
		app.Workers.ScheduleFunc("0 */15 * * * *", "PopularReleases", PopularReleases)
		return nil
	})
}

func PlexSessionUpdates() error {
	sessions, err := app.Plex.GetSessions()
	if err != nil {
		app.Log.Named("PlexSessionUpdates").Error(err)
		return err
	}

	return app.Events.Send("tower.plex_sessions", &EventPlexSessions{sessions})
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
				return fae.Wrap(err, fmt.Sprintf("popular releases %s %s", f, t))
			}

			app.Cache.Set(fmt.Sprintf("releases_popular_%s_%s", f, t), results)
		}
	}

	return nil
}
