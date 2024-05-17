package app

import (
	"context"

	"github.com/dashotv/fae"
	"github.com/dashotv/minion"
)

type PlexWatched struct {
	minion.WorkerDefaults[*PlexWatched]
}

func (j *PlexWatched) Kind() string { return "plex_watched" }
func (j *PlexWatched) Work(ctx context.Context, job *minion.Job[*PlexWatched]) error {
	a := ContextApp(ctx)
	if a == nil {
		return fae.New("app not found")
	}

	_, err := a.Plex.GetAccountsUpdate()
	if err != nil {
		return fae.Wrap(err, "getting accounts")
	}

	list, err := a.Plex.GetHistoryRecent()
	if err != nil {
		return fae.Wrap(err, "getting history")
	}

	if err := a.plexHistoryWatched(list); err != nil {
		return fae.Wrap(err, "getting history media")
	}

	return nil
}

type PlexWatchedAll struct {
	minion.WorkerDefaults[*PlexWatchedAll]
}

func (j *PlexWatchedAll) Kind() string { return "plex_watched_all" }
func (j *PlexWatchedAll) Work(ctx context.Context, job *minion.Job[*PlexWatchedAll]) error {
	a := ContextApp(ctx)
	if a == nil {
		return fae.New("app not found")
	}

	_, err := a.Plex.GetAccountsUpdate()
	if err != nil {
		return fae.Wrap(err, "getting accounts")
	}

	total, err := a.Plex.GetHistoryTotal()
	if err != nil {
		return fae.Wrap(err, "getting history")
	}

	for i := int64(0); i < total; i += 200 {
		list, err := a.Plex.GetHistory(i, 100)
		if err != nil {
			return fae.Wrap(err, "getting history media")
		}

		if err := a.plexHistoryWatched(list); err != nil {
			return fae.Wrap(err, "getting history media")
		}
	}
	return nil
}
