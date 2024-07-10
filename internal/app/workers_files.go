package app

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"sync/atomic"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/samber/lo"

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
	a := ContextApp(ctx)
	l := a.Log.Named("file_walk")
	if !atomic.CompareAndSwapUint32(&walking, 0, 1) {
		l.Warnf("already running")
		return fae.Errorf("already running")
	}
	defer atomic.StoreUint32(&walking, 0)

	_, err := a.DB.File.Collection.UpdateMany(ctx, bson.M{}, bson.M{"$set": bson.M{"exists": false}})
	if err != nil {
		return fae.Wrap(err, "updating")
	}

	// eg := new(errgroup.Group)
	for _, lib := range a.Libs {
		lib := lib
		// eg.Go(func() error {
		err := filepath.WalkDir(lib.Path, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return fae.Wrap(err, "walking")
			}

			if d.IsDir() {
				if filepath.Base(path)[0] == '@' || filepath.Base(path)[0] == '.' {
					// skip directories starting with .
					// skip directories starting with @ (e.g. @eaDir from synology)
					return filepath.SkipDir
				}
				return nil
			}
			if filepath.Base(path)[0] == '.' {
				return nil
			}

			l.Debugf("path: %s", path)
			_, _, file, ext, err := pathParts(path)
			if err != nil {
				l.Warnf("parts: %s: %s", path, err)
				return nil
			}

			f, err := a.DB.FileFindOrCreateByPath(path)
			if err != nil {
				return fae.Wrap(err, "finding or creating")
			}

			i, err := d.Info()
			if err != nil {
				return fae.Wrap(err, "info")
			}

			f.LibraryID = lib.ID
			f.Type = fileType(path)
			f.Name = file
			f.Extension = ext
			f.ModifiedAt = i.ModTime().Unix()
			f.Size = i.Size()
			f.Type = fileType(path)
			f.Exists = true

			// sum, err := sumFile(path)
			// if err != nil {
			// 	return fae.Wrap(err, "summing")
			// }
			// f.Checksum = sum

			if err := a.DB.File.Save(f); err != nil {
				return fae.Wrap(err, "saving")
			}

			return nil
		})
		// })
		if err != nil {
			return fae.Wrap(err, "walking")
		}
	}

	// if err := eg.Wait(); err != nil {
	// 	return fae.Wrap(err, "walking")
	// }

	// app.Workers.Enqueue(&FileMatch{})
	return nil
}

// type FileMatch struct {
// 	minion.WorkerDefaults[*FileMatch]
// }
//
// func (j *FileMatch) Kind() string { return "file_match" }
// func (j *FileMatch) Timeout(job *minion.Job[*FileMatch]) time.Duration {
// 	return 60 * time.Minute
// }
// func (j *FileMatch) Work(ctx context.Context, job *minion.Job[*FileMatch]) error {
// 	l := app.Log.Named("files.match")
//
// 	for _, kind := range KINDS { // TODO: use libraries
// 		dir := filepath.Join(app.Config.DirectoriesCompleted, kind)
// 		if err := app.Workers.Enqueue(&FileMatchDir{Dir: dir}); err != nil {
// 			l.Errorw("enqueue", "error", err)
// 			return fae.Wrap(err, "enqueue")
// 		}
// 	}
//
// 	return nil
// }
//
// type FileMatchMedium struct {
// 	minion.WorkerDefaults[*FileMatchMedium]
// 	ID string
// }
//
// func (j *FileMatchMedium) Kind() string { return "file_match_medium" }
// func (j *FileMatchMedium) Timeout(job *minion.Job[*FileMatchMedium]) time.Duration {
// 	return 60 * time.Minute
// }
// func (j *FileMatchMedium) Work(ctx context.Context, job *minion.Job[*FileMatchMedium]) error {
// 	a := ContextApp(ctx)
// 	if a == nil {
// 		return fae.Errorf("no app context")
// 	}
// 	l := a.Log.Named("files.match.medium")
//
// 	m := &Medium{}
// 	if err := a.DB.Medium.Find(job.Args.ID, m); err != nil {
// 		l.Errorw("find", "error", err)
// 		return fae.Wrap(err, "finding")
// 	}
//
// 	dest := m.Destination()
// 	dir := filepath.Join(a.Config.DirectoriesCompleted, dest)
// 	if err := a.Workers.Enqueue(&FileMatchDir{Dir: dir}); err != nil {
// 		l.Errorw("enqueue", "error", err)
// 		return fae.Wrap(err, "enqueue")
// 	}
//
// 	return nil
// }
//
// type FileMatchDir struct {
// 	minion.WorkerDefaults[*FileMatchDir]
// 	Dir string
// }
//
// func (j *FileMatchDir) Kind() string { return "file_match_dir" }
// func (j *FileMatchDir) Timeout(job *minion.Job[*FileMatchDir]) time.Duration {
// 	return 60 * time.Minute
// }
// func (j *FileMatchDir) Work(ctx context.Context, job *minion.Job[*FileMatchDir]) error {
// 	a := ContextApp(ctx)
// 	if a == nil {
// 		return fae.Errorf("no app context")
// 	}
//
// 	dir := job.Args.Dir
// 	l := app.Log.Named("files.match.dir").With("dir", dir)
// 	l.Debugf("running")
//
// 	if !exists(dir) {
// 		notifier.Log.Warnf("files", "dir not found: %s", dir)
// 		return nil
// 	}
//
// 	return a.fileMatchDir(dir)
// }
//
// type FileMatchAnime struct {
// 	minion.WorkerDefaults[*FileMatchAnime]
// }
//
// func (j *FileMatchAnime) Kind() string { return "file_match_anime" }
// func (j *FileMatchAnime) Work(ctx context.Context, job *minion.Job[*FileMatchAnime]) error {
// 	a := ContextApp(ctx)
// 	if err := a.Workers.Enqueue(&FileMatchKind{Input: "anime"}); err != nil {
// 		return fae.Wrap(err, "enqueue")
// 	}
// 	return nil
// }
//
// type FileMatchDonghua struct {
// 	minion.WorkerDefaults[*FileMatchDonghua]
// }
//
// func (j *FileMatchDonghua) Kind() string { return "file_match_donghua" }
// func (j *FileMatchDonghua) Work(ctx context.Context, job *minion.Job[*FileMatchDonghua]) error {
// 	a := ContextApp(ctx)
// 	if err := a.Workers.Enqueue(&FileMatchKind{Input: "donghua"}); err != nil {
// 		return fae.Wrap(err, "enqueue")
// 	}
// 	return nil
// }
//
// type FileMatchKind struct {
// 	minion.WorkerDefaults[*FileMatchKind]
// 	Input string `bson:"input" json:"input"`
// }
//
// func (j *FileMatchKind) Kind() string { return "file_match_kind" }
// func (j *FileMatchKind) Work(ctx context.Context, job *minion.Job[*FileMatchKind]) error {
// 	kind := job.Args.Input
// 	a := ContextApp(ctx)
//
// 	err := a.DB.Medium.Query().Where("kind", kind).Batch(100, func(media []*Medium) error {
// 		for _, m := range media {
// 			if err := a.Workers.Enqueue(&FileMatchMedium{ID: m.ID.Hex()}); err != nil {
// 				return fae.Wrapf(err, "enqueue medium: %s", m.ID.Hex())
// 			}
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		return fae.Wrap(err, "querying")
// 	}
//
// 	return nil
// }

