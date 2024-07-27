package app

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/samber/lo"

	"github.com/dashotv/fae"
	"github.com/dashotv/flame/nzbget"
	"github.com/dashotv/flame/qbt"
)

var nzbgeekRegex = regexp.MustCompile("^https://api.nzbgeek")
var metubeRegex = regexp.MustCompile("^metube://")
var activeStates = []string{"searching", "loading", "managing", "downloading", "reviewing"}
var downloadStates = []string{"reviewing", "searching", "loading", "managing", "downloading", "done"}

func (c *Connector) DownloadGet(id string) (*Download, error) {
	d := &Download{}
	err := c.Download.Find(id, d)
	if err != nil {
		return nil, err
	}

	c.processDownloads([]*Download{d})
	return d, nil
}

func (d *Download) Error(e error) error {
	d.Status = "reviewing"

	if err := app.DB.Download.Save(d); err != nil {
		return errors.Join(e, fae.Wrap(err, "saving download"))
	}

	return e
}

func (d *Download) StatusIndex() int {
	return slices.Index(downloadStates, d.Status)
}

func (d *Download) GetURL() (string, error) {
	if d.URL != "" {
		return d.URL, nil
	}

	if d.ReleaseID != "" {
		r := &Release{}
		err := app.DB.Release.Find(d.ReleaseID, r)
		if err != nil {
			return "", err
		}

		return r.Download, nil
	}

	return "", fae.New("no url or release")
}

func (d *Download) SortedFileNums(t *qbt.Torrent) ([]string, error) {
	if t == nil {
		return nil, fae.New("no torrent")
	}
	grouped := lo.GroupBy(d.Files, func(df *DownloadFile) string {
		if df.MediumID.IsZero() {
			return fmt.Sprintf("100%03d", df.Num)
		}
		s := df.Medium.SeasonNumber
		if s == 0 {
			s = 100 // sort specials last
		}
		return fmt.Sprintf("%03d%03d", s, df.Medium.EpisodeNumber)
	})

	keys := lo.Keys(grouped)
	sort.Strings(keys)

	list := []string{}
	for _, key := range keys {
		for _, df := range grouped[key] {
			if df.MediumID != primitive.NilObjectID && t.Files[df.Num].Progress < 100 && !df.Medium.Downloaded {
				list = append(list, fmt.Sprintf("%d", df.Num))
			}
		}
	}

	return list, nil
}

func (d *Download) NextFileNums(t *qbt.Torrent, n int) string {
	list, err := d.SortedFileNums(t)
	if err != nil {
		return ""
	}
	if len(list) == 0 {
		return ""
	}

	if len(list) > 3 {
		list = list[:3]
	}
	return strings.Join(list, ",")
}

func (d *Download) HasMedia() bool {
	if !d.Multi {
		return d.MediumID != primitive.NilObjectID
	}

	has := lo.Filter(d.Files, func(f *DownloadFile, _ int) bool {
		return !f.MediumID.IsZero()
	})

	return len(has) > 0
}

func (c *Connector) DownloadByHash(hash string) (*Download, error) {
	list, err := c.Download.Query().In("status", activeStates).Where("thash", hash).Run()
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, fae.Errorf("could not find download by hash: %s", hash)
	}
	if len(list) > 1 {
		return nil, fae.Errorf("multiple downloads found by hash: %s", hash)
	}

	c.processDownloads(list)
	return list[0], nil
}

func (d *Download) IsNzb() bool {
	url, err := d.GetURL()
	if err != nil {
		return false
	}

	if nzbgeekRegex.MatchString(url) {
		return true
	}

	return false
}

func (d *Download) IsMetube() bool {
	url, err := d.GetURL()
	if err != nil {
		return false
	}

	if metubeRegex.MatchString(url) {
		return true
	}

	return false
}

func (d *Download) IsTorrent() bool {
	url, err := d.GetURL()
	if err != nil {
		return false
	}

	if !nzbgeekRegex.MatchString(url) && !metubeRegex.MatchString(url) {
		return true
	}

	return false
}

func (c *Connector) ActiveDownloads() ([]*Download, error) {
	q := c.Download.Query()
	list, err := q.In("status", activeStates).Limit(-1).Run()
	if err != nil {
		return nil, err
	}

	c.processDownloads(list)
	return list, nil
}

