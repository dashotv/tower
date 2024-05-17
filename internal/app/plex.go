package app

import (
	"context"

	"github.com/dashotv/fae"
	"github.com/dashotv/tower/internal/plex"
)

func init() {
	initializers = append(initializers, setupPlex)
	// initializers = append(initializers, setupPlexFiles)
	// starters = append(starters, startPlexFiles)
}

func setupPlex(app *Application) error {
	p := plex.New(&plex.ClientOptions{
		URL:               app.Config.PlexServerURL,
		Token:             app.Config.PlexToken,
		Debug:             false,
		MachineIdentifier: app.Config.PlexMachineIdentifier,
		ClientIdentifier:  app.Config.PlexClientIdentifier,
		Product:           app.Config.PlexAppName,
		Device:            app.Config.PlexDevice,
		AppName:           app.Config.PlexAppName,
	})

	app.Plex = p
	return nil
}

func setupPlexFiles(a *Application) error {
	a.PlexFileCache = &plexFileCache{files: make(map[string]string)}
	return nil
}

func startPlexFiles(_ context.Context, a *Application) error {
	return a.Workers.Enqueue(&PlexMatch{})
}

func (a *Application) plexHistoryWatched(list []*plex.SessionMetadata) error {
	for _, session := range list {
		meta, err := a.Plex.GetMetadataByKey(session.RatingKey)
		if err != nil {
			return fae.Wrap(err, "getting metadata")
		}
		if len(meta) != 1 {
			return fae.Errorf("unexpected metadata count: %d", len(meta))
		}

		for _, me := range meta[0].Media {
			m, err := a.DB.MediumByPlexMedia(me)
			if err != nil {
				return fae.Wrap(err, "getting medium")
			}
			if m != nil {
				user, err := a.plexAccountTitle(session.AccountID)
				if err != nil {
					return fae.Wrap(err, "getting account title")
				}
				if err := a.DB.WatchMedium(m.ID, user, session.ViewedAt); err != nil {
					return fae.Wrap(err, "watch medium")
				}
				break
			}
		}
	}

	return nil
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

// func (a *Application) plexHistoryMedia() error {
// 	list, err := a.Plex.GetHistory()
// 	if err != nil {
// 		return fae.Wrap(err, "getting history")
// 	}
//
// 	for _, session := range list {
// 		meta, err := a.Plex.GetMetadataByKey(session.RatingKey)
// 		if err != nil {
// 			return fae.Wrap(err, "getting metadata")
// 		}
// 		if len(meta) != 1 {
// 			return fae.Errorf("unexpected metadata count: %d", len(meta))
// 		}
//
// 		for _, me := range meta[0].Media {
// 			m, err := a.DB.MediumByPlexMedia(me)
// 			if err != nil {
// 				return fae.Wrap(err, "getting medium")
// 			}
// 			if m != nil {
// 				user, err := a.plexAccountTitle(session.AccountID)
// 				if err != nil {
// 					return fae.Wrap(err, "getting account title")
// 				}
// 				if err := a.DB.WatchMedium(m.ID, user); err != nil {
// 					return fae.Wrap(err, "watch medium")
// 				}
// 				break
// 			}
// 		}
// 	}
//
// 	return nil
// }
