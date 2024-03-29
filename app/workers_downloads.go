package app

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/dashotv/flame/metube"
	"github.com/dashotv/minion"
)

type DownloadsProcess struct {
	minion.WorkerDefaults[*DownloadsProcess]
}

func (j *DownloadsProcess) Kind() string { return "DownloadsProcess" }
func (j *DownloadsProcess) Work(ctx context.Context, job *minion.Job[*DownloadsProcess]) error {
	// notifier.Info("Downloads", "processing downloads")
	funcs := []func() error{
		// j.Create,
		// j.Search,
		j.Load,
		j.Manage,
		j.Move,
	}

	for _, f := range funcs {
		err := f()
		if err != nil {
			return err
		}
	}

	return nil
}

func (j *DownloadsProcess) Create() error {
	list, err := app.DB.UpcomingNow()
	if err != nil {
		return errors.Wrap(err, "failed to get upcoming episodes")
	}

	seriesDownloads, err := app.DB.SeriesDownloadCounts()
	if err != nil {
		return errors.Wrap(err, "failed to get series download counts")
	}

	for _, ep := range list {
		//app.Workers.Log.Debugf("DownloadsProcess: create: %s %s", ep.Title, ep.Display)
		if !ep.Active {
			//app.Workers.Log.Debugf("DownloadsProcess: create: %s %s: not active", ep.Title, ep.Display)
			continue
		}

		if seriesDownloads[ep.SeriesId.Hex()] >= 3 {
			//app.Workers.Log.Debugf("DownloadsProcess: create: %s %s: series downloads", ep.Title, ep.Display)
			continue
		}

		if !ep.Favorite && ep.Unwatched >= 3 {
			// If I'm not watching it, see if others are
			unwatched, err := app.DB.SeriesUnwatchedByID(ep.SeriesId.Hex())
			if err != nil {
				return errors.Wrap(err, "failed to get unwatched")
			}

			if unwatched >= 3 {
				//app.Workers.Log.Debugf("DownloadsProcess: create: %s %s: unwatched >3", ep.Title, ep.Display)
				continue
			}
		}

		app.Workers.Log.Debugf("DownloadsProcess: create: %s %s", ep.Title, ep.Display)
		notifier.Info("Downloads::Create", fmt.Sprintf("%s %s", ep.Title, ep.Display))
		seriesDownloads[ep.SeriesId.Hex()]++

		d := &Download{
			Status:   "searching",
			MediumId: ep.ID,
			Auto:     true,
		}
		err = app.DB.Download.Save(d)
		if err != nil {
			return errors.Wrap(err, "failed to save download")
		}

		err = app.DB.EpisodeSetting(ep.ID.Hex(), "downloaded", true)
		if err != nil {
			return errors.Wrap(err, "failed to save episode")
		}
	}

	return nil
}

func (j *DownloadsProcess) Search() error {
	list, err := app.DB.DownloadByStatus("searching")
	if err != nil {
		return errors.Wrap(err, "failed to get downloads")
	}

	for _, d := range list {
		if d.Medium == nil {
			continue
		}
		if d.Medium.Type != "Episode" {
			//TODO: handle movies
			continue
		}

		//app.Workers.Log.Debugf("DownloadsProcess: search: %s %s", d.Medium.Title, d.Medium.Display)
		match, err := app.Scry.ScrySearchEpisode(d.Medium)
		if err != nil {
			return errors.Wrap(err, "failed to search releases")
		}
		if match == nil {
			continue
		}

		notifier.Info("Downloads::Found", fmt.Sprintf("%s %s", d.Medium.Title, d.Medium.Display))

		d.Status = "loading"
		if match.NZB {
			d.Url = match.Download
		} else {
			d.ReleaseId = match.ID
		}

		err = app.DB.Download.Save(d)
		if err != nil {
			return errors.Wrap(err, "failed to save download")
		}
	}
	return nil
}

func (j *DownloadsProcess) Load() error {
	list, err := app.DB.DownloadByStatus("loading")
	if err != nil {
		return errors.Wrap(err, "failed to get downloads")
	}

	for _, d := range list {
		if d.ReleaseId == "" && d.Url == "" {
			app.DB.Log.Debugf("DownloadsProcess: load: %s %s: no release", d.Medium.Title, d.Medium.Display)
			continue
		}

		url, err := d.GetURL()
		if err != nil {
			return errors.Wrap(err, "failed to get url")
		}

		if nzbgeekRegex.MatchString(url) {
			id, err := app.Flame.LoadNzb(d, url)
			if err != nil {
				return errors.Wrap(err, "failed to load nzb")
			}
			d.Status = "downloading"
			d.Thash = id
		} else if metubeRegex.MatchString(url) {
			autoStart := false
			if app.Config.Production {
				autoStart = true
			}
			app.Log.Named("downloads").Debugf("loading metube: %s", url)
			url = strings.Replace(url, "metube://", "", 1)
			err := app.Flame.LoadMetube(d.ID.Hex(), url, autoStart)
			if err != nil {
				return fmt.Errorf("failed to load metube: %w", err)
			}
			d.Status = "downloading"
			d.Thash = "M"
		} else {
			thash, err := app.Flame.LoadTorrent(d, url)
			if err != nil {
				return errors.Wrap(err, "failed to load torrent")
			}
			d.Status = "managing"
			d.Thash = strings.ToLower(thash)
		}

		err = app.DB.Download.Save(d)
		if err != nil {
			return errors.Wrap(err, "failed to save download")
		}
	}

	return nil
}

