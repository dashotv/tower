package app

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/samber/lo"

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
		return fmt.Errorf("already running")
	}
	defer atomic.StoreUint32(&walking, 0)

	libs, err := app.Plex.GetLibraries()
	if err != nil {
		l.Errorw("libs", "error", err)
		return fmt.Errorf("getting libraries: %w", err)
	}

	w := newWalker(app.DB, l.Named("walker"), libs)
	if err := w.Walk(); err != nil {
		l.Errorw("walk", "error", err)
		return fmt.Errorf("walking: %w", err)
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

	start := time.Now()
	found := 0
	missing := 0
	existing := 0
	defer func() {
		l.Debugf("duration: %d, found: %d, existing: %d, missing: %d", time.Since(start), found, existing, missing)
	}()

	for _, kind := range KINDS {
		dir := filepath.Join(app.Config.DirectoriesCompleted, kind)
		l.Infof("walking: %s", dir)
		err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				l.Errorw("walk", "error", err)
				return fmt.Errorf("walking: %w", err)
			}

			if d.IsDir() {
				return nil
			}

			if filepath.Base(path)[0] == '.' {
				return nil
			}

			ext := filepath.Ext(path)
			if ext == "" || !lo.Contains(app.Config.ExtensionsVideo, ext[1:]) {
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
				existing += 1
				return nil
			}
			if m == nil {
				missing += 1
				l.Warnw("medium", "not found", local)
				return nil
			}

			// l.Warnw("found", "path", path)
			found += 1

			m.Paths = append(m.Paths, &Path{Type: "video", Local: local, Extension: ext})
			if err := app.DB.Medium.Save(m); err != nil {
				l.Errorw("save", "error", err)
				return fmt.Errorf("saving: %w", err)
			}

			return nil
		})
		if err != nil {
			l.Errorw("walk", "error", err)
			return fmt.Errorf("walking: %w", err)
		}
	}

	return nil
}
