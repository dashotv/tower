package app

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/dashotv/minion"
)

var walking uint32

type FileWalk struct {
	minion.WorkerDefaults[*FileWalk]
}

func (j *FileWalk) Kind() string { return "file_walk" }
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
	return nil
}
