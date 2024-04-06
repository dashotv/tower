package app

func (c *Connector) RequestList() ([]*Request, error) {
	list, err := c.Request.Query().Limit(10).Run()
	if err != nil {
		return nil, err
	}

	return list, nil
}
