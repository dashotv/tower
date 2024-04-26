package app

import (
	"context"
	"fmt"

	"github.com/dashotv/fae"
	"github.com/dashotv/minion"
)

type plexFileCache struct {
	files map[string]string
}

type PlexMatch struct {
	minion.WorkerDefaults[*PlexMatch]
}

func (j *PlexMatch) Kind() string { return "plex_match" }
func (j *PlexMatch) Work(ctx context.Context, job *minion.Job[*PlexMatch]) error {
	a := ContextApp(ctx)
	if a == nil {
		return fae.New("no app context")
	}
	// l := a.Log.Named("PlexMatch")

	defer TickTock("PlexMatch")()
	cache := &plexFileCache{files: make(map[string]string)}

	libs, err := app.Plex.GetLibraries()
	if err != nil {
		return fae.Wrap(err, "get libraries")
	}
	for _, lib := range libs {
		t := ""
		if lib.Type == "show" {
			t = "4"
		} else if lib.Type == "movie" {
			t = "1"
		} else {
			continue
		}

		_, total, err := app.Plex.GetLibrarySection(lib.Key, "all", t, 0, 1)
		if err != nil {
			return fae.Wrapf(err, "get library section: %s", lib.Key)
		}

		for i := int64(0); i < total; i += 50 {
			list, _, err := app.Plex.GetLibrarySection(lib.Key, "all", t, int(i), 50)
			if err != nil {
				return fae.Wrap(err, "get library section")
			}
			for _, item := range list {
				if len(item.Media) > 0 {
					for _, media := range item.Media {
						if len(media.Part) > 0 {
							for _, part := range media.Part {
								if part.File != "" {
									if _, ok := cache.files[part.File]; !ok {
										cache.files[part.File] = fmt.Sprintf("%s", item.RatingKey)
									}
								}
							}
						}
					}
				}
			}
		}
	}

	a.PlexFileCache = cache
	return nil
}
