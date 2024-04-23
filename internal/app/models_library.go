package app

func (c *Connector) LibraryGet(id string) (*Library, error) {
	m := &Library{}
	err := c.Library.Find(id, m)
	if err != nil {
		return nil, err
	}

	c.processLibraries([]*Library{m})

	return m, nil
}

func (c *Connector) LibraryList(page, limit int) ([]*Library, error) {
	skip := (page - 1) * limit
	list, err := c.Library.Query().Limit(limit).Skip(skip).Run()
	if err != nil {
		return nil, err
	}

	c.processLibraries(list)
	return list, nil
}

func (c *Connector) processLibraries(list []*Library) {
	for _, l := range list {
		dt, err := c.LibraryTemplateGet(l.LibraryTemplateID.Hex())
		if err != nil {
			continue
		}
		l.LibraryTemplate = dt

		rt, err := c.LibraryTypeGet(l.LibraryTypeID.Hex())
		if err != nil {
			continue
		}
		l.LibraryType = rt
	}
}
