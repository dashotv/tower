package app

import (
	"context"
	"time"

	"github.com/sourcegraph/conc"

	"github.com/dashotv/minion"
)

type UpdateIndexes struct {
	minion.WorkerDefaults[*UpdateIndexes]
}

func (j *UpdateIndexes) Kind() string { return "UpdateIndexes" }
func (j *UpdateIndexes) Timeout(job *minion.Job[*UpdateIndexes]) time.Duration {
	return 60 * time.Minute
}
func (j *UpdateIndexes) Work(ctx context.Context, job *minion.Job[*UpdateIndexes]) error {
	wg := conc.NewWaitGroup()

	wg.Go(func() {
		total, err := app.DB.Series.Query().Limit(-1).Count()
		if err != nil {
			app.Workers.Log.Errorf("getting series count: %s", err)
			return
		}
		for i := 0; i < int(total); i += 100 {
			series, err := app.DB.Series.Query().Desc("created_at").Limit(100).Skip(i).Run()
			if err != nil {
				app.Workers.Log.Errorf("getting series: %s", err)
				return
			}
			for _, s := range series {
				if err := app.DB.Series.Update(s); err != nil {
					app.Workers.Log.Errorf("updating series: %s", err)
				}
			}
		}
	})

	wg.Go(func() {
		total, err := app.DB.Movie.Query().Limit(-1).Count()
		if err != nil {
			app.Workers.Log.Errorf("getting movies count: %s", err)
			return
		}
		for i := 0; i < int(total); i += 100 {
			movies, err := app.DB.Movie.Query().Desc("created_at").Limit(100).Skip(i).Run()
			if err != nil {
				app.Workers.Log.Errorf("getting movie: %s", err)
				return
			}
			for _, m := range movies {
				if err := app.DB.Movie.Update(m); err != nil {
					app.Workers.Log.Errorf("updating movie: %s", err)
				}
			}
		}
	})

	wg.Go(func() {
		total, err := app.DB.Release.Query().Limit(-1).Count()
		if err != nil {
			app.Workers.Log.Errorf("getting releases count: %s", err)
			return
		}
		for i := 0; i < int(total); i += 100 {
			releases, err := app.DB.Release.Query().Desc("created_at").Limit(100).Skip(i).Run()
			if err != nil {
				app.Workers.Log.Errorf("getting release: %s", err)
				return
			}
			for _, r := range releases {
				if err := app.DB.Release.Update(r); err != nil {
					app.Workers.Log.Errorf("updating release: %s", err)
				}
			}
		}
	})

	wg.Wait()

	return nil
}
