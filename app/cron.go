package app

import (
	"time"

	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
)

func (s *Server) Cron() error {
	if ConfigInstance().Cron {
		c := cron.New(cron.WithSeconds())

		//	// every 30 seconds DownloadsProcess
		//	if _, err := c.AddFunc("*/30 * * * * *", s.DownloadsProcess); err != nil {
		//		return errors.Wrap(err, "adding cron function")
		//	}

		// every day at 1am PopularReleasesToday
		if _, err := c.AddFunc("0 0 1 * * *", s.PopularReleasesDaily); err != nil {
			return errors.Wrap(err, "adding cron function: PopularReleasesDaily")
		}

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

func (s *Server) PopularReleasesDaily() {
	App().Log.Info("processing popular releases")

	types := []string{"tv", "anime", "movie"}
	for _, t := range types {
		date := time.Now().AddDate(0, 0, -1)
		limit := 25

		results, err := App().DB.ReleasesPopular("tv", date, limit)
		if err != nil {
			s.Log.Error(errors.Wrap(err, "popular releases today"))
			return
		}

		cache := App().Cache
		cache.Set("releases_popular_daily_"+t, results)
	}
}