func (c *Connector) RecentDownloads(mid string, page int) ([]*Download, int64, error) {
	total, err := c.Download.Query().Count()
	if err != nil {
		return nil, 0, err
	}

	q := c.Download.Query()

	if mid != "" {
		m, err := c.Medium.Get(mid, &Medium{})
		if err != nil {
			return nil, 0, err
		}

		ids := []primitive.ObjectID{m.ID}
		if m.Type == "Series" {
			eps, err := c.SeriesSeasonEpisodesAll(m.ID.Hex())
			if err != nil {
				return nil, 0, err
			}
			for _, e := range eps {
				ids = append(ids, e.ID)
			}
		}

		q = q.In("medium_id", ids)
	}

	results, err := q.
		Desc("updated_at").Desc("created_at").
		Skip((page - 1) * pagesize).
		Limit(pagesize).
		Run()
	if err != nil {
		return nil, 0, err
	}

	c.processDownloads(results)
	return results, total, nil
}

func (c *Connector) DownloadByStatus(status string) ([]*Download, error) {
	list, err := c.Download.Query().Where("status", status).Run()
	if err != nil {
		return nil, err
	}

	c.processDownloads(list)
	return list, nil
}

func (c *Connector) processDownloads(list []*Download) {
	for _, d := range list {
		c.processDownload(d)
	}
}

func (c *Connector) processDownload(d *Download) {
	m := &Medium{}
	err := c.Medium.FindByID(d.MediumID, m)
	if err != nil {
		c.Log.Errorf("could not find medium: %s", d.MediumID)
		return
	}

	d.Title = m.Title
	d.Kind = m.Kind
	d.Source = m.Source
	d.SourceID = m.SourceID
	d.Directory = m.Directory
	d.Active = m.Active
	d.Favorite = m.Favorite

	d.Search = &DownloadSearch{
		SourceID: m.SourceID,
		Title:    m.Search,
		Exact:    false,
	}
	if m.SearchParams != nil {
		d.Search.Type = m.SearchParams.Type
		d.Search.Source = m.SearchParams.Source
		d.Search.Resolution = m.SearchParams.Resolution
		d.Search.Group = m.SearchParams.Group
		d.Search.Website = m.SearchParams.Group
		d.Search.Verified = m.SearchParams.Verified
		d.Search.Uncensored = m.SearchParams.Uncensored
		d.Search.Bluray = m.SearchParams.Bluray
	}

	if m.Type == "Movie" {
		d.Search.SourceID = m.ImdbID
		d.Display = m.Display
		if !m.ReleaseDate.IsZero() {
			d.Search.Year = m.ReleaseDate.Year()
		}
	}

	paths := m.Paths
	if m.Type == "Episode" && !m.SeriesID.IsZero() {
		m.ApplyOverrides()

		s := &Series{}
		err := c.Series.FindByID(m.SeriesID, s)
		if err != nil {
			c.Log.Errorf("could not find series: %s: %s", d.MediumID, err)
			return
		}

		parts := strings.Split(s.Search, ":")
		title := parts[0]
		var shift int64
		if len(parts) > 1 {
			shift, _ = strconv.ParseInt(parts[1], 10, 64)
		}

		d.Title = s.Title
		d.Kind = s.Kind
		d.Source = s.Source
		d.SourceID = s.SourceID
		d.Directory = s.Directory
		d.Active = s.Active
		d.Favorite = s.Favorite

		d.Search.Source = s.Source
		d.Search.SourceID = s.SourceID
		d.Search.Title = title
		d.Search.Type = s.SearchParams.Type
		d.Search.Source = s.SearchParams.Source
		d.Search.Resolution = s.SearchParams.Resolution
		d.Search.Group = s.SearchParams.Group
		d.Search.Website = s.SearchParams.Group
		d.Search.Verified = s.SearchParams.Verified
		d.Search.Uncensored = s.SearchParams.Uncensored
		d.Search.Bluray = s.SearchParams.Bluray

		if isAnimeKind(string(s.Kind)) && m.AbsoluteNumber > 0 {
			n := m.AbsoluteNumber
			if shift > 0 && n > int(shift) {
				n = n - int(shift)
			}
			d.Search.Episode = n
			d.Display = fmt.Sprintf("%02dx%02d #%d %s", m.SeasonNumber, m.EpisodeNumber, m.AbsoluteNumber, m.Title)
		} else {
			d.Search.Season = m.SeasonNumber
			d.Search.Episode = m.EpisodeNumber
			d.Display = fmt.Sprintf("%02dx%02d %s", m.SeasonNumber, m.EpisodeNumber, m.Title)
		}

		unwatched, err := c.SeriesUserUnwatched(s)
		if err != nil {
			c.Log.Errorf("could not get unwatched count: %s: %s", s.ID.Hex(), err)
		}
		d.Unwatched = unwatched

		paths = s.Paths
	}

	for _, p := range paths {
		if p.Type == "cover" {
			d.Cover = fmt.Sprintf("%s/%s.%s", imagesBaseURL, p.Local, p.Extension)
			continue
		}
		if p.Type == "background" {
			d.Background = fmt.Sprintf("%s/%s.%s", imagesBaseURL, p.Local, p.Extension)
			continue
		}
	}

	for j, f := range d.Files {
		if !f.MediumID.IsZero() {
			fm := &Medium{}
			err := c.Medium.FindByID(f.MediumID, fm)
			if err != nil {
				c.Log.Errorf("could not find medium: %s", d.MediumID)
				continue
			}

			d.Files[j].Medium = fm
		}
	}

	// 	completed := lo.Filter(d.Files, func(file *DownloadFile, _ int) bool {
	// 		tf := t[0].Files[file.Num]
	// 		return !file.MediumID.IsZero() && tf.Progress == 100
	// 	})
	// 	g.Files.Completed = len(completed)
	//
	// 	selected := lo.Filter(d.Files, func(file *DownloadFile, _ int) bool {
	// 		return !file.MediumID.IsZero()
	// 	})
	// 	g.Files.Selected = len(selected)

	d.Medium = m
}