func (j *DownloadsProcess) Manage() error {
	list, err := app.DB.DownloadByStatus("managing")
	if err != nil {
		return fmt.Errorf("failed to get downloads: %w", err)
	}

	for _, d := range list {
		if d.Thash == "" {
			continue
		}
		if d.Thash == "M" {
			continue
		}

		if d.IsNzb() {
			continue
		}

		t, err := app.Flame.Torrent(d.Thash)
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

		if len(d.Files) > 1 {
			app.Workers.Log.Warnf("download has multiple files: %s", d.ID.Hex())
			continue
		}

		d.Files[0].MediumId = d.MediumId
		d.Status = "downloading"

		err = app.DB.Download.Save(d)
		if err != nil {
			return errors.Wrap(err, "failed to save download")
		}
	}

	return nil
}

func (j *DownloadsProcess) Move() error {
	l := app.Log.Named("downloads.move")
	list, err := app.DB.DownloadByStatus("downloading")
	if err != nil {
		return errors.Wrap(err, "failed to get downloads")
	}

	moved := []string{}

	for _, d := range list {
		if d.Thash == "" {
			continue
		}

		if d.IsNzb() {
			continue
		}
		if d.IsMetube() {
			j.MetubeMove(d)
			continue
		}

		t, err := app.Flame.Torrent(d.Thash)
		if err != nil {
			l.Errorf("failed to get torrent: %s", err)
			continue
		}

		if t.Progress < 100 {
			continue
		}

		if len(t.Files) == 0 {
			continue
		}

		if len(t.Files) > 1 {
			continue
		}

		// notifier.Info("Downloads::Move", fmt.Sprintf("%s %s", d.Medium.Title, d.Medium.Display))

		tf := t.Files[0]
		// df := d.Files[0]
		ext := filepath.Ext(tf.Name)
		if len(ext) > 0 {
			ext = ext[1:]
		}

		// TODO: move to configurable templates
		dest, err := Destination(d.Medium)
		if err != nil {
			return errors.Wrap(err, "failed to get destination")
		}

		source := fmt.Sprintf("%s/%s", app.Config.DirectoriesIncoming, tf.Name)
		file := strings.ToLower(fmt.Sprintf("%s.%s", dest, ext))
		destination := fmt.Sprintf("%s/%s", app.Config.DirectoriesCompleted, file)

		l.Debugf("mover: %s", source)
		l.Debugf("    -> %s", destination)

		if !app.Config.Production {
			l.Debugf("skipping move in dev mode")
			continue
		}

		if !exists(source) {
			return errors.Errorf("source does not exist: %s", source)
		}
		if exists(destination) {
			if !d.Force {
				return errors.New("destination exists, force false")
			}

			match, err := sumFiles(source, destination)
			if err != nil {
				return errors.Wrap(err, "failed to sum files")
			}
			if match {
				l.Debugf("destination exists, checksums match")
				notifier.Log.Info("Downloads::FileMover", fmt.Sprintf("destination exists, checksums match: %s %s", d.Medium.Title, d.Medium.Display))
				return nil
			}
		}

		if err := FileLink(source, destination, d.Force); err != nil {
			return errors.Wrap(err, "copy")
		}

		d.Status = "done"

		// update medium and add path
		if err := updateMedium(d.Medium.ID.Hex(), dest, ext); err != nil {
			d.Status = "reviewing"
		}

		err = app.DB.Download.Save(d)
		if err != nil {
			return errors.Wrap(err, "failed to save download")
		}

		if err := app.Flame.RemoveTorrent(d.Thash); err != nil {
			return errors.Wrap(err, "failed to remove torrent")
		}

		moved = append(moved, destination)
		notifier.Success("Downloads::Completed", fmt.Sprintf("%s %s", d.Medium.Title, d.Medium.Display))
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
				return errors.Wrap(err, "failed to refresh library")
			}
		}
	}

	return nil
}

