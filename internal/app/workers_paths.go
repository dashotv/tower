package app

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/dashotv/fae"
	"github.com/dashotv/minion"
)

var TYPES = []string{"Movie", "Series", "Episode"}

type PathManageAll struct {
	minion.WorkerDefaults[*PathManageAll]
}

func (j *PathManageAll) Kind() string { return "path_manage_all" }
func (j *PathManageAll) Work(ctx context.Context, job *minion.Job[*PathManageAll]) error {
	a := ContextApp(ctx)
	//l := a.Workers.Log.Named("path_manage_all")
	err := a.DB.Medium.Query().In("_type", []string{"Movie", "Series"}).Where("broken", false).Each(100, func(m *Medium) error {
		return a.Workers.Enqueue(&PathManage{MediumID: m.ID.Hex()})
	})
	if err != nil {
		return fae.Wrap(err, "find media")
	}
	return nil
}

// PathManage removes missing paths and updates path metadata from Plex
type PathManage struct {
	minion.WorkerDefaults[*PathManage]
	MediumID string `bson:"medium_id" json:"medium_id"`
}

func (j *PathManage) Kind() string { return "path_manage" }
func (j *PathManage) Work(ctx context.Context, job *minion.Job[*PathManage]) error {
	a := ContextApp(ctx)
	if a == nil {
		return fae.New("no app in context")
	}

	// l := a.Workers.Log.Named("path_manage")
	MediumID := job.Args.MediumID

	tctx, timeout := context.WithTimeout(ctx, 20*time.Second)
	defer timeout()
	ok := checkContextTimeout(tctx, func() bool {
		return a.PlexFileCache != nil && a.PlexFileCache.files != nil
	})
	if !ok {
		return fae.New("no plex file cache")
	}

	media := []*Medium{}
	medium := &Medium{}
	if err := app.DB.Medium.Find(MediumID, medium); err != nil {
		return fae.Wrap(err, "find medium")
	}
	kind := medium.Kind

	lib, ok := a.Libs[string(medium.Kind)]
	if !ok {
		return fae.Errorf("library not found: %s", medium.Kind)
	}

	dir := fmt.Sprintf("%s/%s", lib.Path, medium.Directory)
	if !exists(fmt.Sprintf("%s/%s", lib.Path, medium.Directory)) {
		medium.Broken = true
		if err := app.DB.Medium.Save(medium); err != nil {
			return fae.Wrap(err, "save medium")
		}
		return fae.Errorf("directory not found: %s", dir)
	}

	if err := a.fileMatchDir(dir); err != nil {
		return fae.Wrap(err, "file match dir")
	}
	if err := a.filePlexmatch(medium); err != nil {
		return fae.Wrap(err, "file plexmatch")
	}

	media = append(media, medium)
	if medium.Type == "Series" {
		// remove any paths that are not covers or backgrounds, videos should be on the episode not the series
		medium.Paths = lo.Filter(medium.Paths, func(p *Path, i int) bool {
			return p.IsCoverBackground()
		})
		// add episodes to list
		err := app.DB.Medium.Query().Where("_type", "Episode").Where("series_id", medium.ID).Each(100, func(e *Medium) error {
			media = append(media, e)
			return nil
		})
		if err != nil {
			return fae.Wrap(err, "find episodes")
		}
	}

	for _, m := range media {
		newPaths := map[string]*Path{}
		for _, path := range m.Paths {
			if !path.Exists() && !path.IsCoverBackground() {
				a.Log.Warnf("path does not exist: %s", path.LocalPath())
				continue
			}

			if newPaths[path.LocalPath()] != nil {
				a.Log.Warnf("duplicate path: %s", path.LocalPath())
				continue
			}

			newPaths[path.LocalPath()] = path
			if !path.IsVideo() {
				// keep path, but don't process
				continue
			}

			if err := a.pathImport(path); err != nil {
				return fae.Wrap(err, "path import")
			}
			if err := a.pathDest(m, kind, path); err != nil {
				return fae.Wrap(err, "path check")
			}
		}

		m.Paths = lo.Values(newPaths)
		if len(m.Paths) > 0 {
			m.Downloaded = true
			m.Completed = true
		}
		if err := app.DB.Medium.Save(m); err != nil {
			return fae.Wrap(err, "save path")
		}
	}
	return nil
}

