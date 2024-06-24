package app

import (
	"context"
	"sync"
	"time"

	"github.com/sourcegraph/conc"
	"go.uber.org/ratelimit"

	"github.com/dashotv/fae"
	"github.com/dashotv/minion"
)

var batchSize = 100
var scryRateLimit = 100 // per second

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
	a := ContextApp(ctx)
	log := app.Log.Named("update_indexes")
	// ctx, cancel := context.WithCancel(ctx)
	// defer cancel()

	rl := ratelimit.New(scryRateLimit) // per second

	wg := conc.NewWaitGroup()
	wg.Go(func() {
		err := a.DB.Series.Query().Desc("created_at").Each(100, func(s *Series) error {
			select {
			case <-ctx.Done():
				return fae.Errorf("cancelled")
			default:
				// proceed
			}

			rl.Take()
			if err := a.DB.Series.Update(s); err != nil {
				return fae.Wrapf(err, "updating series %s", s.ID.Hex())
			}
			return nil
		})
		if err != nil {
			log.Errorf("%s", err)
		}
	})

	wg.Go(func() {
		err := a.DB.Movie.Query().Desc("created_at").Each(100, func(s *Movie) error {
			select {
			case <-ctx.Done():
				return fae.Errorf("cancelled")
			default:
				// proceed
			}

			rl.Take()
			if err := a.DB.Movie.Update(s); err != nil {
				return fae.Wrapf(err, "updating movie %s", s.ID.Hex())
			}
			return nil
		})
		if err != nil {
			log.Errorf("%s", err)
		}
	})

	wg.Go(func() {
		err := a.DB.Episode.Query().Desc("created_at").Each(100, func(s *Episode) error {
			select {
			case <-ctx.Done():
				return fae.Errorf("cancelled")
			default:
				// proceed
			}

			rl.Take()
			if err := a.DB.Episode.Update(s); err != nil {
				return fae.Wrapf(err, "updating episode %s", s.ID.Hex())
			}
			return nil
		})
		if err != nil {
			log.Errorf("%s", err)
		}
	})

	wg.Wait()

	return nil
}
