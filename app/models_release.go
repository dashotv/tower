package app

var releaseTypes = []string{"tv", "anime", "movies"}

func (c *Connector) ReleasesAll() ([]*Release, error) {
	return c.Release.Query().Limit(-1).Run()
}

func (c *Connector) ReleaseSetting(id, setting string, value bool) error {
	release := &Release{}
	err := c.Release.Find(id, release)
	if err != nil {
		return err
	}

	switch setting {
	case "verified":
		release.Verified = value
	}

	return c.Release.Update(release)
}