func (a *Application) pathDest(m *Medium, kind primitive.Symbol, path *Path) error {
	d, err := a.Destinator.Destination(kind, m)
	if err != nil {
		return fae.Wrap(err, "destination")
	}
	dest := fmt.Sprintf("%s.%s", d, path.Extension)

	path.Rename = false
	if path.LocalPath() != dest {
		a.Log.Debugw("pathcheck", "path", path.LocalPath(), "dest", dest)
		path.Rename = true
	}

	return nil
}

func (a *Application) pathImport(path *Path) error {
	f := path.LocalPath()
	if a.PlexFileCache.files[f] == nil {
		a.Log.Warnf("path not in cache: %s", f)
		return nil
	}

	meta := a.PlexFileCache.files[f]
	if len(meta.Media) == 0 {
		a.Log.Warnf("no media in cache: %s", f)
		return nil
	}
	path.Bitrate = int(meta.Media[0].Bitrate)
	path.Resolution = meta.Media[0].GetVideoResolution()

	path.ParseTag()

	if len(meta.Media[0].Part) == 0 {
		a.Log.Warnf("no parts in cache: %s", f)
		return nil
	}
	path.Size = meta.Media[0].Part[0].Size

	return nil
}

func (a *Application) filePlexmatch(medium *Medium) error {
	if medium.Type != "Series" && medium.Type != "Movie" {
		return nil
	}

	lib := a.Libs[string(medium.Kind)]
	if lib == nil {
		return fae.Errorf("library not found: %s", medium.Kind)
	}

	file := fmt.Sprintf("%s/%s/.plexmatch", lib.Path, medium.Directory)
	data := []string{}
	data = append(data, "# PlexMatch - managed by dashotv")
	data = append(data, fmt.Sprintf("Title: %s", medium.DisplayTitle()))
	data = append(data, fmt.Sprintf("Year: %s", medium.Year()))
	if medium.Source == "tvdb" {
		data = append(data, fmt.Sprintf("tvdbid: %s", medium.SourceID))
	} else if medium.Source == "tmdb" {
		data = append(data, fmt.Sprintf("tmdbid: %s", medium.SourceID))
		if medium.ImdbID != "" {
			data = append(data, fmt.Sprintf("imdbid: %s", medium.ImdbID))
		}
	}
	data = append(data, "")
	// TODO: handle episodes as well?

	if err := os.WriteFile(file, []byte(strings.Join(data, "\n")), 0644); err != nil {
		return fae.Wrap(err, "writing plexmatch")
	}
	return nil
}

func (a *Application) fileMatchDir(dir string) error {
	l := a.Log.Named("file_match_dir")
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			l.Errorw("walk", "error", err)
			return fae.Wrap(err, "walking")
		}

		if d.IsDir() {
			if filepath.Base(path)[0] == '@' || filepath.Base(path)[0] == '.' {
				// skip directories starting with .
				// skip directories starting with @ (e.g. @eaDir from synology)
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Base(path)[0] == '.' {
			return nil
		}

		filetype := fileType(path)
		if filetype == "" {
			l.Warnf("path: unknown type: %s", path)
			return nil
		}

		kind, name, file, ext, err := pathParts(path)
		if err != nil {
			l.Errorw("parts", "error", err, "path", path)
			return nil
		}
		local := fmt.Sprintf("%s/%s/%s", kind, name, file)

		m, ok, err := app.DB.MediumBy(kind, name, file, ext)
		if err != nil {
			l.Errorw("medium", "error", err)
			return nil
		}
		if ok {
			return nil
		}
		if m == nil {
			l.Warnw("medium", "not found", local)
			return nil
		}

		p := &Path{Type: primitive.Symbol(filetype), Local: local, Extension: ext}
		// l.Debugw("adding", "path", p.LocalPath(), "ext", ext)
		m.Completed = true
		m.Paths = append(m.Paths, p)
		if err := app.DB.Medium.Save(m); err != nil {
			l.Errorw("save", "error", err)
			return fae.Wrap(err, "saving")
		}

		return nil
	})
	if err != nil {
		l.Errorw("walk", "error", err)
		return fae.Wrap(err, "walking")
	}
	return nil
}

