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
	a := ContextApp(ctx)
	if a == nil {
		return fae.New("no app context")
	}

	id := job.Args.Id

	c, err := a.DB.Collection.Get(id, &Collection{})
	if err != nil {
		return err
	}

	// TODO: sync first, so we don't blow up the paylist?

	if len(c.Media) == 0 {
		return nil
	}
	if c.RatingKey == "" {
		resp, err := a.Plex.CreateCollection(c.Name, c.Library, c.Media[0].RatingKey)
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
	add, remove, err := a.Plex.UpdateCollection(c.Library, c.RatingKey, keys)
	if err != nil {
		return err
	}

	if len(add) > 0 {
		for _, k := range add {
			if !lo.Contains(keys, k) {
				a.Log.Debugf("add %s", k)
			}
		}
	}

	if len(remove) > 0 {
		for _, k := range remove {
			if lo.Contains(keys, k) {
				a.Log.Debugf("remove %s", k)
			}
		}
	}

	c.SyncedAt = time.Now()

	return a.DB.Collection.Save(c)
}
