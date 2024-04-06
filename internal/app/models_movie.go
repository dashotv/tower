package app

import (
	"fmt"

	"github.com/dashotv/fae"
)

func (c *Connector) MovieGet(id string) (*Movie, error) {
	movie, err := c.Movie.Get(id, &Movie{})
	if err != nil {
		return nil, err
	}

	// if err := c.processMovies([]*Movie{movie}); err != nil {
	// 	return nil, err
	// }

	return movie, nil
}

func (c *Connector) MovieList(page, limit int) ([]*Movie, error) {
	skip := (page - 1) * limit
	list, err := c.Movie.Query().Desc("created_at").Limit(limit).Skip(skip).Run()
	if err != nil {
		return nil, fae.Wrap(err, "query failed")
	}

	// if err := c.processMovies(list); err != nil {
	// 	return nil, fae.Wrap(err, "process movies failed")
	// }

	return list, nil
}

func (c *Connector) MoviesAll() ([]*Movie, error) {
	return c.Movie.Query().Limit(-1).Run()
}

func (c *Connector) processMovies(list []*Movie) error {
	for _, m := range list {
		for _, p := range m.Paths {
			if p.Type == "cover" {
				m.Cover = fmt.Sprintf("%s/%s.%s", imagesBaseURL, p.Local, p.Extension)
				continue
			}
			if p.Type == "background" {
				m.Background = fmt.Sprintf("%s/%s.%s", imagesBaseURL, p.Local, p.Extension)
				continue
			}
		}
	}
	return nil
}

func (c *Connector) MovieSetting(id, setting string, value bool) error {
	m := &Movie{}
	err := c.Movie.Find(id, m)
	if err != nil {
		return err
	}

	c.Log.Infof("movie setting: %s %t", setting, value)
	switch setting {
	case "broken":
		m.Broken = value
	case "completed":
		m.Completed = value
	case "downloaded":
		m.Downloaded = value
	}

	return c.Movie.Update(m)
}

func (c *Connector) MovieUpdate(id string, data *Movie) error {
	m := &Movie{}
	err := c.Movie.Find(id, m)
	if err != nil {
		return err
	}

	m.Display = data.Display
	m.Directory = data.Directory
	m.Kind = data.Kind
	m.Source = data.Source
	m.SourceId = data.SourceId
	m.Search = data.Search

	return c.Movie.Update(m)
}

func (c *Connector) MoviePaths(id string) ([]*Path, error) {
	m := &Movie{}
	err := c.Movie.Find(id, m)
	if err != nil {
		return nil, err
	}

	return m.Paths, nil
}
