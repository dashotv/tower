package app

import (
	"context"

	"github.com/pkg/errors"

	"github.com/dashotv/minion"
)

type TmdbUpdateAll struct {
	minion.WorkerDefaults[*TmdbUpdateAll]
}

func (j *TmdbUpdateAll) Kind() string { return "TmdbUpdateAll" }
func (j *TmdbUpdateAll) Work(ctx context.Context, job *minion.Job[*TmdbUpdateAll]) error {
	movies, err := app.DB.Movie.Query().Limit(-1).Run()
	if err != nil {
		return errors.Wrap(err, "querying movies")
	}

	for _, m := range movies {
		app.Log.Infof("updating movie: %s", m.Title)
		app.Workers.Enqueue(&TmdbUpdateMovie{ID: m.ID.Hex(), JustMedia: true})
	}

	return nil
}
