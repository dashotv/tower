package app

func (c *Connector) IndexerGet(id string) (*Indexer, error) {
	m, err := c.Indexer.Get(id, &Indexer{})
	if err != nil {
		return nil, err
	}

	// post process here

	return m, nil
}

func (c *Connector) IndexerList(page, limit int) ([]*Indexer, int64, error) {
	q := c.Indexer.Query().Limit(limit).Skip((page - 1) * limit).Desc("created_at")

	count, err := q.Count()
	if err != nil {
		return nil, 0, err
	}

	indexers, err := q.Run()
	if err != nil {
		return nil, 0, err
	}

	return indexers, count, nil
}
