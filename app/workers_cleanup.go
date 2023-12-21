package app

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/dashotv/minion"
)

// CleanPlexPins removes old pins
type CleanPlexPins struct {
	minion.WorkerDefaults[*CleanPlexPins]
}

func (j *CleanPlexPins) Kind() string { return "CleanPlexPins" }
func (j *CleanPlexPins) Work(ctx context.Context, job *minion.Job[*CleanPlexPins]) error {
	list, err := app.DB.Pin.Query().
		GreaterThan("created_at", time.Now().UTC().AddDate(0, 0, -1)).
		Run()
	if err != nil {
		return errors.Wrap(err, "querying pins")
	}

	for _, p := range list {
		err := app.DB.Pin.Delete(p)
		if err != nil {
			return errors.Wrap(err, "deleting pin")
		}
	}

	return nil
}

type CleanupLogs struct {
	minion.WorkerDefaults[*CleanupLogs]
}

func (j *CleanupLogs) Kind() string { return "cleanup_logs" }
func (j *CleanupLogs) Work(ctx context.Context, job *minion.Job[*CleanupLogs]) error {
	return nil
}

type CleanupJobs struct {
	minion.WorkerDefaults[*CleanupJobs]
}

func (j *CleanupJobs) Kind() string { return "cleanup_jobs" }
func (j *CleanupJobs) Work(ctx context.Context, job *minion.Job[*CleanupJobs]) error {
	return nil
}
