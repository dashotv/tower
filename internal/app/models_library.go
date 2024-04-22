package app

func (c *Connector) LibraryGet(id string) (*Library, error) {
	m := &Library{}
	err := c.Library.Find(id, m)
	if err != nil {
		return nil, err
	}

	// post process here

	return m, nil
}

func (c *Connector) LibraryList(page, limit int) ([]*Library, error) {
	skip := (page - 1) * limit
	list, err := c.Library.Query().Limit(limit).Skip(skip).Run()
	if err != nil {
		return nil, err
	}

	return list, nil
}
