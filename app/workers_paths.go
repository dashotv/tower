package app

import (
	"context"
	"os"

	vidio "github.com/AlexEidt/Vidio"
	"github.com/pkg/errors"
	"github.com/samber/lo"

	"github.com/dashotv/minion"
)

type PathImport struct {
	minion.WorkerDefaults[*PathImport]
	ID     string // medium
	PathID string // path
	Title  string
}

func (j *PathImport) Kind() string { return "PathImport" }
func (j *PathImport) Work(ctx context.Context, job *minion.Job[*PathImport]) error {
	m := &Medium{}
	if err := app.DB.Medium.Find(job.Args.ID, m); err != nil {
		return errors.Wrap(err, "find medium")
	}

	paths := m.Paths
	if m.Type == "Series" {
		eps, err := app.DB.Episode.Query().
			Where("series_id", m.ID).
			Limit(-1).
			Run()
		if err != nil {
			return errors.Wrap(err, "find episodes")
		}

		for _, e := range eps {
			paths = append(paths, e.Paths...)
		}
	}

	list := lo.Filter(paths, func(p *Path, i int) bool {
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

	queuedPaths := map[string]int{}

	newPaths := []*Path{}
	for _, p := range m.Paths {
		if queuedPaths[p.LocalPath()] == 0 {
			app.Workers.Log.Debugf("path import: %s", p.LocalPath())
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
						app.Workers.Log.Debugf("path import: %s", p.LocalPath())
						if err := app.Workers.Enqueue(&PathImport{ID: m.ID.Hex(), PathID: p.Id.Hex(), Title: p.LocalPath()}); err != nil {
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
