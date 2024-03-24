package app

import (
	"context"
	"sync"
	"time"

	"github.com/sourcegraph/conc"
	"go.uber.org/ratelimit"

	"github.com/dashotv/minion"
)

var batchSize = 100
var scryRateLimit = 50 // per second

type Count struct {
	sync.Mutex
	i int
}

func (c *Count) Inc() {
	c.Lock()
	defer c.Unlock()
	c.i++
}

type UpdateIndexes struct {
	minion.WorkerDefaults[*UpdateIndexes]
}

func (j *UpdateIndexes) Kind() string { return "UpdateIndexes" }
func (j *UpdateIndexes) Timeout(job *minion.Job[*UpdateIndexes]) time.Duration {
	return 60 * time.Minute
}
func (j *UpdateIndexes) Work(ctx context.Context, job *minion.Job[*UpdateIndexes]) error {
	log := app.Log.Named("update_indexes")
	// ctx, cancel := context.WithCancel(ctx)
	// defer cancel()

	rl := ratelimit.New(scryRateLimit) // per second

	wg := conc.NewWaitGroup()
	wg.Go(func() {
		count := &Count{}
		total, err := app.DB.Series.Query().Limit(-1).Count()
		if err != nil {
			app.Workers.Log.Errorf("getting series count: %s", err)
			return
		}
		for i := 0; i < int(total); i += batchSize {
			series, err := app.DB.Series.Query().Desc("created_at").Limit(batchSize).Skip(i).Run()
			if err != nil {
				app.Workers.Log.Errorf("getting series: %s", err)
				return
			}
			for _, s := range series {
				rl.Take()
				if err := app.DB.Series.Update(s); err != nil {
					app.Workers.Log.Errorf("updating series: %s", err)
				}
				count.Inc()
			}
			log.Debugf("series: %d/%d", count.i, total)
		}
	})

	wg.Go(func() {
		count := &Count{}
		total, err := app.DB.Movie.Query().Limit(-1).Count()
		if err != nil {
			app.Workers.Log.Errorf("getting movies count: %s", err)
			return
		}
		for i := 0; i < int(total); i += batchSize {
			movies, err := app.DB.Movie.Query().Desc("created_at").Limit(batchSize).Skip(i).Run()
			if err != nil {
				app.Workers.Log.Errorf("getting movie: %s", err)
				return
			}
			for _, m := range movies {
				rl.Take()
				if err := app.DB.Movie.Update(m); err != nil {
					app.Workers.Log.Errorf("updating movie: %s", err)
				}
				count.Inc()
			}
			log.Debugf("series: %d/%d", count.i, total)
		}
	})

	wg.Go(func() {
		count := &Count{}
		total, err := app.DB.Release.Query().Limit(-1).Count()
		if err != nil {
			app.Workers.Log.Errorf("getting releases count: %s", err)
			return
		}
		for i := 0; i < int(total); i += batchSize {
			releases, err := app.DB.Release.Query().Desc("created_at").Limit(batchSize).Skip(i).Run()
			if err != nil {
				app.Workers.Log.Errorf("getting release: %s", err)
				return
			}
			for _, r := range releases {
				rl.Take()
				if err := app.DB.Release.Update(r); err != nil {
					app.Workers.Log.Errorf("updating release: %s", err)
				}
				count.Inc()
			}
			log.Debugf("series: %d/%d", count.i, total)
		}
	})

	wg.Wait()

	return nil
}