func (db *Connector) processDownloadExtra(d *Download, c *FlameCombined) {
	// if thashNumbersRegex.MatchString(d.Thash) && len(c.Nzbs) > 0 {
	if d.IsNzb() && len(c.Nzbs) > 0 {
		n := lo.Filter(c.Nzbs, func(nzb nzbget.Group, _ int) bool {
			return fmt.Sprintf("%d", nzb.ID) == d.Thash
		})

		if len(n) > 0 {
			if err := db.processDownloadExtraNzb(d, n[0], c.NzbStatus); err != nil {
				app.Log.Errorf("error handling nzb: %v", err)
			}
		}
	}

	if d.IsTorrent() && len(c.Torrents) > 0 {
		t := lo.Filter(c.Torrents, func(torrent *qbt.TorrentJSON, _ int) bool {
			return strings.ToLower(torrent.Hash) == strings.ToLower(d.Thash)
		})
		if len(t) > 0 {
			if err := db.processDownloadExtraTorrent(d, t[0]); err != nil {
				app.Log.Errorf("error handling torrent: %v", err)
			}
		}
	}
}

func (db *Connector) processDownloadExtraTorrent(d *Download, t *qbt.TorrentJSON) error {
	d.Torrent = t
	d.TorrentState = t.State
	if t.Queue > 0 {
		d.Queue = t.Queue
	}
	d.Progress = t.Progress
	if t.Finish > 0 && t.Finish != 8640000 {
		d.Eta = time.Now().Add(time.Duration(t.Finish) * time.Second).Format(time.RFC3339)
	}

	// set torrent file on download files
	if len(t.Files) > 0 {
		for _, file := range d.Files {
			if file.Num >= len(t.Files) {
				continue
			}
			file.TorrentFile = t.Files[file.Num]
		}
	}

	if !d.Multi || len(d.Files) == 0 || len(t.Files) == 0 {
		return nil
	}

	{
		completed := lo.Filter(d.Files, func(file *DownloadFile, _ int) bool {
			if file.Num >= len(t.Files) {
				return false
			}
			tf := t.Files[file.Num]
			return !file.MediumID.IsZero() && tf.Progress == 100
		})
		d.FilesCompleted = len(completed)
	}

	{
		selected := lo.Filter(d.Files, func(file *DownloadFile, _ int) bool {
			return !file.MediumID.IsZero()
		})
		d.FilesSelected = len(selected)
	}

	{
		wanted := lo.Filter(t.Files, func(file *qbt.TorrentFile, _ int) bool {
			return file.Priority > 0 && file.Progress < 100
		})
		d.FilesWanted = len(wanted)
	}

	// sort files by torrent name
	selected := lo.Filter(d.Files, func(item *DownloadFile, index int) bool {
		return item.TorrentFile != nil && item.MediumID != primitive.NilObjectID
	})
	ignored := lo.Filter(d.Files, func(item *DownloadFile, index int) bool {
		return item.MediumID == primitive.NilObjectID
	})
	missing := lo.Filter(d.Files, func(item *DownloadFile, index int) bool {
		return item.TorrentFile == nil
	})

	slices.SortFunc(selected, func(a, b *DownloadFile) int {
		if a.Medium == nil || b.Medium == nil {
			return strings.Compare(a.TorrentFile.Name, b.TorrentFile.Name)
		}
		return strings.Compare(
			fmt.Sprintf("%03d %03d %03d %s %s", a.Medium.AbsoluteNumber, a.Medium.SeasonNumber, a.Medium.EpisodeNumber, a.Medium.Title, a.Medium.Display),
			fmt.Sprintf("%03d %03d %03d %s %s", b.Medium.AbsoluteNumber, b.Medium.SeasonNumber, b.Medium.EpisodeNumber, b.Medium.Title, b.Medium.Display),
		)
	})
	slices.SortFunc(ignored, func(a, b *DownloadFile) int {
		return strings.Compare(a.TorrentFile.Name, b.TorrentFile.Name)
	})
	d.Files = append(selected, ignored...)
	d.Files = append(d.Files, missing...)

	return nil
}

