package app

import (
	"context"
	"time"

	"github.com/dashotv/fae"
	"github.com/dashotv/minion"
)

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
		a.downloadsCreate,
		a.downloadsSearch,
		a.downloadsLoad,
		a.downloadsManage,
		a.downloadsMove,
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
