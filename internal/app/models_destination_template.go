package app

func (c *Connector) DestinationTemplateGet(id string) (*DestinationTemplate, error) {
	m := &DestinationTemplate{}
	err := c.DestinationTemplate.Find(id, m)
	if err != nil {
		return nil, err
	}

	// post process here

	return m, nil
}

func (c *Connector) DestinationTemplateList(page, limit int) ([]*DestinationTemplate, error) {
	skip := (page - 1) * limit
	list, err := c.DestinationTemplate.Query().Limit(limit).Skip(skip).Run()
	if err != nil {
		return nil, err
	}

	return list, nil
}
