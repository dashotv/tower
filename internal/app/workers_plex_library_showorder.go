package app

import (
	"context"
	"strings"

	"github.com/dashotv/fae"
	"github.com/dashotv/minion"
)

type PlexLibraryShoworder struct {
	minion.WorkerDefaults[*PlexLibraryShoworder]
}

func (j *PlexLibraryShoworder) Kind() string { return "plex_library_showorder" }
func (j *PlexLibraryShoworder) Work(ctx context.Context, job *minion.Job[*PlexLibraryShoworder]) error {
	a := ContextApp(ctx)
	l := a.Workers.Log.Named("plex_library_showorder")
	libs, err := a.Plex.GetLibraries()
	if err != nil {
		return fae.Wrap(err, "getting plex libs")
	}
	for _, lib := range libs {
		l.Debugw("lib", "title", lib.Title, "type", lib.Type)
		if lib.Type != "show" || !isAnimeKind(strings.ToLower(lib.Title)) {
			continue
		}

		t := plexLibType(lib.Type)
		if t == "" {
			continue
		}

		_, total, err := a.Plex.Search("", lib.Key, map[string]string{}, 0, 1)
		if err != nil {
			return fae.Wrapf(err, "get library section: %s", lib.Key)
		}

		for i := int64(0); i < total; i += 50 {
			list, _, err := a.Plex.Search("", lib.Key, map[string]string{}, int(i), 50)
			if err != nil {
				return fae.Wrap(err, "get library section")
			}
			for _, item := range list {
				if err := a.Plex.PutMetadataPrefs(item.RatingKey, map[string]string{"showOrdering": "tvdbAbsolute"}); err != nil {
					l.Errorf("put metadata prefs: [%s] %s: %v", item.RatingKey, item.Title, err)
				}
			}
		}
	}
	return nil
}
