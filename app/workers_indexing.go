package app

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/dashotv/minion"
)

type UpdateIndexes struct{}

func (j *UpdateIndexes) Kind() string { return "UpdateIndexes" }
func (j *UpdateIndexes) Work(ctx context.Context, job *minion.Job[*UpdateIndexes]) error {
	series, err := db.SeriesAll()
	if err != nil {
		return errors.Wrap(err, "getting series")
	}
	for _, s := range series {
		<-time.After(1 * time.Second)
		if err := workers.Enqueue(&TvdbUpdateSeries{s.ID.Hex(), false, false, false}); err != nil {
			return errors.Wrap(err, "enqueueing series")
		}
	}

	movies, err := db.MoviesAll()
	if err != nil {
		return errors.Wrap(err, "getting movies")
	}
	log.Infof("updating %d movies", len(movies))
	for _, m := range movies {
		<-time.After(1 * time.Second)
		if err := workers.Enqueue(&TmdbUpdateMovie{m.ID.Hex(), false}); err != nil {
			return errors.Wrap(err, "enqueueing movie")
		}
	}

	return nil
}
