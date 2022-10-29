package app

func (c *Connector) MovieSetting(id, setting string, value bool) error {
	m := &Movie{}
	err := c.Movie.Find(id, m)
	if err != nil {
		return err
	}

	App().Log.Infof("movie setting: %s %t", setting, value)
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

func (c *Connector) MoviePaths(id string) ([]Path, error) {
	m := &Movie{}
	err := App().DB.Movie.Find(id, m)
	if err != nil {
		return nil, err
	}

	return m.Paths, nil
}
