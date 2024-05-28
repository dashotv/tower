package app

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/dashotv/fae"
	"github.com/dashotv/flame/qbt"
)

var titleRegex = regexp.MustCompile(`(?i)^(?:episode|chapter)`)
var downloadMultiFiles = 3

func Extension(path string) string {
	ext := filepath.Ext(path)
	if len(ext) > 0 && ext[0] == '.' {
		ext = ext[1:]
	}
	return ext
}

// func Destination(m *Medium) (string, error) {
// 	switch m.Type {
// 	case "Series", "Movie":
// 		return fmt.Sprintf("%s/%s/%s", m.Kind, m.Directory, m.Directory), nil
// 	case "Episode":
// 		return destinationEpisode(m)
// 	default:
// 		return "", fae.Errorf("unknown type: %s", m.Type)
// 	}
// }
//
// func destinationEpisode(m *Medium) (string, error) {
// 	s := &Series{}
// 	err := app.DB.Series.FindByID(m.SeriesID, s)
// 	if err != nil {
// 		return "", err
// 	}
//
// 	e := &Episode{}
// 	err = app.DB.Episode.FindByID(m.ID, e)
// 	if err != nil {
// 		return "", err
// 	}
//
// 	out := ""
// 	if isAnimeKind(string(s.Kind)) && m.AbsoluteNumber > 0 {
// 		out = fmt.Sprintf("%s/%s/%s - %02dx%02d #%03d", s.Kind, s.Directory, s.Directory, m.SeasonNumber, m.EpisodeNumber, m.AbsoluteNumber)
// 	} else {
// 		out = fmt.Sprintf("%s/%s/%s - %02dx%02d", s.Kind, s.Directory, s.Directory, m.SeasonNumber, m.EpisodeNumber)
// 	}
//
// 	if e.Title != "" && !titleRegex.MatchString(e.Title) {
// 		out = fmt.Sprintf("%s - %s", out, path(e.Title))
// 	}
//
// 	return out, nil
// }

func updateMedium(m *Medium, files []*MoverFile) error {
	// fmt.Printf("updateMedium: %s %+v\n", m.ID.Hex(), files)
	// only mark downloaded/completed if there are videos (so subtitles don't trigger it)
	videos := lo.Filter(files, func(f *MoverFile, i int) bool {
		return fileType(f.Destination) == "video"
	})
	if len(videos) > 0 {
		m.Downloaded = true
		m.Completed = true
	}

	for _, f := range files {
		path := m.AddPathByFullpath(f.Destination)

		if err := app.Workers.Enqueue(&PathImport{ID: m.ID.Hex(), PathID: path.ID.Hex(), Title: f.Destination}); err != nil {
			return fae.Errorf("enqueue path: %s", err)
		}
	}

	if err := app.DB.Medium.Save(m); err != nil {
		return fae.Wrap(err, "failed to save medium")
	}

	return nil
}

func updateMedia(files []*MoverFile) error {
	for _, f := range files {
		if err := updateMedium(f.Medium, []*MoverFile{f}); err != nil {
			return fae.Wrapf(err, "update medium: %s", f.Destination)
		}
	}

	return nil
}