func (db *Connector) processDownloadExtraNzb(d *Download, n nzbget.Group, status nzbget.Status) error {
	s := 0
	if status.DownloadRate > 0 {
		s = ((n.RemainingSizeMB * 1024) / (status.DownloadRate / 1024)) * 1000
	}
	d.Queue = float64(n.ID)
	d.Progress = 100.0 - (float64(n.RemainingSizeMB)/float64(n.FileSizeMB))*100.0
	if s > 0 {
		d.Eta = time.Now().Add(time.Duration(s) * time.Second).Format(time.RFC3339)
	}

	return nil
}

func (c *Connector) DownloadSetting(id, setting string, value bool) error {
	d := &Download{}
	err := c.Download.Find(id, d)
	if err != nil {
		return err
	}

	switch setting {
	case "auto":
		d.Auto = value
	case "multi":
		d.Multi = value
	case "force":
		d.Force = value
	}

	return c.Download.Update(d)
}

func (c *Connector) DownloadSelect(id, mediumID string, num int) error {
	download := &Download{}
	err := c.Download.Find(id, download)
	if err != nil {
		return err
	}

	for _, f := range download.Files {
		if f.Num == num {
			if mediumID == "" {
				f.MediumID = primitive.NilObjectID
				return c.Download.Update(download)
			}

			mid, err := primitive.ObjectIDFromHex(mediumID)
			if err != nil {
				return err
			}
			f.MediumID = mid
			return c.Download.Update(download)
		}
	}

	return fae.New("could not match num with download file")
}
func (c *Connector) DownloadClear(id string, nums string) error {
	list := strings.Split(nums, ",")
	if len(list) == 0 {
		return fae.New("no nums")
	}

	download := &Download{}
	err := c.Download.Find(id, download)
	if err != nil {
		return err
	}

	files := lo.Filter(download.Files, func(f *DownloadFile, _ int) bool {
		return lo.Contains(list, fmt.Sprintf("%d", f.Num)) && f.MediumID != primitive.NilObjectID
	})

	for _, f := range files {
		f.MediumID = primitive.NilObjectID
	}

	return c.Download.Save(download)

	//	for _, f := range download.Files {
	//		if f.Num == num {
	//			if mediumID == "" {
	//				f.MediumID = primitive.NilObjectID
	//				return c.Download.Update(download)
	//			}
	//
	//			mid, err := primitive.ObjectIDFromHex(mediumID)
	//			if err != nil {
	//				return err
	//			}
	//			f.MediumID = mid
	//			return c.Download.Update(download)
	//		}
	//	}
	//
	// return fae.New("could not match num with download file")
}
