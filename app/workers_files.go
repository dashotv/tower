package app

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"sync/atomic"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/dashotv/fae"
	"github.com/dashotv/minion"
)

var KINDS = []string{"movies", "movies3d", "movies4k", "movies4h", "kids", "tv", "anime", "donghua", "ecchi"}

var walking uint32

type FileWalk struct {
	minion.WorkerDefaults[*FileWalk]
}

func (j *FileWalk) Kind() string { return "file_walk" }
func (j *FileWalk) Timeout(job *minion.Job[*FileWalk]) time.Duration {
	return 60 * time.Minute
}
func (j *FileWalk) Work(ctx context.Context, job *minion.Job[*FileWalk]) error {
	l := app.Log.Named("file_walk")
	if !atomic.CompareAndSwapUint32(&walking, 0, 1) {
		l.Warnf("walkFiles: already running")
		return fae.Errorf("already running")
	}
	defer atomic.StoreUint32(&walking, 0)

	libs, err := app.Plex.GetLibraries()
	if err != nil {
		l.Errorw("libs", "error", err)
		return fae.Wrap(err, "getting libraries")
	}

	w := newWalker(app.DB, l.Named("walker"), libs)
	if err := w.Walk(); err != nil {
		l.Errorw("walk", "error", err)
		return fae.Wrap(err, "walking")
	}

	app.Workers.Enqueue(&FileMatch{})
	return nil
}

type FileMatch struct {
	minion.WorkerDefaults[*FileMatch]
}

func (j *FileMatch) Kind() string { return "file_match" }
func (j *FileMatch) Timeout(job *minion.Job[*FileMatch]) time.Duration {
	return 60 * time.Minute
}
func (j *FileMatch) Work(ctx context.Context, job *minion.Job[*FileMatch]) error {
	l := app.Log.Named("files.match")

	for _, kind := range KINDS {
		dir := filepath.Join(app.Config.DirectoriesCompleted, kind)
		if err := app.Workers.Enqueue(&FileMatchDir{Dir: dir}); err != nil {
			l.Errorw("enqueue", "error", err)
			return fae.Wrap(err, "enqueue")
		}
	}

	return nil
}

type FileMatchMedium struct {
	minion.WorkerDefaults[*FileMatchMedium]
	ID string
}

func (j *FileMatchMedium) Kind() string { return "file_match_medium" }
func (j *FileMatchMedium) Timeout(job *minion.Job[*FileMatchMedium]) time.Duration {
	return 60 * time.Minute
}
func (j *FileMatchMedium) Work(ctx context.Context, job *minion.Job[*FileMatchMedium]) error {
	l := app.Log.Named("files.match.medium")

	m := &Medium{}
	if err := app.DB.Medium.Find(job.Args.ID, m); err != nil {
		l.Errorw("find", "error", err)
		return fae.Wrap(err, "finding")
	}

	dest := m.Destination()
	dir := filepath.Join(app.Config.DirectoriesCompleted, dest)
	if err := app.Workers.Enqueue(&FileMatchDir{Dir: dir}); err != nil {
		l.Errorw("enqueue", "error", err)
		return fae.Wrap(err, "enqueue")
	}

	return nil
}

type FileMatchDir struct {
	minion.WorkerDefaults[*FileMatchDir]
	Dir string
}

func (j *FileMatchDir) Kind() string { return "file_match_dir" }
func (j *FileMatchDir) Timeout(job *minion.Job[*FileMatchDir]) time.Duration {
	return 60 * time.Minute
}
func (j *FileMatchDir) Work(ctx context.Context, job *minion.Job[*FileMatchDir]) error {
	dir := job.Args.Dir
	l := app.Log.Named("files.match.dir").With("dir", dir)
	l.Debugf("running")

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			l.Errorw("walk", "error", err)
			return fae.Wrap(err, "walking")
		}

		if d.IsDir() {
			if path != dir {
				return app.Workers.Enqueue(&FileMatchDir{Dir: path})
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
			l.Errorw("parts", "error", err)
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

		m.Completed = true
		m.Paths = append(m.Paths, &Path{Type: primitive.Symbol(filetype), Local: local, Extension: ext})
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
