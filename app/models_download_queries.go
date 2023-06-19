package app

import (
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var activeStates = []string{"searching", "loading", "managing", "downloading", "reviewing", "paused"}

func (c *Connector) ActiveDownloads() ([]*Download, error) {
	q := c.Download.Query()
	list, err := q.In("status", activeStates).Run()
	if err != nil {
		return nil, err
	}

	processDownloads(list)
	return list, nil
}

func (c *Connector) RecentDownloads(page int) ([]*Download, error) {
	q := c.Download.Query()
	list, err := q.Where("status", "done").
		Desc("updated_at").Desc("created_at").
		Skip((page - 1) * pagesize).
		Limit(pagesize).
		Run()
	if err != nil {
		return nil, err
	}

	processDownloads(list)
	return list, nil
}

func processDownloads(list []*Download) {
	for i, d := range list {
		m := &Medium{}
		err := App().DB.Medium.FindByID(d.MediumId, m)
		if err != nil {
			App().Log.Errorf("could not find medium: %s", d.MediumId)
			continue
		}

		paths := m.Paths
		m.Display = m.Type
		if m.Type == "Episode" && !m.SeriesId.IsZero() {
			s := &Series{}
			err := App().DB.Series.FindByID(m.SeriesId, s)
			if err != nil {
				App().Log.Errorf("could not find series: %s: %s", d.MediumId, err)
				continue
			}

			unwatched, err := App().DB.SeriesAllUnwatched(s)
			if err != nil {
				App().Log.Errorf("could not get unwatched count: %s: %s", s.ID.Hex(), err)
			}

			m.Kind = s.Kind
			m.Unwatched = unwatched
			m.Display = fmt.Sprintf("%dx%d %s", m.SeasonNumber, m.EpisodeNumber, m.Title)
			m.Title = s.Title
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
				err := App().DB.Medium.FindByID(f.MediumId, fm)
				if err != nil {
					App().Log.Errorf("could not find medium: %s", d.MediumId)
					continue
				}

				list[i].Files[j].Medium = fm
			}
		}

		list[i].Medium = *m
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
	err := App().DB.Download.Find(id, download)
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

	return errors.New("could not match num with download file")
}
