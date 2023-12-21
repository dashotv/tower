package app

func (c *Connector) PinList() ([]*Pin, error) {
	list, err := c.Pin.Query().Limit(10).Run()
	if err != nil {
		return nil, err
	}

	return list, nil
}
