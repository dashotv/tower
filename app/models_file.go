package app

import "fmt"

func (c *Connector) FileGet(id string) (*File, error) {
	m := &File{}
	err := c.File.Find(id, m)
	if err != nil {
		return nil, err
	}

	// post process here

	return m, nil
}

func (c *Connector) FileList() ([]*File, error) {
	list, err := c.File.Query().Limit(10).Run()
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (c *Connector) FileCount() (int64, error) {
	return c.File.Query().Count()
}

func (c *Connector) FileByPath(path string) (*File, error) {
	list, err := c.File.Query().Where("path", path).Run()
	if err != nil {
		return nil, err
	}

	if len(list) > 1 {
		return nil, fmt.Errorf("more than one file found for path: %s", path)
	}

	if len(list) == 0 {
		return &File{Path: path}, nil
	}

	return list[0], nil
}
