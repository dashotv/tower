package app

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/dashotv/minion"
)

type UpdateIndexes struct {
	minion.WorkerDefaults[*UpdateIndexes]
}

func (j *UpdateIndexes) Kind() string { return "UpdateIndexes" }
func (j *UpdateIndexes) Work(ctx context.Context, job *minion.Job[*UpdateIndexes]) error {
	series, err := app.DB.SeriesAll()
	if err != nil {
		return errors.Wrap(err, "getting series")
	}
	for _, s := range series {
		<-time.After(1 * time.Second)
		if err := app.Workers.Enqueue(&TvdbUpdateSeries{ID: s.ID.Hex(), JustMedia: true}); err != nil {
			return errors.Wrap(err, "enqueueing series")
		}
	}

	movies, err := app.DB.MoviesAll()
	if err != nil {
		return errors.Wrap(err, "getting movies")
	}
	for _, m := range movies {
		<-time.After(1 * time.Second)
		if err := app.Workers.Enqueue(&TmdbUpdateMovie{ID: m.ID.Hex(), JustMedia: true}); err != nil {
			return errors.Wrap(err, "enqueueing movie")
		}
	}

	return nil
}
