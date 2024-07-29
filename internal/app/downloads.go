package app

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"

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

func (a *Application) updateMedium(m *Medium, files []*MoverFile) error {
	// fmt.Printf("updateMedium: %s %+v\n", m.ID.Hex(), files)
	// only mark downloaded/completed if there are videos (so subtitles don't trigger it)
	videos := lo.Filter(files, func(f *MoverFile, i int) bool {
		return fileType(f.Destination) == "video"
	})
	if len(videos) > 0 {
		m.Broken = false
		m.Downloaded = true
		m.Completed = true
	}

	for _, f := range files {
		path := m.AddPathByFullpath(f.Destination)

		if err := a.Workers.Enqueue(&PathImport{ID: m.ID.Hex(), PathID: path.ID.Hex(), Title: f.Destination}); err != nil {
			return fae.Errorf("enqueue path: %s", err)
		}
	}

	if err := a.DB.Medium.Save(m); err != nil {
		return fae.Wrap(err, "failed to save medium")
	}

	return nil
}

func (a *Application) updateMedia(files []*MoverFile) error {
	for _, f := range files {
		if err := a.updateMedium(f.Medium, []*MoverFile{f}); err != nil {
			return fae.Wrapf(err, "update medium: %s", f.Destination)
		}
	}

	return nil
}

func (a *Application) downloadsCreate() error {
	// defer TickTock("DownloadsProcess: Create")()
	seriesDownloads, err := a.DB.SeriesDownloadCounts()
	if err != nil {
		return fae.Wrap(err, "failed to get series download counts")
	}
	seriesMulti, err := a.DB.SeriesMultiDownloads()
	if err != nil {
		return fae.Wrap(err, "failed to get series multi")
	}

	list, err := a.DB.UpcomingNow()
	if err != nil {
		return fae.Wrap(err, "failed to get upcoming episodes")
	}

	for _, ep := range list {
		//app.Workers.Log.Debugf("DownloadsProcess: create: %s %s", ep.Title, ep.Display)
		if !ep.SeriesActive {
			//app.Workers.Log.Debugf("DownloadsProcess: create: %s %s: not active", ep.Title, ep.Display)
			continue
		}

		unwatched, err := a.DB.SeriesUnwatchedByID(ep.SeriesID.Hex())
		if err != nil {
			return fae.Wrap(err, "failed to get unwatched")
		}

		if unwatched+seriesDownloads[ep.SeriesID.Hex()] >= 3 || seriesMulti[ep.SeriesID.Hex()] {
			continue
		}

		a.Workers.Log.Debugf("download created %s - %s", ep.SeriesTitle, ep.Display)
		seriesDownloads[ep.SeriesID.Hex()]++

		d := &Download{
			Status:   "searching",
			MediumID: ep.ID,
			Auto:     true,
		}
		err = a.DB.Download.Save(d)
		if err != nil {
			return fae.Wrap(err, "failed to save download")
		}

		err = a.DB.EpisodeSetting(ep.ID.Hex(), "downloaded", true)
		if err != nil {
			return fae.Wrap(err, "failed to save episode")
		}
	}

	return nil
}

func (a *Application) downloadsSearch() error {
	// defer TickTock("DownloadsProcess: Search")()
	list, err := a.DB.DownloadByStatus("searching")
	if err != nil {
		return fae.Wrap(err, "failed to get downloads")
	}

	for _, d := range list {
		if d.Medium == nil {
			continue
		}
		if d.Medium.Type != "Episode" {
			//movies handled in downloadsSearchMovies
			continue
		}

		match, err := a.ScrySearchEpisode(d.Search)
		if err != nil {
			return fae.Wrap(err, "failed to search releases")
		}
		if match == nil {
			continue
		}

		a.Workers.Log.Debugf("download found %s - %s", d.Title, d.Display)
		notifier.Info("SearchFound", fmt.Sprintf("%s - %s", d.Title, d.Display))

		d.Status = "loading"
		if !a.Config.Production {
			d.Status = "reviewing"
		}
		d.URL = match.Download
		tags := []string{}
		if match.Group != "" {
			tags = append(tags, match.Group)
		}
		if match.Website != "" {
			tags = append(tags, match.Website)
		}
		if match.Resolution != "" {
			tags = append(tags, match.Resolution+"p")
		}
		d.Tag = strings.Join(tags, " ")

		err = a.DB.Download.Save(d)
		if err != nil {
			return fae.Wrap(err, "failed to save download")
		}
	}

	return nil
}