func updateMedium(id, dest, ext string) error {
	m, err := app.DB.Medium.Get(id, &Medium{})
	if err != nil {
		return fmt.Errorf("get medium: %w", err)
	}

	m.Completed = true
	m.Paths = append(m.Paths, &Path{
		Local:     dest,
		Extension: ext[1:],
		Type:      primitive.Symbol(fileType(fmt.Sprintf("%s.%s", dest, ext[1:]))),
	})
	err = app.DB.Medium.Update(m)
	if err != nil {
		return fmt.Errorf("update medium: %s", err)
	}
	// enqueue path import
	path, ok := lo.Find(m.Paths, func(p *Path) bool {
		return p.Local == dest
	})
	if ok && path != nil {
		if err := app.Workers.Enqueue(&PathImport{ID: m.ID.Hex(), PathID: path.Id.Hex(), Title: path.Local}); err != nil {
			return fmt.Errorf("enqueue path: %s", err)
		}
	}

	return nil
}

func (j *DownloadsProcess) MetubeMove(download *Download) error {
	l := app.Log.Named("downloads.metube")
	files := []string{}
	moved := []string{}

	history, err := app.Flame.MetubeHistory()
	if err != nil {
		return fmt.Errorf("history: %w", err)
	}

	done, ok := lo.Find(history.Done, func(h *metube.Download) bool {
		return h.CustomNamePrefix == download.ID.Hex()
	})
	if !ok || done == nil {
		l.Debugf("not done: %s", download.ID.Hex())
		return nil
	}

	if download.Medium == nil || download.Medium.Type != "Episode" {
		l.Debugf("not episode: %s", download.ID.Hex())
		return nil
	}

	err = filepath.WalkDir(app.Config.DirectoriesMetube, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if strings.Contains(path, download.ID.Hex()) {
			files = append(files, path)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("walk: %w", err)
	}

	files = lo.Filter(files, func(s string, i int) bool {
		return shouldDownloadFile(s)
	})

	if len(files) == 0 {
		return nil
	}

	for _, f := range files {
		ext := filepath.Ext(f)

		dest, err := Destination(download.Medium)
		if err != nil {
			return fmt.Errorf("destination: %w", err)
		}

		destination := fmt.Sprintf("%s/%s.%s", app.Config.DirectoriesCompleted, dest, ext[1:])
		l.Debugf("move:  %s", filepath.Base(f))
		l.Debugf("    -> %s", destination)

		if exists(destination) && !download.Force {
			return errors.New("exists, force false")
		}

		if !app.Config.Production {
			l.Debugf("dev mode")
			continue
		}

		if err := FileLink(f, destination, download.Force); err != nil {
			l.Errorf("copy: %s", err)
			return fmt.Errorf("copy: %w", err)
		}

		download.Status = "done"

		// update medium and add path
		if err := updateMedium(download.Medium.ID.Hex(), dest, ext); err != nil {
			download.Status = "reviewing"
		}

		err = app.DB.Download.Save(download)
		if err != nil {
			return errors.Wrap(err, "failed to save download")
		}

		moved = append(moved, destination)
		notifier.Success("Downloads::Completed", fmt.Sprintf("%s %s", download.Medium.Title, download.Medium.Display))
	}

	if len(moved) > 0 {
		dirs := lo.Map(moved, func(s string, i int) string {
			return filepath.Dir(s)
		})
		dirs = lo.Uniq(dirs)

		for _, dir := range dirs {
			err := app.Plex.RefreshLibraryPath(dir)
			if err != nil {
				return errors.Wrap(err, "failed to refresh library")
			}
		}
	}

	return nil
}

//
// type DownloadFileMover struct {
// 	ID string
// }
//
// func (j *DownloadFileMover) Kind() string { return "DownloadsFileMove" }
// func (j *DownloadFileMover) Work(ctx context.Context, job *minion.Job[*DownloadFileMover]) error {
// 	d, err :=app.DB.DownloadGet(job.Args.ID)
// 	if err != nil {
// 		return errors.Wrap(err, "failed to get download")
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
// 		return errors.New("source does not exist")
// 	}
// 	if exists(destination) {
// 		if !d.Force {
// 			return errors.New("destination exists, force false")
// 		}
//
// 		match, err := sumFiles(source, destination)
// 		if err != nil {
// 			return errors.Wrap(err, "failed to sum files")
// 		}
// 		if match {
//app.Workers.Log.Debugf("destination exists, checksums match")
// 			notifier.Log.Info("Downloads::FileMover", fmt.Sprintf("destination exists, checksums match: %s %s", d.Medium.Title, d.Medium.Display))
// 			return nil
// 		}
// 	}
//
// 	if err := FileCopy(source, destination); err != nil {
// 		return errors.Wrap(err, "copy")
// 	}
//
// 	return nil
// }
