package app

import (
	"context"
	"time"

	"github.com/samber/lo"

	"github.com/dashotv/fae"
	"github.com/dashotv/minion"
)

type PlexCollectionUpdate struct {
	minion.WorkerDefaults[*PlexCollectionUpdate]
	Id string `bson:"id" json:"id"`
}

func (j *PlexCollectionUpdate) Kind() string { return "plex_collection_update" }
func (j *PlexCollectionUpdate) Work(ctx context.Context, job *minion.Job[*PlexCollectionUpdate]) error {
	id := job.Args.Id

	c, err := app.DB.Collection.Get(id, &Collection{})
	if err != nil {
		return err
	}

	if len(c.Media) == 0 {
		return nil
	}
	if c.RatingKey == "" {
		resp, err := app.Plex.CreateCollection(c.Name, c.Library, c.Media[0].RatingKey)
		if err != nil {
			return err
		}
		if len(resp.MediaContainer.Directory) == 0 {
			return fae.New("api response did not contain a directory")
		}

		c.RatingKey = resp.MediaContainer.Directory[0].RatingKey
	}

	keys := lo.Map(c.Media, func(m *CollectionMedia, i int) string {
		return m.RatingKey
	})
	if err := app.Plex.UpdateCollection(c.Library, c.RatingKey, keys); err != nil {
		return err
	}

	c.SyncedAt = time.Now()

	return app.DB.Collection.Save(c)
}
