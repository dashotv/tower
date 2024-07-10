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

	if err := a.PlexFileCache.build(ctx); err != nil {
		return fae.Wrap(err, "build plex cache")
	}

	return nil
}

type PlexFilesPartial struct {
	minion.WorkerDefaults[*PlexFilesPartial]
	Title   string `bson:"title" json:"title"`
	Section string `bson:"section" json:"section"`
	Libtype string `bson:"libtype" json:"libtype"`
}

func (j *PlexFilesPartial) Kind() string { return "plex_files_partial" }
func (j *PlexFilesPartial) Work(ctx context.Context, job *minion.Job[*PlexFilesPartial]) error {
	a := ContextApp(ctx)
	if a == nil {
		return fae.New("no app context")
	}

	//l := a.Workers.Log.Named("plex_files_partial")
	if err := a.PlexFileCache.update(ctx, job.Args.Title, job.Args.Section, job.Args.Libtype); err != nil {
		return fae.Wrap(err, "update plex cache")
	}

	return nil
}
