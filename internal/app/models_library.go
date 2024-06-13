package app

import "github.com/dashotv/fae"

func (c *Connector) LibraryGet(id string) (*Library, error) {
	m := &Library{}
	err := c.Library.Find(id, m)
	if err != nil {
		return nil, err
	}

	c.processLibraries([]*Library{m})

	return m, nil
}

func (c *Connector) LibraryGetByKind(kind string) (*Library, error) {
	list, err := c.Library.Query().Where("kind", kind).Run()
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, fae.Errorf("library not found for kind: %s", kind)
	}
	if len(list) > 1 {
		return nil, fae.Errorf("multiple libraries found for kind: %s", kind)
	}

	c.processLibraries(list)
	return list[0], nil
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

func (c *Connector) LibraryDestination(m *Medium) (string, error) {
	lib, err := c.LibraryGetByKind(string(m.Kind))
	if err != nil {
		return "", fae.Wrap(err, "failed to get library")
	}

	return lib.LibraryTemplate.Name, nil
}
