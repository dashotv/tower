package app

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/dashotv/fae"
)

var nzbgeekRegex = regexp.MustCompile("^https://api.nzbgeek")
var metubeRegex = regexp.MustCompile("^metube://")
var activeStates = []string{"searching", "loading", "managing", "downloading", "reviewing", "paused"}

func (c *Connector) DownloadGet(id string) (*Download, error) {
	d := &Download{}
	err := c.Download.Find(id, d)
	if err != nil {
		return nil, err
	}

	c.processDownloads([]*Download{d})
	return d, nil
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

func (c *Connector) DownloadByHash(hash string) (*Download, error) {
	list, err := c.Download.Query().Where("thash", hash).Run()
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
	list, err := q.In("status", activeStates).Run()
	if err != nil {
		return nil, err
	}

	c.processDownloads(list)
	return list, nil
}

func (c *Connector) RecentDownloads(mid string, page int) ([]*Download, int64, error) {
	total, err := app.DB.Download.Query().Where("status", "done").Count()
	if err != nil {
		return nil, 0, err
	}

	q := app.DB.Download.Query()

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

	results, err := q.Where("status", "done").
		Desc("updated_at").Desc("created_at").
		Skip((page - 1) * pagesize).
		Limit(pagesize).
		Run()
	if err != nil {
		return nil, 0, err
	}

	app.DB.processDownloads(results)
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
	for i, d := range list {
		m := &Medium{}
		err := app.DB.Medium.FindByID(d.MediumID, m)
		if err != nil {
			c.Log.Errorf("could not find medium: %s", d.MediumID)
			continue
		}

		d.Title = m.Title
		d.Kind = m.Kind
		d.Source = m.Source
		d.SourceID = m.SourceID
		d.Directory = m.Directory
		d.Active = m.Active
		d.Favorite = m.Favorite

		d.Search = &DownloadSearch{
			Type:       m.SearchParams.Type,
			Source:     m.SearchParams.Source,
			SourceID:   m.SourceID,
			Title:      m.Search,
			Resolution: m.SearchParams.Resolution,
			Group:      m.SearchParams.Group,
			Website:    m.SearchParams.Group,
			Exact:      false,
			Verified:   m.SearchParams.Verified,
			Uncensored: m.SearchParams.Uncensored,
			Bluray:     m.SearchParams.Bluray,
		}

		if m.Type == "Movie" {
			d.Search.SourceID = m.ImdbID
		}

		paths := m.Paths
		if m.Type == "Episode" && !m.SeriesID.IsZero() {
			s := &Series{}
			err := app.DB.Series.FindByID(m.SeriesID, s)
			if err != nil {
				c.Log.Errorf("could not find series: %s: %s", d.MediumID, err)
				continue
			}

			parts := strings.Split(s.Title, ":")
			title := parts[0]
			var shift int64
			if len(parts) > 1 {
				shift, _ = strconv.ParseInt(parts[1], 10, 64)
			}

			d.Title = title
			d.Kind = s.Kind
			d.Source = s.Source
			d.SourceID = s.SourceID
			d.Directory = s.Directory
			d.Active = s.Active
			d.Favorite = s.Favorite

			d.Search.Source = s.Source
			d.Search.SourceID = s.SourceID
			d.Search.Title = s.Search
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
				d.Display = fmt.Sprintf("#%d %s", m.AbsoluteNumber, m.Title)
			} else {
				d.Search.Season = m.SeasonNumber
				d.Search.Episode = m.EpisodeNumber
				d.Display = fmt.Sprintf("%02dx%02d %s", m.SeasonNumber, m.EpisodeNumber, m.Title)
			}

			unwatched, err := app.DB.SeriesUserUnwatched(s)
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
				err := app.DB.Medium.FindByID(f.MediumID, fm)
				if err != nil {
					c.Log.Errorf("could not find medium: %s", d.MediumID)
					continue
				}

				list[i].Files[j].Medium = fm
			}
		}

		list[i].Medium = m
	}
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
	err := app.DB.Download.Find(id, download)
	if err != nil {
		return err
	}

	for _, f := range download.Files {
		if f.Num == num {
			mid := primitive.ObjectID{}

			if mediumID != "" {
				mid, err = primitive.ObjectIDFromHex(mediumID)
				if err != nil {
					return err
				}
			}
			f.MediumID = mid

			return c.Download.Update(download)
		}
	}

	return fae.New("could not match num with download file")
}
