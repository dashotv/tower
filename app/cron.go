package app

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
)

func (s *Server) Cron() error {
	if cfg.Cron {
		c := cron.New(cron.WithSeconds())

		// every 30 seconds DownloadsProcess
		// if _, err := c.AddFunc("*/5 * * * * *", s.DownloadsProcess); err != nil {
		// 	return errors.Wrap(err, "adding cron function")
		// }

		// every 5 minutes
		if _, err := c.AddFunc("0 */5 * * * *", s.PopularReleases); err != nil {
			return errors.Wrap(err, "adding cron function: PopularReleases")
		}
		// every 15 minutes
		// if _, err := c.AddFunc("0 */15 * * * *", s.ProcessFeeds); err != nil {
		// 	return errors.Wrap(err, "adding cron function: PopularReleases")
		// }

		go func() {
			s.Log.Info("starting cron...")
			c.Start()
		}()
	}

	return nil
}

func (s *Server) DownloadsProcess() {
	s.Log.Info("processing downloads")
}

func (s *Server) ProcessFeeds() {
	s.Log.Info("processing feeds")
	db.ProcessFeeds()
}

func (s *Server) PopularReleases() {
	limit := 25
	intervals := map[string]int{
		"daily":   1,
		"weekly":  7,
		"monthly": 30,
	}

	start := time.Now()
	for f, i := range intervals {
		for _, t := range releaseTypes {
			date := time.Now().AddDate(0, 0, -i)

			results, err := db.ReleasesPopularQuery(t, date, limit)
			if err != nil {
				s.Log.Error(errors.Wrap(err, fmt.Sprintf("popular releases %s %s", f, t)))
				return
			}

			cache.Set(fmt.Sprintf("releases_popular_%s_%s", f, t), results)
		}
	}

	diff := time.Since(start)
	s.Log.Infof("PopularReleases: took %s", diff)
}
