package app

import (
	"context"
	"net/url"
	"time"

	"github.com/dashotv/fae"
	"github.com/dashotv/tower/internal/plex"
)

var buildPlexCacheMutex = &CtxMutex{ch: make(chan struct{}, 1)}

func init() {
	initializers = append(initializers, setupPlex)
	starters = append(starters, startPlexFiles)
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

func startPlexFiles(ctx context.Context, a *Application) error {
	a.PlexFileCache = &plexFileCache{}
	return a.Workers.Enqueue(&PlexFiles{})
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
	libs    []*plex.Library
	files   map[string]*plex.LeavesMetadata // absolute path -> plex file metadata
	parents map[string]string               // title -> parent ID
}

func plexLibType(t string) string {
	switch t {
	case "show":
		return "4"
	case "movie":
		return "1"
	}
	return ""
}

func (c *plexFileCache) build(ctx context.Context) error {
	muctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if !buildPlexCacheMutex.Lock(muctx) {
		app.Log.Named("buildPlexCache").Warn("failed to lock mutex")
		return nil
	}
	defer buildPlexCacheMutex.Unlock()

	a := ContextApp(ctx)
	if a == nil {
		return fae.New("no app context")
	}

	libs, err := a.Plex.GetLibraries()
	if err != nil {
		return fae.Wrap(err, "get libraries")
	}
	c.libs = libs

	if err := c.buildFiles(ctx); err != nil {
		return fae.Wrap(err, "build files")
	}
	if err := c.buildFolders(ctx); err != nil {
		return fae.Wrap(err, "build folders")
	}
	return nil
}

// func (c *plexFileCache) buildPlexCache(ctx context.Context) (*plexFileCache, error) {
//
// 	a := ContextApp(ctx)
// 	if a == nil {
// 		return fae.New("no app context")
// 	}
//
// 	for _, lib := range c.libs {
// 		t := plexLibType(lib.Type)
// 		if t == "" {
// 			continue
// 		}
//
// 		folders, total, err := a.Plex.GetLibrarySectionFolder(lib.Key, "", t, 0, 500)
// 		if err != nil {
// 			return fae.Wrapf(err, "get library section folder: %s", lib.Key)
// 		}
//
// 		_, total, err := a.Plex.GetLibrarySection(lib.Key, "all", t, 0, 1)
// 		if err != nil {
// 			return fae.Wrapf(err, "get library section: %s", lib.Key)
// 		}
//
// 		for i := int64(0); i < total; i += 50 {
// 			list, _, err := a.Plex.GetLibrarySection(lib.Key, "all", t, int(i), 50)
// 			if err != nil {
// 				return fae.Wrap(err, "get library section")
// 			}
// 			for _, item := range list {
// 				if len(item.Media) == 0 {
// 					continue
// 				}
//
// 				for _, media := range item.Media {
// 					if len(media.Part) == 0 {
// 						continue
// 					}
//
// 					for _, part := range media.Part {
// 						if part.File == "" {
// 							continue
// 						}
//
// 						if _, ok := cache.files[part.File]; !ok {
// 							cache.files[part.File] = item
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
//
// 	return cache, nil
// }

func (c *plexFileCache) buildFiles(ctx context.Context) error {
	a := ContextApp(ctx)

	files := make(map[string]*plex.LeavesMetadata)

	for _, lib := range c.libs {
		t := plexLibType(lib.Type)
		if t == "" {
			continue
		}

		_, total, err := a.Plex.GetLibrarySection(lib.Key, "all", t, 0, 1)
		if err != nil {
			return fae.Wrapf(err, "get library section: %s", lib.Key)
		}

		for i := int64(0); i < total; i += 50 {
			list, _, err := a.Plex.GetLibrarySection(lib.Key, "all", t, int(i), 50)
			if err != nil {
				return fae.Wrap(err, "get library section")
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

						files[part.File] = item
					}
				}
			}
		}
	}

	c.files = files
	return nil
}

func (c *plexFileCache) update(ctx context.Context, title string, section string, libtype string) error {
	a := ContextApp(ctx)
	parent := c.parents[title]
	t := plexLibType(libtype)
	if t == "" {
		return fae.Errorf("unknown library type: %s", libtype)
	}

	_, total, err := a.Plex.GetLibrarySectionFolder(section, parent, t, 0, 1)
	if err != nil {
		return fae.Wrapf(err, "get library section: %s", section)
	}

	for i := int64(0); i < total; i += 50 {
		list, _, err := a.Plex.GetLibrarySection(section, parent, t, int(i), 50)
		if err != nil {
			return fae.Wrap(err, "get library section")
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

					c.files[part.File] = item
				}
			}
		}
	}
	return nil
}

func (c *plexFileCache) buildFolders(ctx context.Context) error {
	a := ContextApp(ctx)
	parents := make(map[string]string)

	for _, lib := range c.libs {
		t := plexLibType(lib.Type)
		if t == "" {
			continue
		}

		_, total, err := a.Plex.GetLibrarySectionFolder(lib.Key, "", t, 0, 1)
		if err != nil {
			return fae.Wrapf(err, "get library section folder: %s", lib.Key)
		}

		for i := int64(0); i < total; i += 50 {
			list, _, err := a.Plex.GetLibrarySectionFolder(lib.Key, "all", t, int(i), 50)
			if err != nil {
				return fae.Wrap(err, "get library section")
			}
			for _, item := range list {
				u, err := url.Parse(item.Key)
				if err != nil {
					return fae.Wrap(err, "parse uri")
				}
				m, err := url.ParseQuery(u.RawQuery)
				if err != nil {
					return fae.Wrap(err, "parse query")
				}
				parent := m["parent"][0]

				// a.Log.Debugf("folder: %s -> %s", item.Title, parent)
				if _, ok := parents[item.Title]; !ok {
					parents[item.Title] = parent
				}
			}
		}
	}

	c.parents = parents
	return nil
}