type PathImport struct {
	minion.WorkerDefaults[*PathImport]
	ID     string `bson:"id" json:"id"`           // medium
	PathID string `bson:"path_id" json:"path_id"` // path
	Title  string `bson:"title" json:"title"`
}

func (j *PathImport) Kind() string                                       { return "PathImport" }
func (j *PathImport) Timeout(job *minion.Job[*PathImport]) time.Duration { return 300 * time.Second }
func (j *PathImport) Work(ctx context.Context, job *minion.Job[*PathImport]) error {
	a := ContextApp(ctx)
	if a == nil {
		return fae.New("no app in context")
	}

	m := &Medium{}
	if err := a.DB.Medium.Find(job.Args.ID, m); err != nil {
		return fae.Wrap(err, "find medium")
	}

	list := lo.Filter(m.Paths, func(p *Path, i int) bool {
		return p.ID.Hex() == job.Args.PathID
	})
	if len(list) != 1 {
		return fae.New("matching path")
	}

	path := list[0]
	if !path.Exists() {
		return nil
	}

	if !path.IsVideo() {
		// keep path, but don't process
		return nil
	}

	if err := a.pathImport(path); err != nil {
		return fae.Wrap(err, "path import")
	}

	if err := app.DB.Medium.Save(m); err != nil {
		return fae.Wrap(err, "save path")
	}

	return nil
}

// type PathCleanup struct {
// 	minion.WorkerDefaults[*PathCleanup]
// 	ID string // medium
// }
//
// func (j *PathCleanup) Kind() string { return "PathCleanup" }
// func (j *PathCleanup) Work(ctx context.Context, job *minion.Job[*PathCleanup]) error {
// 	l := app.Log.Named("path.cleanup")
// 	m := &Medium{}
// 	if err := app.DB.Medium.Find(job.Args.ID, m); err != nil {
// 		l.Errorf("find medium: %s", err)
// 		return fae.Wrap(err, "find medium")
// 	}
//
// 	queuedPaths := map[string]int{}
// 	newPaths := []*Path{}
// 	for _, p := range m.Paths {
// 		if !p.Exists() {
// 			continue
// 		}
// 		if queuedPaths[p.LocalPath()] == 0 {
// 			if err := app.Workers.Enqueue(&PathImport{ID: m.ID.Hex(), PathID: p.ID.Hex(), Title: p.LocalPath()}); err != nil {
// 				return fae.Wrap(err, "enqueue path import")
// 			}
// 			queuedPaths[p.LocalPath()]++
// 			newPaths = append(newPaths, p)
// 		}
// 	}
//
// 	m.Paths = newPaths
// 	if err := app.DB.Medium.Save(m); err != nil {
// 		return fae.Wrap(err, "save medium")
// 	}
//
// 	if m.Type == "Series" {
// 		err := app.DB.Episode.Query().Where("series_id", m.ID).Batch(100, func(results []*Episode) error {
// 			for _, e := range results {
// 				if err := app.Workers.Enqueue(&PathCleanup{ID: e.ID.Hex()}); err != nil {
// 					return fae.Wrap(err, "enqueue media paths")
// 				}
// 			}
//
// 			return nil
// 		})
// 		if err != nil {
// 			return fae.Wrap(err, "series batch")
// 		}
// 	}
//
// 	return nil
// }
//
// type PathCleanupAll struct {
// 	minion.WorkerDefaults[*PathCleanupAll]
// }
//
// func (j *PathCleanupAll) Kind() string { return "PathCleanupAll" }
// func (j *PathCleanupAll) Work(ctx context.Context, job *minion.Job[*PathCleanupAll]) error {
// 	for _, t := range TYPES {
// 		total, err := app.DB.Medium.Query().Where("_type", t).Count()
// 		if err != nil {
// 			return fae.Wrap(err, "count media")
// 		}
// 		if total == 0 {
// 			continue
// 		}
// 		for skip := 0; skip < int(total); skip += 100 {
// 			media, err := app.DB.Medium.Query().Limit(100).Skip(skip).Run()
// 			if err != nil {
// 				return fae.Wrap(err, "find media")
// 			}
//
// 			for _, m := range media {
// 				if err := app.Workers.Enqueue(&PathCleanup{ID: m.ID.Hex()}); err != nil {
// 					return fae.Wrap(err, "enqueue path cleanup")
// 				}
// 			}
// 		}
// 	}
//
// 	return nil
// }

