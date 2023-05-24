package app

import (
	"github.com/robfig/cron/v3"
)

func (s *Server) Cron() error {
	if ConfigInstance().Cron {
		c := cron.New(cron.WithSeconds())

		//	// every 30 seconds DownloadsProcess
		//	if _, err := c.AddFunc("*/30 * * * * *", s.DownloadsProcess); err != nil {
		//		return errors.Wrap(err, "adding cron function")
		//	}

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
