package app

func (c *Connector) MinionList() ([]*Minion, error) {
	list, err := c.Minion.Query().Limit(10).Run()
	if err != nil {
		return nil, err
	}

	return list, nil
}
