package app

import (
	"time"

	"github.com/madflojo/tasks"
	"go.uber.org/zap"
)

var workers *Workers

type Workers struct {
	log       *zap.SugaredLogger
	scheduler *tasks.Scheduler
}

func (w *Workers) Add(task *tasks.Task) (string, error) {
	return w.scheduler.Add(task)
}

// func setupPlex() (err error) {
// 	plexClient, err = plex.New(cfg.PlexToken, cfg.PlexURL, cfg.PlexMetadata, cfg.PlexTV)
// 	if err != nil {
// 		return errors.Wrap(err, "creating plex instance")
// 	}
// 	return nil
// }

func setupWorkers() error {
	workers = &Workers{
		log:       log.Named("workers"),
		scheduler: tasks.New(),
	}
	return nil
}

func schedulePinTask(p *Pin) (string, error) {
	return workers.Add(&tasks.Task{
		Interval: 60 * time.Second,
		RunOnce:  true,
		TaskFunc: func() error {
			workers.log.Info("polling pin: ", p.ID)

			// check if pin is done
			// if done, stop polling
			// if not done, continue polling
			return nil
		},
	})
}
