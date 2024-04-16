package app

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/dashotv/fae"
	"github.com/dashotv/minion"
)

var downloadProcessMutex = &CtxMutex{ch: make(chan struct{}, 1)}
var downloadMultiFiles = 3

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

	// notifier.Info("Downloads", "processing downloads")
	funcs := []func() error{
		j.Create,
		j.Search,
		j.Load,
		j.Manage,
		j.Move,
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

func (j *DownloadsProcess) Create() error {
	seriesDownloads, err := app.DB.SeriesDownloadCounts()
	if err != nil {
		return fae.Wrap(err, "failed to get series download counts")
	}

	list, err := app.DB.UpcomingNow()
	if err != nil {
		return fae.Wrap(err, "failed to get upcoming episodes")
	}

	for _, ep := range list {
		//app.Workers.Log.Debugf("DownloadsProcess: create: %s %s", ep.Title, ep.Display)
		if !ep.SeriesActive {
			//app.Workers.Log.Debugf("DownloadsProcess: create: %s %s: not active", ep.Title, ep.Display)
			continue
		}

		unwatched, err := app.DB.SeriesUnwatchedByID(ep.SeriesID.Hex())
		if err != nil {
			return fae.Wrap(err, "failed to get unwatched")
		}

		if unwatched+seriesDownloads[ep.SeriesID.Hex()] >= 3 {
			continue
		}

		app.Workers.Log.Debugf("DownloadsProcess: create: %s - %s", ep.SeriesTitle, ep.Display)
		notifier.Info("Downloads::Create", fmt.Sprintf("%s - %s", ep.SeriesTitle, ep.Display))
		seriesDownloads[ep.SeriesID.Hex()]++

		d := &Download{
			Status:   "searching",
			MediumID: ep.ID,
			Auto:     true,
		}
		err = app.DB.Download.Save(d)
		if err != nil {
			return fae.Wrap(err, "failed to save download")
		}

		err = app.DB.EpisodeSetting(ep.ID.Hex(), "downloaded", true)
		if err != nil {
			return fae.Wrap(err, "failed to save episode")
		}
	}

	return nil
}

func (j *DownloadsProcess) Search() error {
	list, err := app.DB.DownloadByStatus("searching")
	if err != nil {
		return fae.Wrap(err, "failed to get downloads")
	}

	for _, d := range list {
		if d.Medium == nil {
			continue
		}
		if d.Medium.Type != "Episode" {
			//TODO: handle movies
			continue
		}

		match, err := app.ScrySearchEpisode(d.Search)
		if err != nil {
			return fae.Wrap(err, "failed to search releases")
		}
		if match == nil {
			continue
		}

		app.Workers.Log.Debugf("DownloadsProcess: found: %s - %s", d.Title, d.Display)
		notifier.Info("Downloads::Found", fmt.Sprintf("%s - %s", d.Title, d.Display))

		d.Status = "loading" // TODO: review
		if !app.Config.Production {
			d.Status = "reviewing"
		}
		d.URL = match.Download

		err = app.DB.Download.Save(d)
		if err != nil {
			return fae.Wrap(err, "failed to save download")
		}
	}
	return nil
}

func (j *DownloadsProcess) Load() error {
	list, err := app.DB.DownloadByStatus("loading")
	if err != nil {
		return fae.Wrap(err, "failed to get downloads")
	}

	for _, d := range list {
		if d.ReleaseID == "" && d.URL == "" {
			app.DB.Log.Debugf("DownloadsProcess: load: %s %s: no release", d.Title, d.Display)
			continue
		}

		res, err := app.FlameAdd(d)
		if err != nil {
			return fae.Wrap(err, "failed to add to flame")
		}

		d.Status = "downloading"
		if d.IsTorrent() {
			d.Status = "managing"
		}
		d.Thash = res

		err = app.DB.Download.Save(d)
		if err != nil {
			return fae.Wrap(err, "failed to save download")
		}
	}

	return nil
}

func (j *DownloadsProcess) Manage() error {
	list, err := app.DB.DownloadByStatus("managing")
	if err != nil {
		return fae.Wrap(err, "get downloads")
	}

	for _, d := range list {
		// TODO: manage metube? show files while downloading?
		if d.Thash == "" || !d.IsTorrent() {
			continue
		}

		t, err := app.FlameTorrent(d.Thash)
		if err != nil {
			app.Log.Named("downloads.manage").Errorf("failed to get torrent: %s", err)
			continue
		}

		if len(t.Files) == 0 {
			continue
		}

		dfs := d.Files
		numToDf := map[int]*DownloadFile{}
		for _, df := range dfs {
			numToDf[df.Num] = df
		}

		for _, f := range t.Files {
			if shouldDownloadFile(f.Name) {
				if _, ok := numToDf[f.ID]; !ok {
					// if it doesn't already exist, add it
					d.Files = append(d.Files, &DownloadFile{Num: f.ID})
				}
			}
		}

		if len(d.Files) == 0 {
			app.Workers.Log.Warnf("download has no files: %s", d.ID.Hex())
			continue
		}

		if len(d.Files) == 1 {
			d.Files[0].MediumID = d.MediumID
			d.Status = "downloading"

			err = app.DB.Download.Save(d)
			if err != nil {
				return fae.Wrap(err, "failed to save download")
			}

			continue
		}

		if !d.Multi {
			app.Workers.Log.Warnf("multiple files, but not multi", d.ID.Hex())

			d.Status = "reviewing"
			err = app.DB.Download.Save(d)
			if err != nil {
				return fae.Wrap(err, "failed to save download")
			}

			continue
		}

		for _, df := range d.Files {
			if df.MediumID != primitive.NilObjectID {
				// already has media
				continue
			}

			if d.Medium.Type != "Series" {
				// only handle series for now
				app.Workers.Log.Warnf("multi not series", d.ID.Hex())

				d.Status = "reviewing"
				err = app.DB.Download.Save(d)
				if err != nil {
					return fae.Wrap(err, "failed to save download")
				}
			}

			file := t.Files[df.Num]

			// find the episode based on the name
			ep, err := app.RunicFindEpisode(d.MediumID, file.Name, "tv")
			if err != nil {
				return fae.Wrap(err, "failed to find episode")
			}

			if ep == nil {
				app.Workers.Log.Warnf("episode not found: %s", file.Name)
				continue
			}

			df.MediumID = ep.ID
		}

		// TODO: handle want more / none / etc
		wanted := false
		for _, f := range t.Files {
			if f.Priority > 0 {
				wanted = true
				break
			}
		}

		if wanted && t.Progress != 100 {
			err := app.FlameTorrentWant(d.Thash, "none")
			if err != nil {
				return fae.Wrap(err, "want none")
			}
		}

		if d.HasMedia() {
			nums := d.NextFileNums(t, downloadMultiFiles)
			if nums != "" {
				err := app.FlameTorrentWant(d.Thash, nums)
				if err != nil {
					return fae.Wrap(err, "want next")
				}
			}

			// save updates to download files
			d.Status = "downloading"
		}

		if err := app.DB.Download.Save(d); err != nil {
			return fae.Wrap(err, "failed to save download")
		}
	}

	return nil
}

func (j *DownloadsProcess) Move() error {
	list, err := app.DB.DownloadByStatus("downloading")
	if err != nil {
		return fae.Wrap(err, "failed to get downloads")
	}

	moved := []string{}

	if len(list) == 0 {
		return nil
	}

	for _, d := range list {
		if d.Medium == nil || d.Thash == "" || d.IsNzb() {
			continue
		}

		if d.IsMetube() {
			files, err := DownloadMove(d)
			if err != nil {
				return fae.Wrap(err, "metube move")
			}

			moved = append(moved, files...)
			continue
		} else {
			t, err := app.FlameTorrent(d.Thash)
			if err != nil {
				app.Log.Debugf("error: %+v", err)
				return fae.Wrap(err, "getting torrent")
			}

			mover := NewMover(app.Log.Named("mover"), d, t)
			files, err := mover.Move()
			// files, err := DownloadMove(d)
			if err != nil {
				app.Log.Debugf("error: %+v", err)
				return fae.Wrap(err, "move download")
			}

			if files == nil || len(files) == 0 {
				continue
			}

			moved = append(moved, files...)
			// update medium and add path
			if err := updateMedium(d.Medium, files); err != nil {
				d.Status = "reviewing"
			}

			if d.Multi {
				nums := d.NextFileNums(t, 3)
				if nums != "" {
					err := app.FlameTorrentWant(d.Thash, nums)
					if err != nil {
						return fae.Wrap(err, "want next")
					}
				}

				continue
			}
		}

		notifier.Success("Downloads::Completed", fmt.Sprintf("%s - %s", d.Title, d.Display))
		if d.Status != "reviewing" {
			d.Status = "done"
		}

		err = app.DB.Download.Save(d)
		if err != nil {
			return fae.Wrap(err, "failed to save download")
		}

		if d.IsTorrent() {
			if err := app.FlameTorrentRemove(d.Thash); err != nil {
				return fae.Wrap(err, "failed to remove torrent")
			}
		}
	}

	if len(moved) > 0 {
		dirs := lo.Map(moved, func(s string, i int) string {
			return filepath.Dir(s)
		})
		dirs = lo.Uniq(dirs)

		for _, dir := range dirs {
			notifier.Log.Info("downloads: refresh: ", dir)
			err := app.Plex.RefreshLibraryPath(dir)
			if err != nil {
				return fae.Wrap(err, "failed to refresh library")
			}
		}
	}

	return nil
}

func updateMedium(m *Medium, files []string) error {
	m.Completed = true

	for _, f := range files {
		path := m.AddPathByFullpath(f)

		if err := app.Workers.Enqueue(&PathImport{ID: m.ID.Hex(), PathID: path.ID.Hex(), Title: path.Local}); err != nil {
			return fae.Errorf("enqueue path: %s", err)
		}
	}

	if err := app.DB.Medium.Save(m); err != nil {
		return fae.Wrap(err, "failed to save medium")
	}

	return nil
}

// TODO: handle context?
func DownloadMove(d *Download) ([]string, error) {
	l := app.Log.Named("download.move:" + d.ID.Hex())
	out := []string{}

	files, err := Files(d)
	if err != nil {
		return nil, fae.Wrap(err, "failed to get files")
	}
	if len(files) == 0 {
		return nil, nil
	}

	for _, source := range files {
		var dest string
		var err error
		ext := Extension(source)

		if d.Medium.Type == "Series" {
			dest, err = Destination(d.Medium)
		} else {
			dest, err = Destination(d.Medium)
		}
		if err != nil {
			return nil, fae.Wrap(err, "failed to get destination")
		}

		file := strings.ToLower(fmt.Sprintf("%s.%s", dest, ext))
		destination := fmt.Sprintf("%s/%s", app.Config.DirectoriesCompleted, file)

		l.Debugf("mover: %s", source)
		l.Debugf("    -> %s", destination)

		if !exists(source) {
			return nil, fae.Errorf("source does not exist: %s", source)
		}
		if exists(destination) {
			if !d.Force {
				notifier.Log.Warn("DownloadMove", fmt.Sprintf("destination exists, force false: %s", destination))
				return nil, nil
			}

			match, err := sumFiles(source, destination)
			if err != nil {
				return nil, fae.Errorf("failed to sum files")
			}
			if match {
				notifier.Log.Warn("DownloadMove", fmt.Sprintf("destination exists, sums match: %s", destination))
				return nil, nil
			}
		}

		if !app.Config.Production {
			l.Debugf("skipping move in dev mode")
			continue
		}

		if err := FileLink(source, destination, d.Force); err != nil {
			return nil, fae.Wrap(err, "link")
		}

		out = append(out, destination)
	}

	return out, nil
}

// func (j *DownloadsProcess) MetubeMove(download *Download) error {
// 	l := app.Log.Named("downloads.metube")
// 	files := []string{}
// 	moved := []string{}
//
// 	history, err := app.Flame.MetubeHistory()
// 	if err != nil {
// 		return fae.Wrap(err, "history")
// 	}
//
// 	done, ok := lo.Find(history.Done, func(h *metube.Download) bool {
// 		return h.CustomNamePrefix == download.ID.Hex()
// 	})
// 	if !ok || done == nil {
// 		l.Debugf("not done: %s", download.ID.Hex())
// 		return nil
// 	}
//
// 	if download.Medium == nil || download.Medium.Type != "Episode" {
// 		l.Debugf("not episode: %s", download.ID.Hex())
// 		return nil
// 	}
//
// 	err = filepath.WalkDir(app.Config.DirectoriesMetube, func(path string, d fs.DirEntry, err error) error {
// 		if err != nil {
// 			return err
// 		}
//
// 		if d.IsDir() {
// 			return nil
// 		}
//
// 		if strings.Contains(path, download.ID.Hex()) {
// 			files = append(files, path)
// 		}
//
// 		return nil
// 	})
// 	if err != nil {
// 		return fae.Wrap(err, "walk")
// 	}
//
// 	files = lo.Filter(files, func(s string, i int) bool {
// 		return shouldDownloadFile(s)
// 	})
//
// 	if len(files) == 0 {
// 		return nil
// 	}
//
// 	l.Debugf("%s: files: %d", download.ID.Hex(), len(files))
// 	for _, f := range files {
// 		ext := filepath.Ext(f)
// 		if ext[0] == '.' {
// 			ext = ext[1:]
// 		}
//
// 		dest, err := Destination(download.Medium)
// 		if err != nil {
// 			return fae.Wrap(err, "destination")
// 		}
//
// 		destination := fmt.Sprintf("%s/%s.%s", app.Config.DirectoriesCompleted, dest, ext)
// 		l.Debugf("move:  %s", filepath.Base(f))
// 		l.Debugf("    -> %s", destination)
//
// 		if exists(destination) && !download.Force {
// 			return fae.New("exists, force false")
// 		}
//
// 		if !app.Config.Production {
// 			l.Debugf("dev mode")
// 			continue
// 		}
//
// 		if err := FileLink(f, destination, download.Force); err != nil {
// 			l.Errorf("copy: %s", err)
// 			return fae.Wrap(err, "copy")
// 		}
//
// 		download.Status = "done"
//
// 		// update medium and add path
// 		if err := updateMedium(download.Medium.ID.Hex(), dest, ext); err != nil {
// 			download.Status = "reviewing"
// 		}
//
// 		err = app.DB.Download.Save(download)
// 		if err != nil {
// 			return fae.Wrap(err, "failed to save download")
// 		}
//
// 		moved = append(moved, destination)
// 		notifier.Success("Downloads::Completed", fmt.Sprintf("%s %s", download.Medium.Title, download.Medium.Display))
// 	}
//
// 	if len(moved) > 0 {
// 		dirs := lo.Map(moved, func(s string, i int) string {
// 			return filepath.Dir(s)
// 		})
// 		dirs = lo.Uniq(dirs)
//
// 		for _, dir := range dirs {
// 			err := app.Plex.RefreshLibraryPath(dir)
// 			if err != nil {
// 				return fae.Wrap(err, "failed to refresh library")
// 			}
// 		}
// 	}
//
// 	return nil
// }

//
// type DownloadFileMover struct {
// 	ID string
// }
//
// func (j *DownloadFileMover) Kind() string { return "DownloadsFileMove" }
// func (j *DownloadFileMover) Work(ctx context.Context, job *minion.Job[*DownloadFileMover]) error {
// 	d, err :=app.DB.DownloadGet(job.Args.ID)
// 	if err != nil {
// 		return fae.Wrap(err, "failed to get download")
// 	}
// 	notifier.Info("Downloads::Move", fmt.Sprintf("%s %s", d.Medium.Title, d.Medium.Display))
//
// 	t, err := flameClient.Torrent(d.Thash)
// 	if err != nil {
// 		app.Log.Named("downloads.manage").Errorf("failed to get torrent: %s", err)
// 		continue
// 	}
//
// 	tf := t.Files[0]
// 	// df := d.Files[0]
// 	kind := d.Medium.Kind
// 	dir := d.Medium.Directory
// 	ext := filepath.Ext(tf.Name)
// 	if len(ext) > 0 {
// 		ext = ext[1:]
// 	}
//
// 	source := fmt.Sprintf("%s/%s", cfg.Directories.Incoming, tf.Name)
// 	file := strings.ToLower(fmt.Sprintf("%s/%s/%s %s.%s", kind, dir, dir, d.Medium.Display, ext))
// 	destination := fmt.Sprintf("%s/%s", cfg.Directories.Completed, file)
//
//app.Workers.Log.Debugf("mover: %s", source)
//app.Workers.Log.Debugf("    -> %s", destination)
//
// 	if !exists(source) {
// 		return fae.New("source does not exist")
// 	}
// 	if exists(destination) {
// 		if !d.Force {
// 			return fae.New("destination exists, force false")
// 		}
//
// 		match, err := sumFiles(source, destination)
// 		if err != nil {
// 			return fae.Wrap(err, "failed to sum files")
// 		}
// 		if match {
//app.Workers.Log.Debugf("destination exists, checksums match")
// 			notifier.Log.Info("Downloads::FileMover", fmt.Sprintf("destination exists, checksums match: %s %s", d.Medium.Title, d.Medium.Display))
// 			return nil
// 		}
// 	}
//
// 	if err := FileCopy(source, destination); err != nil {
// 		return fae.Wrap(err, "copy")
// 	}
//
// 	return nil
// }