type PathDeleteAll struct {
	minion.WorkerDefaults[*PathDeleteAll]
	MediumID string `bson:"medium_id" json:"medium_id"`
}

func (j *PathDeleteAll) Kind() string { return "path_delete_all" }
func (j *PathDeleteAll) Work(ctx context.Context, job *minion.Job[*PathDeleteAll]) error {
	id := job.Args.MediumID

	m := &Medium{}
	if err := app.DB.Medium.Find(id, m); err != nil {
		return fae.Wrap(err, "find medium")
	}

	paths := m.Paths
	if m.Type == "Series" {
		err := app.DB.Episode.Query().Where("series_id", m.ID).Batch(100, func(list []*Episode) error {
			for _, e := range list {
				paths = append(paths, e.Paths...)
			}
			return nil
		})
		if err != nil {
			return fae.Wrap(err, "listing episodes")
		}
	}

	for _, p := range paths {
		if !p.Exists() {
			continue
		}
		if err := os.Remove(p.LocalPath()); err != nil {
			return fae.Wrap(err, "remove path")
		}
	}

	if err := app.DB.Medium.Delete(m); err != nil { // TODO: why?
		return fae.Wrap(err, "delete medium")
	}

	return nil
}

type PathDelete struct {
	minion.WorkerDefaults[*PathDelete]
	MediumID string `bson:"medium_id" json:"medium_id"` // medium
	PathID   string `bson:"path_id" json:"path_id"`     // path
	Title    string `bson:"title" json:"title"`
}

func (j *PathDelete) Kind() string { return "path_delete" }
func (j *PathDelete) Work(ctx context.Context, job *minion.Job[*PathDelete]) error {
	a := ContextApp(ctx)
	if a == nil {
		return fae.New("no app in context")
	}

	medium_id := job.Args.MediumID
	path_id := job.Args.PathID

	oid, err := primitive.ObjectIDFromHex(path_id)
	if err != nil {
		return fae.Wrap(err, "invalid id")
	}

	media, err := a.DB.Medium.Query().Where("paths._id", oid).Run()
	if err != nil {
		return fae.Wrap(err, "error querying media")
	}
	if len(media) == 0 {
		return fae.Wrap(err, "no media found")
	}
	if len(media) > 1 {
		return fae.Wrap(err, "duplicate media found")
	}

	m := media[0]
	list := lo.Filter(m.Paths, func(p *Path, i int) bool {
		return p.ID != oid
	})

	removed, _ := lo.Difference(m.Paths, list)
	if len(removed) > 0 {
		for _, p := range removed {
			a.Log.Named("path_delete").Debugf("removing path: %s %s", path_id, p.LocalPath())
			if p.Exists() {
				if err := os.Remove(p.LocalPath()); err != nil {
					return fae.Wrap(err, "removing path")
				}
			}
		}
	}

	m.Paths = list
	if err := a.DB.Medium.Save(m); err != nil {
		return fae.Wrap(err, "error saving Paths")
	}

	if medium_id != m.ID.Hex() {
		medium, err := a.DB.Medium.Get(medium_id, &Medium{})
		if err != nil {
			return fae.Wrap(err, "error getting medium")
		}
		if err := a.DB.Medium.Save(medium); err != nil {
			return fae.Wrap(err, "error saving medium")
		}
	}

	return nil
}
