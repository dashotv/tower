package app

func (c *Connector) MovieList() ([]*Movie, error) {
	list, err := c.Movie.Query().Limit(10).Run()
	if err != nil {
		return nil, err
	}

	return list, nil
}