func (a *Application) downloadsSearchMovies() error {
	// l := a.Workers.Log.Named("downloads_movies")

	query := a.DB.Movie.Query().
		Where("downloaded", false).Where("completed", false).
		LessThanEqual("release_date", time.Now()).
		NotIn("kind", []string{"movies3d", "movies4k"}).Desc("created_at")
	err := query.Batch(100, func(movies []*Movie) error {
		for _, movie := range movies {
			d := &Download{MediumID: movie.ID, Status: "searching", Auto: true}
			a.DB.processDownload(d)

			found, err := a.ScrySearchMovie(d.Search)
			if err != nil {
				return fae.Wrap(err, "failed to search releases")
			}
			if found == nil {
				continue
			}

			name := d.Display
			if name == "" {
				name = d.Title
			}
			a.Workers.Log.Debugf("download found (movie) %s", name)
			notifier.Info("SearchFound", fmt.Sprintf("%s (movie)", name))

			d.Status = "loading"
			if !a.Config.Production {
				d.Status = "reviewing"
			}
			d.URL = found.Download

			if err := a.DB.Download.Save(d); err != nil {
				return fae.Wrap(err, "failed to save download")
			}

			movie.Downloaded = true
			if err := a.DB.Movie.Save(movie); err != nil {
				return fae.Wrap(err, "failed to save movie")
			}
		}
		return nil
	})
	if err != nil {
		return fae.Wrap(err, "failed to query movies")
	}

	return nil
}

func (a *Application) downloadsLoad() (err error) {
	// defer TickTock("DownloadsProcess: Load")()
	list, err := a.DB.DownloadByStatus("loading")
	if err != nil {
		return fae.Wrap(err, "failed to get downloads")
	}

	for _, d := range list {
		if d.ReleaseID == "" && d.URL == "" {
			a.DB.Log.Debugf("DownloadsProcess: load: %s %s: no release", d.Title, d.Display)
			continue
		}

		res, err := a.FlameAdd(d)
		if err != nil {
			return d.Error(fae.Wrap(err, "failed to add to flame"))
		}

		d.Status = "downloading"
		if d.IsTorrent() {
			d.Status = "managing"
		}
		d.Thash = res

		err = a.DB.Download.Save(d)
		if err != nil {
			return fae.Wrap(err, "failed to save download")
		}
	}

	return nil
}

