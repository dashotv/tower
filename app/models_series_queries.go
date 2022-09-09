package app

import "github.com/samber/lo"

func (c *Connector) SeriesActive() ([]*Series, error) {
	return c.Series.Query().
		Where("_type", "Series").
		Where("active", true).
		Limit(1000).
		Run()
}

func (c *Connector) SeriesAllUnwatched(s *Series) (int, error) {
	list, err := c.Episode.Query().GreaterThan("season_number", 0).Where("completed", true).Run()
	if err != nil {
		return 0, err
	}

	// get ids of all episodes from query above
	ids := []string{}
	for _, e := range list {
		ids = append(ids, e.ID.Hex())
	}
	// get distinct list of ids
	ids = lo.Uniq[string](ids)

	// get watches for those ids
	watches, err := c.Watch.Query().In("medium_id", ids).Run()
	if err != nil {
		return 0, err
	}

	// return total episodes - total watches
	return len(ids) - len(watches), nil
}
