package app

import (
	"time"
)

func (c *Connector) Upcoming() ([]*Episode, error) {
	// TODO: add series counts check
	// Get Active Series
	series, err := c.SeriesActive()
	if err != nil {
		return nil, err
	}

	// Create a slice of ids
	ids := make([]string, len(series))
	for _, s := range series {
		ids = append(ids, s.ID.Hex())
	}

	// Get upcoming episodes
	q2 := c.Episode.Query()
	now := time.Now()
	since := time.Now().Add(-time.Hour * 24 * 7)
	//fmt.Println("time between ", since, " and ", now)
	return q2.Where("_type", "Episode").
		In("series_id", ids).
		Where("downloaded", false).
		Where("completed", false).
		Where("skipped", false).
		//In("missing", []interface{}{false, nil}).
		GreaterThan("season_number", 0).
		GreaterThan("episode_number", 0).
		LessThanEqual("release_date", now).
		GreaterThanEqual("release_date", since).
		Asc("release_date").Asc("season_number").Asc("episode_number").
		Limit(25).
		Run()
}
