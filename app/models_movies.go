package app

import "fmt"

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