func (a *Application) downloadsCreate() error {
	// defer TickTock("DownloadsProcess: Create")()
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

		app.Workers.Log.Debugf("download created %s - %s", ep.SeriesTitle, ep.Display)
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

func (a *Application) downloadsSearch() error {
	// defer TickTock("DownloadsProcess: Search")()
	list, err := app.DB.DownloadByStatus("searching")
	if err != nil {
		return fae.Wrap(err, "failed to get downloads")
	}

	for _, d := range list {
		if d.Medium == nil {
			continue
		}
		if d.Medium.Type != "Episode" {
			//TODO: #4 handle movies
			continue
		}

		match, err := app.ScrySearchEpisode(d.Search)
		if err != nil {
			return fae.Wrap(err, "failed to search releases")
		}
		if match == nil {
			continue
		}

		app.Workers.Log.Debugf("download found %s - %s", d.Title, d.Display)
		notifier.Info("Downloads::Found", fmt.Sprintf("%s - %s", d.Title, d.Display))

		d.Status = "loading"
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

func (a *Application) downloadsLoad() (err error) {
	// defer TickTock("DownloadsProcess: Load")()
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
			return d.Error(fae.Wrap(err, "failed to add to flame"))
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

func (a *Application) downloadsManage() error {
	// defer TickTock("DownloadsProcess: Manage")()
	list, err := app.DB.DownloadByStatus("managing")
	if err != nil {
		return fae.Wrap(err, "get downloads")
	}

	for _, d := range list {
		// TODO: manage metube? show files while downloading?
		if d.Thash == "" || !d.IsTorrent() {
			continue
		}

		if d.Medium == nil {
			app.Workers.Log.Warnf("no medium", d.ID.Hex())
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

		if err := a.downloadsManageOne(d, t); err != nil {
			app.Log.Errorf("failed to manage download: %s", err)
		}
	}

	return nil
}

func (a *Application) downloadsManageOne(d *Download, t *qbt.Torrent) error {
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
		return nil
	}

	// TODO: handle downloads with single media file and multiple subtitles

	if len(d.Files) == 1 {
		d.Files[0].MediumID = d.MediumID
		d.Status = "downloading"

		if err := app.DB.Download.Save(d); err != nil {
			return fae.Wrap(err, "failed to save download")
		}

		return nil
	}

	if !d.Multi {
		app.Workers.Log.Warnf("multiple files, but not multi", d.ID.Hex())

		d.Status = "reviewing"
		if err := app.DB.Download.Save(d); err != nil {
			return fae.Wrap(err, "failed to save download")
		}

		return nil
	}

	if d.Medium.Type != "Series" {
		// only handle series for now
		app.Workers.Log.Warnf("multi not series", d.ID.Hex())

		d.Status = "reviewing"
		if err := app.DB.Download.Save(d); err != nil {
			return fae.Wrap(err, "failed to save download")
		}

		return nil
	}

	for _, df := range d.Files {
		if df.MediumID != primitive.NilObjectID {
			// already has media
			continue
		}

		file := t.Files[df.Num]

		// find the episode based on the name
		title := filepath.Base(file.Name)
		app.Log.Debugf("searching for episode: %s %s", d.Search.Type, title)
		ep, err := app.RunicFindEpisode(d.MediumID, title, d.Search.Type)
		if err != nil {
			return fae.Wrap(err, "failed to find episode")
		}

		if ep == nil {
			app.Workers.Log.Warnf("episode not found: %s", file.Name)
			continue
		}
		app.Log.Debugf("found: %s", ep.Title)

		df.MediumID = ep.ID
	}

	app.DB.processDownloads([]*Download{d})

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

	return nil
}

func (a *Application) downloadsMove() error {
	// defer TickTock("DownloadsProcess: Move")()
	list, err := app.DB.DownloadByStatus("downloading")
	if err != nil {
		return fae.Wrap(err, "failed to get downloads")
	}

	moved := []*MoverFile{}

	if len(list) == 0 {
		return nil
	}

	for _, d := range list {
		if d.Medium == nil || d.Thash == "" || d.IsNzb() {
			continue
		}

		var t *qbt.Torrent
		var err error
		if d.IsTorrent() {
			t, err = app.FlameTorrent(d.Thash)
			if err != nil {
				app.Log.Debugf("error: %+v", err)
				return fae.Wrap(err, "getting torrent")
			}
		}

		mover := NewMover(app.Log.Named("mover"), d, t)
		files, err := mover.Move()
		if err != nil {
			app.Log.Debugf("error: %+v", err)
			return d.Error(fae.Wrap(err, "move download"))
		}

		if d.Multi && d.Medium.Type == "Series" {
			// update medium and add path
			if files != nil && len(files) > 0 {
				moved = append(moved, files...)
				if err := updateMedia(files); err != nil {
					return d.Error(fae.Wrap(err, "update medium"))
				}
			}

			nums := d.NextFileNums(t, 3)
			if nums != "" {
				err := app.FlameTorrentWant(d.Thash, nums)
				if err != nil {
					return d.Error(fae.Wrap(err, "want next"))
				}
			}

			continue
		}

		if files == nil || len(files) == 0 {
			continue
		}

		moved = append(moved, files...)

		// update medium and add path
		if err := updateMedia(files); err != nil {
			return d.Error(fae.Wrap(err, "update medium"))
		}

		if d.IsTorrent() {
			if err := app.FlameTorrentRemove(d.Thash); err != nil {
				return fae.Wrap(err, "failed to remove torrent")
			}
		}

		d.Status = "done"
		err = app.DB.Download.Save(d)
		if err != nil {
			return fae.Wrap(err, "failed to save download")
		}

		notifier.Success("Downloads::Completed", fmt.Sprintf("%s - %s", d.Title, d.Display))
	}

	if len(moved) > 0 {
		dirs := lo.Map(moved, func(f *MoverFile, i int) string {
			return filepath.Dir(f.Destination)
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
