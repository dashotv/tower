package app

import (
	"context"
	"os"
	"time"

	"github.com/samber/lo"

	"github.com/dashotv/fae"
	"github.com/dashotv/minion"
)

var TYPES = []string{"Movie", "Series", "Episode"}

type PathImport struct {
	minion.WorkerDefaults[*PathImport]
	ID     string `bson:"id" json:"id"`           // medium
	PathID string `bson:"path_id" json:"path_id"` // path
	Title  string `bson:"title" json:"title"`
}

func (j *PathImport) Kind() string                                       { return "PathImport" }
func (j *PathImport) Timeout(job *minion.Job[*PathImport]) time.Duration { return 300 * time.Second }
func (j *PathImport) Work(ctx context.Context, job *minion.Job[*PathImport]) error {
	m := &Medium{}
	if err := app.DB.Medium.Find(job.Args.ID, m); err != nil {
		return fae.Wrap(err, "find medium")
	}

	list := lo.Filter(m.Paths, func(p *Path, i int) bool {
		return p.ID.Hex() == job.Args.PathID
	})
	if len(list) == 0 {
		return fae.New("no matching path in list")
	}
	if len(list) > 1 {
		return fae.New("multiple paths found")
	}

	path := list[0]
	if !path.Exists() {
		return fae.Errorf("path does not exist: %s", path.LocalPath())
	}

	stat, err := os.Stat(path.LocalPath())
	if err != nil {
		return fae.Wrap(err, "stat path")
	}

	path.UpdatedAt = stat.ModTime()
	path.Size = int(stat.Size())

	// if path.IsVideo() && lo.Contains(app.Config.ExtensionsSubtitles, path.Extension) {
	// 	path.Type = "subtitle"
	// }

	// 	if path.IsVideo() {
	// 		if sum, err := sumFile(path.LocalPath()); err == nil {
	// 			path.Checksum = sum
	// 		} else {
	// 			app.Workers.Log.Errorf("failed to checksum file: %s", err)
	// 		}
	//
	// 		if v, err := vidio.NewVideo(path.LocalPath()); err == nil {
	// 			path.Resolution = v.Height()
	// 			path.Bitrate = v.Bitrate()
	// 		} else {
	// 			app.Workers.Log.Warnf("failed to get video info: %s", err)
	// 		}
	// 	}

	if err := app.DB.Medium.Save(m); err != nil {
		return fae.Wrap(err, "save path")
	}

	return nil
}

type PathCleanup struct {
	minion.WorkerDefaults[*PathCleanup]
	ID string // medium
}

func (j *PathCleanup) Kind() string { return "PathCleanup" }
func (j *PathCleanup) Work(ctx context.Context, job *minion.Job[*PathCleanup]) error {
	l := app.Log.Named("path.cleanup")
	m := &Medium{}
	if err := app.DB.Medium.Find(job.Args.ID, m); err != nil {
		l.Errorf("find medium: %s", err)
		return fae.Wrap(err, "find medium")
	}

	queuedPaths := map[string]int{}
	newPaths := []*Path{}
	for _, p := range m.Paths {
		if !p.Exists() {
			continue
		}
		if queuedPaths[p.LocalPath()] == 0 {
			// app.Workers.Log.Debugf("path import: %s", p.LocalPath())
			if err := app.Workers.Enqueue(&PathImport{ID: m.ID.Hex(), PathID: p.ID.Hex(), Title: p.LocalPath()}); err != nil {
				return fae.Wrap(err, "enqueue path import")
			}
			queuedPaths[p.LocalPath()]++
			newPaths = append(newPaths, p)
		}
	}

	m.Paths = newPaths
	if err := app.DB.Medium.Save(m); err != nil {
		return fae.Wrap(err, "save medium")
	}

	if m.Type == "Series" {
		q := app.DB.Episode.Query().Where("series_id", m.ID)

		count, err := q.Count()
		if err != nil {
			return fae.Wrap(err, "count episodes")
		}

		for skip := 0; skip < int(count); skip += 100 {
			eps, err := q.Limit(100).Skip(skip).Run()
			if err != nil {
				return fae.Wrap(err, "find episodes")
			}

			for _, e := range eps {
				if err := app.Workers.Enqueue(&PathCleanup{ID: e.ID.Hex()}); err != nil {
					return fae.Wrap(err, "enqueue media paths")
				}
			}
		}
	}

	return nil
}

type PathCleanupAll struct {
	minion.WorkerDefaults[*PathCleanupAll]
}

func (j *PathCleanupAll) Kind() string { return "PathCleanupAll" }
func (j *PathCleanupAll) Work(ctx context.Context, job *minion.Job[*PathCleanupAll]) error {
	for _, t := range TYPES {
		total, err := app.DB.Medium.Query().Where("_type", t).Count()
		if err != nil {
			return fae.Wrap(err, "count media")
		}
		if total == 0 {
			continue
		}
		for skip := 0; skip < int(total); skip += 100 {
			media, err := app.DB.Medium.Query().Limit(100).Skip(skip).Run()
			if err != nil {
				return fae.Wrap(err, "find media")
			}

			for _, m := range media {
				if err := app.Workers.Enqueue(&PathCleanup{ID: m.ID.Hex()}); err != nil {
					return fae.Wrap(err, "enqueue path cleanup")
				}
			}
		}
	}

	return nil
}

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

	return nil
}
