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

type historyResponse struct {
	AccountID int64
	Medium    *Medium
}
