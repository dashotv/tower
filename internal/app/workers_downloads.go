package app

import (
	"context"
	"time"

	"github.com/dashotv/fae"
	"github.com/dashotv/minion"
)

// TODO: Create Movie Downloads
// TODO: Use runic for searches
// TODO: Remove Downloads for episodes whose release date change
// TODO: shift more towards want / evented downloads

var downloadProcessMutex = &CtxMutex{ch: make(chan struct{}, 1)}

type DownloadsProcess struct {
	minion.WorkerDefaults[*DownloadsProcess]
}

func (j *DownloadsProcess) Kind() string { return "DownloadsProcess" }
func (j *DownloadsProcess) Work(ctx context.Context, job *minion.Job[*DownloadsProcess]) error {
	muctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	if !downloadProcessMutex.Lock(muctx) {
		app.Log.Named("DownloadsProcess").Warn("failed to lock mutex")
		return nil
	}
	defer downloadProcessMutex.Unlock()

	a := ContextApp(ctx)

	// defer TickTock("DownloadsProcess: start")()
	// notifier.Info("Downloads", "processing downloads")
	funcs := []func() error{
		a.downloadsMove,
		a.downloadsCreate,
		a.downloadsSearch,
		a.downloadsLoad,
		a.downloadsManage,
	}

	for _, f := range funcs {
		err := f()
		if err != nil {
			app.Log.Named("DownloadsProcess").Errorf("failed to process downloads: %s", err)
			return fae.Wrap(err, "failed to process downloads")
		}
	}

	return nil
}

type DownloadsProcessLoad struct {
	minion.WorkerDefaults[*DownloadsProcessLoad]
}

func (j *DownloadsProcessLoad) Kind() string { return "downloads_process_load" }
func (j *DownloadsProcessLoad) Work(ctx context.Context, job *minion.Job[*DownloadsProcessLoad]) error {
	a := ContextApp(ctx)
	//args := job.Args
	return a.downloadsLoad()
}

type DownloadsMovies struct {
	minion.WorkerDefaults[*DownloadsMovies]
}

func (j *DownloadsMovies) Kind() string { return "downloads_movies" }
func (j *DownloadsMovies) Work(ctx context.Context, job *minion.Job[*DownloadsMovies]) error {
	a := ContextApp(ctx)
	return a.downloadsSearchMovies()
}
