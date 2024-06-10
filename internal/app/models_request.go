package app

func (c *Connector) RequestList(page, limit int) ([]*Request, int64, error) {
	skip := (page - 1) * limit
	count, err := c.Request.Query().Count()
	if err != nil {
		return nil, -1, err
	}

	list, err := c.Request.Query().Desc("created_at").Limit(limit).Skip(skip).Run()
	if err != nil {
		return nil, -1, err
	}

	return list, count, nil
}
