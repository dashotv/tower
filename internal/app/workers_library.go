package app

import (
	"context"

	"github.com/dashotv/fae"
	"github.com/dashotv/minion"
)

type LibraryCounts struct {
	minion.WorkerDefaults[*LibraryCounts]
}

func (j *LibraryCounts) Kind() string { return "library_counts" }
func (j *LibraryCounts) Work(ctx context.Context, job *minion.Job[*LibraryCounts]) error {
	a := ContextApp(ctx)
	//l := a.Workers.Log.Named("library_counts")
	list, err := a.DB.LibraryList(1, -1)
	if err != nil {
		return fae.Wrap(err, "library list")
	}

	for _, lib := range list {
		total, err := a.DB.File.Query().Where("library_id", lib.ID).Count()
		if err != nil {
			return fae.Wrap(err, "file count")
		}
		lib.Count = total
		if err := a.DB.Library.Save(lib); err != nil {
			return fae.Wrap(err, "library update")
		}
	}
	return nil
}
