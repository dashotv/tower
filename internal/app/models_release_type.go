package app

func (c *Connector) ReleaseTypeGet(id string) (*ReleaseType, error) {
	m := &ReleaseType{}
	err := c.ReleaseType.Find(id, m)
	if err != nil {
		return nil, err
	}

	// post process here

	return m, nil
}

func (c *Connector) ReleaseTypeList(page, limit int) ([]*ReleaseType, error) {
	skip := (page - 1) * limit
	list, err := c.ReleaseType.Query().Limit(limit).Skip(skip).Run()
	if err != nil {
		return nil, err
	}

	return list, nil
}
