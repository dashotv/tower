package app

import (
	"context"

	"github.com/dashotv/fae"
	"github.com/dashotv/minion"
)

type MovieDelete struct {
	minion.WorkerDefaults[*MovieDelete]
	ID string `bson:"id" json:"id"`
}

func (j *MovieDelete) Kind() string { return "movie_delete" }
func (j *MovieDelete) Work(ctx context.Context, job *minion.Job[*MovieDelete]) error {
	id := job.Args.ID

	movie := &Movie{}
	if err := app.DB.Movie.Find(id, movie); err != nil {
		return fae.Wrap(err, "finding movie")
	}

	if err := app.DB.Movie.Delete(movie); err != nil {
		return fae.Wrap(err, "deleting series")
	}

	if err := app.Workers.Enqueue(&PathDeleteAll{MediumID: movie.ID.Hex()}); err != nil {
		return fae.Wrap(err, "enqueueing paths")
	}
	return nil
}
