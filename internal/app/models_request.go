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

func (c *Connector) RequestExists(guid string) (bool, error) {
	source, source_id := guidSplit(guid)
	list, err := c.Request.Query().Where("source", source).Where("source_id", source_id).Run()
	if err != nil {
		return false, err
	}
	if len(list) == 0 {
		return false, nil
	}

	return true, nil
}
