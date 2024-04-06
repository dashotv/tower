package app

func (c *Connector) UserList() ([]*User, error) {
	list, err := c.User.Query().Limit(10).Run()
	if err != nil {
		return nil, err
	}

	return list, nil
}
