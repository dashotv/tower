package app

import (
	"context"
	"fmt"
	"os"
	"time"

	vidio "github.com/AlexEidt/Vidio"
	"github.com/pkg/errors"
	"github.com/samber/lo"

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
		return errors.Wrap(err, "find medium")
	}

	list := lo.Filter(m.Paths, func(p *Path, i int) bool {
		return p.Id.Hex() == job.Args.PathID
	})
	if len(list) == 0 {
		return errors.New("no matching path in list")
	}
	if len(list) > 1 {
		return errors.New("multiple paths found")
	}

	path := list[0]
	if !path.Exists() {
		return errors.Errorf("path does not exist: %s", path.LocalPath())
	}

	stat, err := os.Stat(path.LocalPath())
	if err != nil {
		return errors.Wrap(err, "stat path")
	}

	path.UpdatedAt = stat.ModTime()
	path.Size = int(stat.Size())

	if path.IsVideo() && lo.Contains(app.Config.ExtensionsSubtitles, path.Extension) {
		path.Type = "subtitle"
	}

	if path.IsVideo() {
		if sum, err := sumFile(path.LocalPath()); err == nil {
			path.Checksum = sum
		} else {
			app.Workers.Log.Errorf("failed to checksum file: %s", err)
		}

		if v, err := vidio.NewVideo(path.LocalPath()); err == nil {
			path.Resolution = v.Height()
			path.Bitrate = v.Bitrate()
		} else {
			app.Workers.Log.Warnf("failed to get video info: %s", err)
		}
	}

	if err := app.DB.Medium.Save(m); err != nil {
		return errors.Wrap(err, "save path")
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
		return fmt.Errorf("find medium: %w", err)
	}

	queuedPaths := map[string]int{}
	newPaths := []*Path{}
	for _, p := range m.Paths {
		if !p.Exists() {
			continue
		}
		if queuedPaths[p.LocalPath()] == 0 {
			// app.Workers.Log.Debugf("path import: %s", p.LocalPath())
			if err := app.Workers.Enqueue(&PathImport{ID: m.ID.Hex(), PathID: p.Id.Hex(), Title: p.LocalPath()}); err != nil {
				return errors.Wrap(err, "enqueue path import")
			}
			queuedPaths[p.LocalPath()]++
			newPaths = append(newPaths, p)
		}
	}

	m.Paths = newPaths
	if err := app.DB.Medium.Save(m); err != nil {
		return errors.Wrap(err, "save medium")
	}

	if m.Type == "Series" {
		q := app.DB.Episode.Query().Where("series_id", m.ID)

		count, err := q.Count()
		if err != nil {
			return errors.Wrap(err, "count episodes")
		}

		for skip := 0; skip < int(count); skip += 100 {
			eps, err := q.Limit(100).Skip(skip).Run()
			if err != nil {
				return errors.Wrap(err, "find episodes")
			}

			for _, e := range eps {
				if err := app.Workers.Enqueue(&MediaPaths{ID: e.ID.Hex()}); err != nil {
					return errors.Wrap(err, "enqueue media paths")
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
			return fmt.Errorf("count media: %w", err)
		}
		if total == 0 {
			continue
		}
		for skip := 0; skip < int(total); skip += 100 {
			media, err := app.DB.Medium.Query().Limit(100).Skip(skip).Run()
			if err != nil {
				return fmt.Errorf("find media: %w", err)
			}

			for _, m := range media {
				if err := app.Workers.Enqueue(&PathCleanup{ID: m.ID.Hex()}); err != nil {
					return fmt.Errorf("enqueue path cleanup: %w", err)
				}
			}
		}
	}

	return nil
}

type MediaPaths struct {
	minion.WorkerDefaults[*MediaPaths]
	ID string // medium
}

func (j *MediaPaths) Kind() string { return "MediaPaths" }
func (j *MediaPaths) Work(ctx context.Context, job *minion.Job[*MediaPaths]) error {
	m := &Medium{}
	if err := app.DB.Medium.Find(job.Args.ID, m); err != nil {
		return errors.Wrap(err, "find medium")
	}

	err := j.Cleanup(m)
	if err != nil {
		return errors.Wrap(err, "cleanup")
	}

	queuedPaths := map[string]int{}

	newPaths := []*Path{}
	for _, p := range m.Paths {
		if queuedPaths[p.LocalPath()] == 0 {
			// app.Workers.Log.Debugf("path import: %s", p.LocalPath())
			if err := app.Workers.Enqueue(&PathImport{ID: m.ID.Hex(), PathID: p.Id.Hex(), Title: p.LocalPath()}); err != nil {
				return errors.Wrap(err, "enqueue path import")
			}
			queuedPaths[p.LocalPath()]++
			newPaths = append(newPaths, p)
		}
	}
	m.Paths = newPaths
	if err := app.DB.Medium.Save(m); err != nil {
		return errors.Wrap(err, "save medium")
	}

	if m.Type == "Series" {
		eps, err := app.DB.Episode.Query().
			Where("series_id", m.ID).
			Desc("season_number").Desc("episode_number").Desc("absolute_number").
			Limit(-1).
			Run()
		if err != nil {
			return errors.Wrap(err, "find episodes")
		}

		for _, e := range eps {
			if len(e.Paths) > 0 {
				newPaths := []*Path{}
				for _, p := range e.Paths {
					if queuedPaths[p.LocalPath()] == 0 {
						// app.Workers.Log.Debugf("path import: %s", p.LocalPath())
						if err := app.Workers.Enqueue(&PathImport{ID: e.ID.Hex(), PathID: p.Id.Hex(), Title: p.LocalPath()}); err != nil {
							return errors.Wrap(err, "enqueue path import")
						}
						queuedPaths[p.LocalPath()]++
						newPaths = append(newPaths, p)
					}
				}
				e.Paths = newPaths
				if err := app.DB.Episode.Save(e); err != nil {
					return errors.Wrap(err, "save episode")
				}
			}
		}
	}

	return nil
}

func (j *MediaPaths) Cleanup(m *Medium) error {
	paths := []*Path{}
	for _, p := range m.Paths {
		if p.IsVideo() && lo.Contains(app.Config.ExtensionsSubtitles, p.Extension) {
			p.Type = "subtitle"
		}
		if p.Exists() {
			paths = append(paths, p)
		}
	}

	m.Paths = paths
	if err := app.DB.Medium.Save(m); err != nil {
		return errors.Wrap(err, "save medium")
	}

	if m.Type == "Series" {
		eps, err := app.DB.Episode.Query().
			Where("series_id", m.ID).
			Limit(-1).
			Run()
		if err != nil {
			return errors.Wrap(err, "find episodes")
		}

		for _, e := range eps {
			paths := []*Path{}
			for _, p := range e.Paths {
				if p.IsVideo() && lo.Contains(app.Config.ExtensionsSubtitles, p.Extension) {
					p.Type = "subtitle"
				}
				if p.Exists() {
					paths = append(paths, p)
				}
			}

			e.Paths = paths
			if err := app.DB.Episode.Save(e); err != nil {
				return errors.Wrap(err, "save episode")
			}
		}
	}

	return nil
}
