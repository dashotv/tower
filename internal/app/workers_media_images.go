package app

import (
	"context"

	"github.com/dashotv/fae"
	"github.com/dashotv/minion"
)

type MediaImages struct {
	minion.WorkerDefaults[*MediaImages]
}

func (j *MediaImages) Kind() string { return "media_images" }
func (j *MediaImages) Work(ctx context.Context, job *minion.Job[*MediaImages]) error {
	a := ContextApp(ctx)
	// l := a.Workers.Log.Named("media_images")
	err := a.DB.Series.Query().Where("paths.type", "").Each(100, func(m *Series) error {
		// l.Infof("media_images: %s", m.ID)
		return a.Workers.Enqueue(&SeriesUpdate{ID: m.ID.Hex(), SkipImages: false, Title: m.Title})
	})
	if err != nil {
		return fae.Wrap(err, "querying media")
	}
	err = a.DB.Movie.Query().Where("paths.type", "").Each(100, func(m *Movie) error {
		// l.Infof("media_images: %s", m.ID)
		return a.Workers.Enqueue(&MovieUpdate{ID: m.ID.Hex(), SkipImages: false, Title: m.Title})
	})
	if err != nil {
		return fae.Wrap(err, "querying media")
	}
	return nil
}
