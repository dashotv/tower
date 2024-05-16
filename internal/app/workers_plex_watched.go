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

	resp, err := a.plexHistoryMedia()
	if err != nil {
		return fae.Wrap(err, "getting history media")
	}

	for _, h := range resp {
		a.Log.Named("plex_watched").Infof("watched: %d %s", h.AccountID, h.Medium.Title)
		user, err := a.plexAccountTitle(h.AccountID)
		if err != nil {
			return fae.Wrap(err, "getting account title")
		}
		if err := a.DB.WatchMedium(h.Medium.ID, user); err != nil {
			return fae.Wrap(err, "watch medium")
		}
	}

	return nil
}

type historyResponse struct {
	AccountID int64
	Medium    *Medium
}

func (a *Application) plexHistoryMedia() ([]*historyResponse, error) {
	media := []*historyResponse{}

	list, err := a.Plex.GetHistory()
	if err != nil {
		return nil, fae.Wrap(err, "getting history")
	}

	for _, session := range list {
		meta, err := a.Plex.GetMetadataByKey(session.RatingKey)
		if err != nil {
			return nil, fae.Wrap(err, "getting metadata")
		}
		if len(meta) != 1 {
			return nil, fae.Errorf("unexpected metadata count: %d", len(meta))
		}

		for _, me := range meta[0].Media {
			m, err := a.DB.MediumByPlexMedia(me)
			if err != nil {
				return nil, fae.Wrap(err, "getting medium")
			}
			if m != nil {
				media = append(media, &historyResponse{AccountID: session.AccountID, Medium: m})
				break
			}
		}
	}

	return media, nil
}

func (a *Application) plexAccountTitle(id int64) (string, error) {
	if id == 1 {
		return a.Config.PlexUsername, nil
	}
	u, err := a.Plex.GetAccount(id)
	if err != nil {
		return "", fae.Wrap(err, "getting account")
	}

	return u.Title, nil
}
