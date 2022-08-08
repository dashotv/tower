package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *MediumStore) Upcoming() ([]Medium, error) {
	q := s.Query()
	series, err := q.
		Where("_type", "Series").
		Where("active", true).
		Limit(10000).
		Run()
	if err != nil {
		return nil, err
	}

	ids := make([]primitive.ObjectID, len(series))
	for _, s := range series {
		ids = append(ids, s.ID)
	}

	q2 := s.Query()
	return q2.
		Where("_type", "Episode").
		In("series_id", ids).
		Where("downloaded", false).
		Where("completed", false).
		Where("skipped", false).
		Where("missing", false).
		GreaterThan("season_number", 0).
		GreaterThan("episode_number", 0).
		GreaterThan("release_date", time.Now()).
		Asc("release_date").Asc("season_number").Asc("episode_number").
		Run()
}

func (m *Medium) Background() string {
	for _, p := range m.Paths {
		if p.Type == "background" {
			return p.Local
		}
	}
	return ""
}

func (m *Medium) Cover() string {
	for _, p := range m.Paths {
		if p.Type == "cover" {
			return p.Local
		}
	}
	return ""
}
