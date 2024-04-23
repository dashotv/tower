package app

func (c *Connector) LibraryTypeGet(id string) (*LibraryType, error) {
	m := &LibraryType{}
	err := c.LibraryType.Find(id, m)
	if err != nil {
		return nil, err
	}

	// post process here

	return m, nil
}

func (c *Connector) LibraryTypeList(page, limit int) ([]*LibraryType, error) {
	skip := (page - 1) * limit
	list, err := c.LibraryType.Query().Limit(limit).Skip(skip).Run()
	if err != nil {
		return nil, err
	}

	return list, nil
}
