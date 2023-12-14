package app

import (
	"context"
	"os"

	vidio "github.com/AlexEidt/Vidio"
	"github.com/pkg/errors"
	"github.com/samber/lo"

	"github.com/dashotv/minion"
)

type MediaPaths struct {
	ID string // medium
}

func (j *MediaPaths) Kind() string { return "MediaPaths" }
func (j *MediaPaths) Work(ctx context.Context, job *minion.Job[*MediaPaths]) error {
	m := &Medium{}
	if err := db.Medium.Find(job.Args.ID, m); err != nil {
		return errors.Wrap(err, "find medium")
	}

	queuedPaths := map[string]int{}

	for _, p := range m.Paths {
		if queuedPaths[p.LocalPath()] == 0 {
			workers.Log.Debugf("path import: %s", p.LocalPath())
			if err := workers.Enqueue(&PathImport{m.ID.Hex(), p.Id.Hex(), p.LocalPath()}); err != nil {
				return errors.Wrap(err, "enqueue path import")
			}
			queuedPaths[p.LocalPath()]++
		}
	}

	if m.Type == "Series" {
		eps, err := db.Episode.Query().
			Where("series_id", m.ID).
			Desc("season_number").Desc("episode_number").Desc("absolute_number").
			Limit(-1).
			Run()
		if err != nil {
			return errors.Wrap(err, "find episodes")
		}

		for _, e := range eps {
			if len(e.Paths) > 0 {
				for _, p := range e.Paths {
					if queuedPaths[p.LocalPath()] == 0 {
						workers.Log.Debugf("path import: %s", p.LocalPath())
						if err := workers.Enqueue(&PathImport{e.ID.Hex(), p.Id.Hex(), p.LocalPath()}); err != nil {
							return errors.Wrap(err, "enqueue path import")
						}
						queuedPaths[p.LocalPath()]++
					}
				}
			}
		}
	}

	return nil
}

type PathImport struct {
	ID     string // medium
	PathID string // path
	Title  string
}

func (j *PathImport) Kind() string { return "PathImport" }
func (j *PathImport) Work(ctx context.Context, job *minion.Job[*PathImport]) error {
	m := &Medium{}
	if err := db.Medium.Find(job.Args.ID, m); err != nil {
		return errors.Wrap(err, "find medium")
	}

	list := lo.Filter(m.Paths, func(p *Path, i int) bool {
		return p.Id.Hex() == job.Args.PathID
	})
	if len(list) == 0 {
		return errors.New("path not found")
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
		sum, err := sumFile(path.LocalPath())
		if err != nil {
			return errors.Wrap(err, "sum file")
		}
		path.Checksum = sum

		if v, err := vidio.NewVideo(path.LocalPath()); err == nil {
			path.Resolution = v.Height()
			path.Bitrate = v.Bitrate()
		} else {
			workers.Log.Warnf("failed to get video info: %s", err)
		}
	}

	if err := db.Medium.Save(m); err != nil {
		return errors.Wrap(err, "save path")
	}

	return nil
}
