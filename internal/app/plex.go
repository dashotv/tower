package app

import (
	"context"
	"fmt"

	"github.com/dashotv/fae"
	"github.com/dashotv/tower/internal/plex"
)

func init() {
	initializers = append(initializers, setupPlex)
	initializers = append(initializers, setupPlexFiles)
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

type plexFileCache struct {
	files map[string]string
}

func buildPlexCache(ctx context.Context) (*plexFileCache, error) {
	a := ContextApp(ctx)
	if a == nil {
		return nil, fae.New("no app context")
	}

	cache := &plexFileCache{files: make(map[string]string)}

	libs, err := a.Plex.GetLibraries()
	if err != nil {
		return nil, fae.Wrap(err, "get libraries")
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

		_, total, err := a.Plex.GetLibrarySection(lib.Key, "all", t, 0, 1)
		if err != nil {
			return nil, fae.Wrapf(err, "get library section: %s", lib.Key)
		}

		for i := int64(0); i < total; i += 50 {
			list, _, err := a.Plex.GetLibrarySection(lib.Key, "all", t, int(i), 50)
			if err != nil {
				return nil, fae.Wrap(err, "get library section")
			}
			for _, item := range list {
				if len(item.Media) == 0 {
					continue
				}

				for _, media := range item.Media {
					if len(media.Part) == 0 {
						continue
					}

					for _, part := range media.Part {
						if part.File == "" {
							continue
						}

						if _, ok := cache.files[part.File]; !ok {
							cache.files[part.File] = fmt.Sprintf("%s", item.RatingKey)
						}
					}
				}
			}
		}
	}

	return cache, nil
}
