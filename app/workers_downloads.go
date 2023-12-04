package app

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/dashotv/minion"
)

type DownloadsProcess struct{}

func (j *DownloadsProcess) Kind() string { return "DownloadsProcess" }
func (j *DownloadsProcess) Work(ctx context.Context, job *minion.Job[*DownloadsProcess]) error {
	workers.Log.Debugf("DownloadsProcess: %s", job.ID)
	notifier.Info("Downloads", "processing downloads")
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
			return err
		}
	}

	return nil
}

func (j *DownloadsProcess) Create() error {
	list, err := db.UpcomingNow()
	if err != nil {
		return errors.Wrap(err, "failed to get upcoming episodes")
	}

	seriesDownloads, err := db.SeriesDownloadCounts()
	if err != nil {
		return errors.Wrap(err, "failed to get series download counts")
	}

	for _, ep := range list {
		// workers.Log.Debugf("DownloadsProcess: create: %s %s", ep.Title, ep.Display)
		if !ep.Active {
			// workers.Log.Debugf("DownloadsProcess: create: %s %s: not active", ep.Title, ep.Display)
			continue
		}

		if seriesDownloads[ep.SeriesId.Hex()] >= 3 {
			// workers.Log.Debugf("DownloadsProcess: create: %s %s: series downloads", ep.Title, ep.Display)
			continue
		}

		if !ep.Favorite && ep.Unwatched >= 3 {
			// If I'm not watching it, see if others are
			unwatched, err := db.SeriesUnwatchedByID(ep.SeriesId.Hex())
			if err != nil {
				return errors.Wrap(err, "failed to get unwatched")
			}

			if unwatched >= 3 {
				// workers.Log.Debugf("DownloadsProcess: create: %s %s: unwatched >3", ep.Title, ep.Display)
				continue
			}
		}

		workers.Log.Debugf("DownloadsProcess: create: %s %s", ep.Title, ep.Display)
		notifier.Info("Downloads::Create", fmt.Sprintf("%s %s", ep.Title, ep.Display))
		seriesDownloads[ep.SeriesId.Hex()]++

		d := &Download{
			Status:   "searching",
			MediumId: ep.ID,
			Auto:     true,
		}
		err = db.Download.Save(d)
		if err != nil {
			return errors.Wrap(err, "failed to save download")
		}

		err = db.EpisodeSetting(ep.ID.Hex(), "downloaded", true)
		if err != nil {
			return errors.Wrap(err, "failed to save episode")
		}
	}

	return nil
}

func (j *DownloadsProcess) Search() error {
	list, err := db.DownloadByStatus("searching")
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

		// workers.Log.Debugf("DownloadsProcess: search: %s %s", d.Medium.Title, d.Medium.Display)
		match, err := scryClient.ScrySearchEpisode(d.Medium)
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

		err = db.Download.Save(d)
		if err != nil {
			return errors.Wrap(err, "failed to save download")
		}
	}
	return nil
}

func (j *DownloadsProcess) Load() error {
	list, err := db.DownloadByStatus("loading")
	if err != nil {
		return errors.Wrap(err, "failed to get downloads")
	}

	for _, d := range list {
		if d.ReleaseId == "" && d.Url == "" {
			db.log.Debugf("DownloadsProcess: load: %s %s: no release", d.Medium.Title, d.Medium.Display)
			continue
		}

		url, err := d.GetURL()
		if err != nil {
			return errors.Wrap(err, "failed to get url")
		}

		if nzbgeekRegex.MatchString(url) {
			id, err := flameClient.LoadNzb(d, url)
			if err != nil {
				return errors.Wrap(err, "failed to load nzb")
			}
			d.Status = "downloading"
			d.Thash = id
		} else {
			thash, err := flameClient.LoadTorrent(d, url)
			if err != nil {
				return errors.Wrap(err, "failed to load torrent")
			}
			d.Status = "managing"
			d.Thash = strings.ToLower(thash)
		}

		err = db.Download.Save(d)
		if err != nil {
			return errors.Wrap(err, "failed to save download")
		}
	}

	return nil
}

func (j *DownloadsProcess) Manage() error {
	list, err := db.DownloadByStatus("managing")
	if err != nil {
		return errors.Wrap(err, "failed to get downloads")
	}

	for _, d := range list {
		if d.Thash == "" {
			continue
		}

		if d.IsNzb() {
			continue
		}

		t, err := flameClient.Torrent(d.Thash)
		if err != nil {
			return errors.Wrap(err, "failed to get torrent")
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
			workers.Log.Warnf("download has no files: %s", d.ID.Hex())
			continue
		}

		if len(d.Files) > 1 {
			workers.Log.Warnf("download has multiple files: %s", d.ID.Hex())
			continue
		}

		d.Files[0].MediumId = d.MediumId
		d.Status = "downloading"

		err = db.Download.Save(d)
		if err != nil {
			return errors.Wrap(err, "failed to save download")
		}
	}

	return nil
}

func (j *DownloadsProcess) Move() error {
	list, err := db.DownloadByStatus("downloading")
	if err != nil {
		return errors.Wrap(err, "failed to get downloads")
	}

	for _, d := range list {
		if d.Thash == "" {
			continue
		}

		if d.IsNzb() {
			continue
		}

		t, err := flameClient.Torrent(d.Thash)
		if err != nil {
			return errors.Wrap(err, "failed to get torrent")
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

		notifier.Info("Downloads::Move", fmt.Sprintf("%s %s", d.Medium.Title, d.Medium.Display))

		tf := t.Files[0]
		// df := d.Files[0]
		kind := d.Medium.Kind
		dir := d.Medium.Directory
		ext := filepath.Ext(tf.Name)
		if len(ext) > 0 {
			ext = ext[1:]
		}

		source := fmt.Sprintf("%s/%s", cfg.DirectoriesIncoming, tf.Name)
		file := strings.ToLower(fmt.Sprintf("%s/%s/%s %s.%s", kind, dir, dir, d.Medium.Display, ext))
		destination := fmt.Sprintf("%s/%s", cfg.DirectoriesCompleted, file)

		workers.Log.Debugf("mover: %s", source)
		workers.Log.Debugf("    -> %s", destination)

		if !exists(source) {
			return errors.New("source does not exist")
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
				workers.Log.Debugf("destination exists, checksums match")
				notifier.Log.Info("Downloads::FileMover", fmt.Sprintf("destination exists, checksums match: %s %s", d.Medium.Title, d.Medium.Display))
				return nil
			}
		}

		if err := FileCopy(source, destination); err != nil {
			return errors.Wrap(err, "copy")
		}

		err = db.EpisodeSetting(d.MediumId.Hex(), "completed", true)
		if err != nil {
			return errors.Wrap(err, "failed to save episode")
		}

		d.Status = "done"
		err = db.Download.Save(d)
		if err != nil {
			return errors.Wrap(err, "failed to save download")
		}

		notifier.Success("Downloads::Completed", fmt.Sprintf("%s %s", d.Medium.Title, d.Medium.Display))
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
// 	d, err := db.DownloadGet(job.Args.ID)
// 	if err != nil {
// 		return errors.Wrap(err, "failed to get download")
// 	}
// 	notifier.Info("Downloads::Move", fmt.Sprintf("%s %s", d.Medium.Title, d.Medium.Display))
//
// 	t, err := flameClient.Torrent(d.Thash)
// 	if err != nil {
// 		return errors.Wrap(err, "failed to get torrent")
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
// 	workers.Log.Debugf("mover: %s", source)
// 	workers.Log.Debugf("    -> %s", destination)
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
// 			workers.Log.Debugf("destination exists, checksums match")
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
