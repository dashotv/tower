package app

func (c *Connector) LibraryTemplateGet(id string) (*LibraryTemplate, error) {
	m := &LibraryTemplate{}
	err := c.LibraryTemplate.Find(id, m)
	if err != nil {
		return nil, err
	}

	// post process here

	return m, nil
}

func (c *Connector) LibraryTemplateList(page, limit int) ([]*LibraryTemplate, error) {
	skip := (page - 1) * limit
	list, err := c.LibraryTemplate.Query().Limit(limit).Skip(skip).Run()
	if err != nil {
		return nil, err
	}

	return list, nil
}