func (a *Application) downloadsManage() error {
	// defer TickTock("DownloadsProcess: Manage")()
	list, err := a.DB.DownloadByStatus("managing")
	if err != nil {
		return fae.Wrap(err, "get downloads")
	}

	for _, d := range list {
		// TODO: manage metube? show files while downloading?
		if d.Thash == "" || !d.IsTorrent() {
			continue
		}

		if d.Medium == nil {
			a.Workers.Log.Warnf("no medium", d.ID.Hex())
			continue
		}

		t, err := a.FlameTorrent(d.Thash)
		if err != nil {
			a.Log.Named("downloads.manage").Errorf("failed to get torrent: %s", err)
			continue
		}

		if len(t.Files) == 0 {
			continue
		}

		if err := a.downloadsManageOne(d, t); err != nil {
			a.Log.Errorf("failed to manage download: %s", err)
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
		a.Workers.Log.Warnf("download has no files: %s", d.ID.Hex())
		return nil
	}

	// TODO: handle downloads with single media file and multiple subtitles

	if len(d.Files) == 1 {
		d.Files[0].MediumID = d.MediumID
		d.Status = "downloading"

		if err := a.DB.Download.Save(d); err != nil {
			return fae.Wrap(err, "failed to save download")
		}

		return nil
	}

	if !d.Multi {
		a.Workers.Log.Warnf("multiple files, but not multi", d.ID.Hex())

		d.Status = "reviewing"
		if err := a.DB.Download.Save(d); err != nil {
			return fae.Wrap(err, "failed to save download")
		}

		return nil
	}

	if d.Medium.Type != "Series" {
		// only handle series for now
		a.Workers.Log.Warnf("multi not series", d.ID.Hex())

		d.Status = "reviewing"
		if err := a.DB.Download.Save(d); err != nil {
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
		a.Log.Debugf("searching for episode: %s %s", d.Search.Type, title)
		ep, err := a.RunicFindEpisode(d.MediumID, title, d.Search.Type)
		if err != nil {
			return fae.Wrap(err, "failed to find episode")
		}

		if ep == nil {
			a.Workers.Log.Warnf("episode not found: %s", file.Name)
			continue
		}
		a.Log.Debugf("found: %s", ep.Title)

		df.MediumID = ep.ID
	}

	a.DB.processDownloads([]*Download{d})

	// TODO: handle want more / none / etc
	wanted := false
	for _, f := range t.Files {
		if f.Priority > 0 {
			wanted = true
			break
		}
	}

	if wanted && t.Progress != 100 {
		err := a.FlameTorrentWant(d.Thash, "none")
		if err != nil {
			return fae.Wrap(err, "want none")
		}
	}

	if d.HasMedia() {
		nums := d.NextFileNums(t, downloadMultiFiles)
		if nums != "" {
			err := a.FlameTorrentWant(d.Thash, nums)
			if err != nil {
				return fae.Wrap(err, "want next")
			}
		}

		// save updates to download files
		d.Status = "downloading"
	}

	if err := a.DB.Download.Save(d); err != nil {
		return fae.Wrap(err, "failed to save download")
	}

	return nil
}

func (a *Application) downloadsMove() error {
	// defer TickTock("DownloadsProcess: Move")()
	list, err := a.DB.DownloadByStatus("downloading")
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
			t, err = a.FlameTorrent(d.Thash)
			if err != nil {
				notifier.Log.Errorf("Downloads::Move", "failed to get torrent: %s", err)
				d.Error(err)
				continue
			}
		}

		mover := NewMover(a.Log.Named("mover"), d, t)
		files, err := mover.Move()
		if err != nil {
			a.Log.Debugf("error: %+v", err)
			return d.Error(fae.Wrap(err, "move download"))
		}

		if d.Multi && d.Medium.Type == "Series" {
			// update medium and add path
			if files != nil && len(files) > 0 && a.Config.Production {
				moved = append(moved, files...)
				if err := a.updateMedia(files); err != nil {
					return d.Error(fae.Wrap(err, "update medium"))
				}
			}

			wanted := lo.Filter(t.Files, func(f *qbt.TorrentFile, i int) bool {
				return f.Priority > 0 && f.Progress < 100
			})
			if len(wanted) < 3 {
				nums := d.NextFileNums(t, 3)
				if nums != "" {
					err := a.FlameTorrentWant(d.Thash, nums)
					if err != nil {
						return d.Error(fae.Wrap(err, "want next"))
					}
				}
			}

			continue
		}

		if files == nil || len(files) == 0 {
			continue
		}

		moved = append(moved, files...)

		if a.Config.Production {
			// update medium and add path
			if err := a.updateMedia(files); err != nil {
				return d.Error(fae.Wrap(err, "update medium"))
			}

			if d.IsTorrent() {
				if err := a.FlameTorrentRemove(d.Thash); err != nil {
					return fae.Wrap(err, "failed to remove torrent")
				}
			}

			d.Status = "done"
			err = a.DB.Download.Save(d)
			if err != nil {
				return fae.Wrap(err, "failed to save download")
			}

			notifier.Success("Downloads::Completed", fmt.Sprintf("%s - %s", d.Title, d.Display))
		}
	}

	if len(moved) > 0 && a.Config.Production {
		dirs := lo.Map(moved, func(f *MoverFile, i int) string {
			return filepath.Dir(f.Destination)
		})
		dirs = lo.Uniq(dirs)

		for _, dir := range dirs {
			notifier.Log.Info("refresh", dir)
			err := a.Plex.RefreshLibraryPath(dir)
			if err != nil {
				return fae.Wrap(err, "failed to refresh library")
			}
		}
	}

	return nil
}
