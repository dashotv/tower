package app

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

func (s *Server) Cron() error {
	if cfg.Cron {
		c := cron.New(cron.WithSeconds())

		// every 5 seconds DownloadsProcess
		// if _, err := c.AddFunc("*/5 * * * * *", s.DownloadsProcess); err != nil {
		// 	return errors.Wrap(err, "adding cron function")
		// }
		// if _, err := c.AddFunc("*/5 * * * * *", s.CausingErrors); err != nil {
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

		// every day at 3am
		if _, err := c.AddFunc("0 0 3 * * *", s.CleanPlexPins); err != nil {
			return errors.Wrap(err, "adding cron function: CleanPlexPins")
		}
		// every day at 3am
		if _, err := c.AddFunc("0 0 3 * * *", s.CleanJobs); err != nil {
			return errors.Wrap(err, "adding cron function: CleanJobs")
		}

		//TODO: clean up plex pins

		go func() {
			s.Log.Info("starting cron...")
			c.Start()
		}()
	}

	return nil
}

func (s *Server) CausingErrors() {
	minion.Add("causing errors", func(id int, log *zap.SugaredLogger) error {
		log.Info("causing error")
		return errors.Errorf("trying error")
	})
}
func (s *Server) DownloadsProcess() {
	minion.Add("downloads process", func(id int, log *zap.SugaredLogger) error {
		log.Info("processing downloads")
		return errors.Errorf("trying error")
	})
}

func (s *Server) ProcessFeeds() {
	s.Log.Info("processing feeds")
	db.ProcessFeeds()
}

func (s *Server) CleanPlexPins() {
	minion.Add("clean plex pins", func(id int, log *zap.SugaredLogger) error {
		list, err := db.Pin.Query().
			GreaterThan("created_at", time.Now().UTC().AddDate(0, 0, -1)).
			Run()
		if err != nil {
			return errors.Wrap(err, "querying pins")
		}

		for _, p := range list {
			err := db.Pin.Delete(p)
			if err != nil {
				return errors.Wrap(err, "deleting pin")
			}
		}

		return nil
	})
}

func (s *Server) CleanJobs() {
	minion.Add("clean jobs", func(id int, log *zap.SugaredLogger) error {
		list, err := db.MinionJob.Query().
			GreaterThan("created_at", time.Now().UTC().AddDate(0, 0, -1)).
			Run()
		if err != nil {
			return errors.Wrap(err, "querying jobs")
		}

		for _, j := range list {
			err := db.MinionJob.Delete(j)
			if err != nil {
				return errors.Wrap(err, "deleting job")
			}
		}

		return nil
	})
}

func (s *Server) PopularReleases() {
	minion.Add("popular releases", func(id int, log *zap.SugaredLogger) error {
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
					return errors.Wrap(err, fmt.Sprintf("popular releases %s %s", f, t))
				}

				cache.Set(fmt.Sprintf("releases_popular_%s_%s", f, t), results)
			}
		}

		diff := time.Since(start)
		log.Infof("PopularReleases: took %s", diff)

		return nil
	})
}
