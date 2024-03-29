package app

import (
	"fmt"
	"regexp"

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
	if d.Url != "" {
		return d.Url, nil
	}

	if d.ReleaseId != "" {
		r := &Release{}
		err := app.DB.Release.Find(d.ReleaseId, r)
		if err != nil {
			return "", err
		}

		return r.Download, nil
	}

	return "", fae.New("no url or release")
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
		err := app.DB.Medium.FindByID(d.MediumId, m)
		if err != nil {
			c.Log.Errorf("could not find medium: %s", d.MediumId)
			continue
		}

		paths := m.Paths
		if m.Type == "Episode" && !m.SeriesId.IsZero() {
			s := &Series{}
			err := app.DB.Series.FindByID(m.SeriesId, s)
			if err != nil {
				c.Log.Errorf("could not find series: %s: %s", d.MediumId, err)
				continue
			}

			unwatched, err := app.DB.SeriesUserUnwatched(s)
			if err != nil {
				c.Log.Errorf("could not get unwatched count: %s: %s", s.ID.Hex(), err)
			}
			m.Unwatched = unwatched

			if isAnimeKind(string(s.Kind)) {
				m.Display = fmt.Sprintf("#%d %s", m.AbsoluteNumber, m.Title)
			} else {
				m.Display = fmt.Sprintf("%02dx%02d %s", m.SeasonNumber, m.EpisodeNumber, m.Title)
			}
			m.Title = s.Title
			m.Kind = s.Kind

			m.Source = s.Source
			m.SourceId = s.SourceId
			m.Search = s.Search
			if m.Search == "" {
				m.Search = s.Title
			}
			m.SearchParams = s.SearchParams
			m.Directory = s.Directory
			m.Active = s.Active
			m.Favorite = s.Favorite
			paths = s.Paths
		}

		for _, p := range paths {
			if p.Type == "cover" {
				m.Cover = fmt.Sprintf("%s/%s.%s", imagesBaseURL, p.Local, p.Extension)
				continue
			}
			if p.Type == "background" {
				m.Background = fmt.Sprintf("%s/%s.%s", imagesBaseURL, p.Local, p.Extension)
				continue
			}
		}

		for j, f := range d.Files {
			if !f.MediumId.IsZero() {
				fm := &Medium{}
				err := app.DB.Medium.FindByID(f.MediumId, fm)
				if err != nil {
					c.Log.Errorf("could not find medium: %s", d.MediumId)
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

func (c *Connector) DownloadSelect(id, mediumId string, num int) error {
	download := &Download{}
	err := app.DB.Download.Find(id, download)
	if err != nil {
		return err
	}

	for _, f := range download.Files {
		if f.Num == num {
			mid := primitive.ObjectID{}

			if mediumId != "" {
				mid, err = primitive.ObjectIDFromHex(mediumId)
				if err != nil {
					return err
				}
			}
			f.MediumId = mid

			return c.Download.Update(download)
		}
	}

	return fae.New("could not match num with download file")
}
