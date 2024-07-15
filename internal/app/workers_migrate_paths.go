package app

import (
	"context"
	"fmt"
	"regexp"

	"github.com/dashotv/fae"
	"github.com/dashotv/minion"
)

var imageRegex = regexp.MustCompile(`\/(cover|background)`)

type MigratePaths struct {
	minion.WorkerDefaults[*MigratePaths]
}

func (j *MigratePaths) Kind() string { return "migrate_paths" }
func (j *MigratePaths) Work(ctx context.Context, job *minion.Job[*MigratePaths]) error {
	a := ContextApp(ctx)
	l := a.Workers.Log.Named("migrate_paths")

	found := []string{}
	failed := []string{}

	libs, err := a.DB.LibraryMap()
	if err != nil {
		return fae.Wrap(err, "library map")
	}

	q := a.DB.Medium.Query().In("_type", []string{"Movie", "Episode"}).Desc("created_at")
	total, err := q.Count()
	if err != nil {
		return fae.Wrap(err, "medium count")
	}
	l.Debugw("migrate paths", "total", total)
	defer TickTock("migrate_paths")()
	err = q.Batch(100, func(results []*Medium) error {
		for _, m := range results {
			for _, p := range m.Paths {
				if p.IsCoverBackground() {
					continue
				}

				f, err := a.DB.FileFindOrCreateByPath(p.LocalPath())
				if err != nil {
					return fae.Wrap(err, "failed to find or create file")
				}

				kind, _, file, ext, err := pathParts(fmt.Sprintf("%s.%s", p.Local, p.Extension))
				if err != nil {
					l.Errorf("failed to get path parts: %s: %s", p.LocalPath(), err)
					failed = append(failed, p.LocalPath())
					continue
				}

				found = append(found, p.LocalPath())

				lib := libs[kind]
				f.LibraryID = lib.ID
				f.MediumID = m.ID

				f.Path = p.LocalPath()
				f.Type = fileType(p.LocalPath())
				f.Name = file
				f.Extension = ext
				f.Size = p.Size
				f.Resolution = p.Resolution
				f.Checksum = p.Checksum

				// l.Debugw("file", "path", f.Path, "name", f.Name, "extension", f.Extension, "size", f.Size, "resolution", f.Resolution, "checksum", f.Checksum)
				if err := a.DB.File.Save(f); err != nil {
					return fae.Wrap(err, "saving file")
				}
			}
		}
		return nil
	})
	if err != nil {
		return fae.Wrap(err, "medium batch")
	}

	l.Infow("migrate paths", "found", len(found), "failed", len(failed))
	for _, f := range failed {
		l.Warnw("failed", "path", f)
	}
	return nil
}
