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
		if _, err := c.AddFunc("0 0 1 * * *", s.PopularReleasesToday); err != nil {
			return errors.Wrap(err, "adding cron function: PopularReleasesToday")
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

func (s *Server) PopularReleasesToday() {
	App().Log.Info("processing popular releases")

	date := time.Now().AddDate(0, 0, -1)
	limit := 25

	results, err := App().DB.ReleasesPopular(date, limit)
	if err != nil {
		s.Log.Error(errors.Wrap(err, "popular releases today"))
		return
	}

	cache := App().Cache
	cache.Set("popular_releases_today", results)
}