type FilesRename struct {
	minion.WorkerDefaults[*FilesRename]
}

func (j *FilesRename) Kind() string { return "files_rename" }
func (j *FilesRename) Work(ctx context.Context, job *minion.Job[*FilesRename]) error {
	a := ContextApp(ctx)
	// l := a.Workers.Log.Named("files_rename")
	id := "65a4943c175ec2916ae45688"

	if err := a.Workers.Enqueue(&FilesRenameMedium{ID: id}); err != nil {
		return fae.Wrap(err, "enqueue")
	}

	return nil
}

type FilesRenameMedium struct {
	minion.WorkerDefaults[*FilesRenameMedium]
	ID string `bson:"id" json:"id"`
}

func (j *FilesRenameMedium) Kind() string { return "files_rename_medium" }
func (j *FilesRenameMedium) Work(ctx context.Context, job *minion.Job[*FilesRenameMedium]) error {
	a := ContextApp(ctx)
	l := a.Workers.Log.Named("files_rename_medium")
	ID := job.Args.ID

	m, err := a.DB.Medium.Get(ID, &Medium{})
	if err != nil {
		return fae.Wrap(err, "getting medium")
	}
	if m == nil {
		return fae.Errorf("medium not found")
	}
	if m.Type != "Series" {
		return fae.Errorf("not a series")
	}

	q := a.DB.Medium.Query().Where("_type", "Episode").Where("series_id", m.ID)
	total, err := q.Count()
	if err != nil {
		return fae.Wrap(err, "counting")
	}

	l.Debugf("medium: %s:(%d) %s: %s", m.ID.Hex(), total, m.Title, m.Destination())

	err = q.Batch(100, func(episodes []*Medium) error {
		for _, e := range episodes {
			existingPaths := lo.Map(e.Paths, func(p *Path, _ int) string { return p.LocalPath() })
			newPaths := make([]*Path, 0)
			for _, p := range e.Paths {
				if p.IsCoverBackground() {
					continue
				}

				d, err := a.Destinator.Destination(m.Kind, e)
				if err != nil {
					return fae.Wrap(err, "destination")
				}
				dest := fmt.Sprintf("%s.%s", d, p.Extension)

				if p.LocalPath() != dest {
					a.Log.Warnw("rename", "from", p.Local, "to", dest)

					kind, name, file, ext, err := pathParts(dest)
					if err != nil {
						return fae.Wrap(err, "parts")
					}

					p.Old = true

					if err := FileLink(p.LocalPath(), dest, false); err != nil {
						l.Errorf("link: %s: %s", p.Local, err)
						continue
					}

					if lo.Contains(existingPaths, fmt.Sprintf("%s/%s/%s", kind, name, file)) {
						continue
					}

					np := &Path{
						Type:       p.Type,
						Local:      fmt.Sprintf("%s/%s/%s", kind, name, file),
						Extension:  ext,
						Size:       p.Size,
						UpdatedAt:  p.UpdatedAt,
						Resolution: p.Resolution,
						Bitrate:    p.Bitrate,
						Checksum:   p.Checksum,
					}
					newPaths = append(newPaths, np)
				}
			}

			e.Paths = append(e.Paths, newPaths...)

			if err := a.DB.Medium.Save(e); err != nil {
				return fae.Wrap(err, "saving episode")
			}
		}
		return nil
	})
	if err != nil {
		return fae.Wrap(err, "episode batch")
	}

	dir := filepath.Join(a.Config.DirectoriesCompleted, m.Destination())
	if err := a.Plex.RefreshLibraryPath(dir); err != nil {
		return fae.Wrap(err, "failed to refresh library")
	}

	if err := a.DB.Medium.Save(m); err != nil {
		return fae.Wrap(err, "saving series")
	}

	return nil
}
