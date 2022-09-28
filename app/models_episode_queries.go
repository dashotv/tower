package app

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *Connector) Upcoming() ([]*Episode, error) {
	// TODO: add series counts check
	seriesMap := map[string]*Series{}
	// Get Active Series
	series, err := c.SeriesActive()
	if err != nil {
		return nil, err
	}

	// Create a slice of ids
	ids := make([]primitive.ObjectID, len(series))
	for _, s := range series {
		ids = append(ids, s.ID)
		if seriesMap[s.ID.Hex()] == nil {
			seriesMap[s.ID.Hex()] = s
		}
	}

	// Get upcoming episodes
	q2 := c.Episode.Query()
	now := time.Now()
	since := time.Now().Add(-time.Hour * 24 * 7)
	//App().Log.Println("ids count ", len(ids))
	//App().Log.Println("time between ", since, " and ", now)
	list, err := q2.
		Where("_type", "Episode").
		In("series_id", ids).
		Where("downloaded", false).
		Where("completed", false).
		Where("skipped", false).
		In("missing", []interface{}{false, nil}).
		GreaterThan("season_number", 0).
		GreaterThan("episode_number", 0).
		LessThanEqual("release_date", now).
		GreaterThanEqual("release_date", since).
		Asc("release_date").Asc("season_number").Asc("episode_number").
		Limit(25).
		Run()
	if err != nil {
		return nil, err
	}

	// Copy the paths (images) from Series to Episode
	for _, e := range list {
		sid := e.SeriesId.Hex()
		if seriesMap[sid] != nil {
			e.Paths = seriesMap[sid].Paths
		}
	}

	return list, nil
}
