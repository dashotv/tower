package app

func (c *Connector) MovieSetting(id, setting string, value bool) error {
	m := &Movie{}
	err := c.Movie.Find(id, m)
	if err != nil {
		return err
	}

	c.log.Infof("movie setting: %s %t", setting, value)
	switch setting {
	case "active":
		m.Active = value
	case "favorite":
		m.Favorite = value
	case "broken":
		m.Broken = value
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

	err = c.Movie.Update(m)
	if err != nil {
		return err
	}

	return events.Send("tower.movies", &EventTowerMovie{Event: "update", ID: m.ID.Hex(), Movie: m})
}
func (c *Connector) MoviePaths(id string) ([]*Path, error) {
	m := &Movie{}
	err := db.Movie.Find(id, m)
	if err != nil {
		return nil, err
	}

	return m.Paths, nil
}
