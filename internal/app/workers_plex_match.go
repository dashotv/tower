package app

import (
	"context"

	"github.com/dashotv/fae"
	"github.com/dashotv/minion"
)

type PlexFiles struct {
	minion.WorkerDefaults[*PlexFiles]
}

func (j *PlexFiles) Kind() string { return "plex_files" }
func (j *PlexFiles) Work(ctx context.Context, job *minion.Job[*PlexFiles]) error {
	a := ContextApp(ctx)
	if a == nil {
		return fae.New("no app context")
	}

	defer TickTock("PlexFiles")()

	cache, err := buildPlexCache(ctx)
	if err != nil {
		return fae.Wrap(err, "build plex cache")
	}

	a.PlexFileCache = cache
	return nil
}
