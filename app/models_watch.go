package app

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *Connector) MediumWatched(id primitive.ObjectID) bool {
	// TODO: add user name to config
	watches, _ := c.Watch.Query().Where("medium_id", id).Where("username", "xenonsoul").Run()
	return len(watches) > 0
}
func (c *Connector) MediumWatchedAny(id primitive.ObjectID) bool {
	// TODO: add user name to config
	watches, _ := c.Watch.Query().Where("medium_id", id).Run()
	return len(watches) > 0
}

func (c *Connector) Watches(mediumId, username string) ([]*Watch, error) {
	query := c.Watch.Query().Limit(100).Desc("watched_at")
	if username != "" {
		query = query.Where("username", username)
	}

	if mediumId != "" {
		id, err := primitive.ObjectIDFromHex(mediumId)
		if err != nil {
			return nil, err
		}

		m := &Medium{}
		if err := c.Medium.Find(mediumId, m); err != nil {
			return nil, err
		}

		if m.Type == "Series" {
			episodes, err := c.Episode.Query().
				Where("series_id", m.ID).
				Desc("episode_number").Desc("series_number").
				Limit(-1).Run()
			if err != nil {
				return nil, err
			}

			ids := make([]primitive.ObjectID, len(episodes))
			for i, e := range episodes {
				ids[i] = e.ID
			}

			query = query.Limit(len(ids)).In("medium_id", ids)
		} else {
			query = query.Where("medium_id", id)
		}
	}

	watches, err := query.Run()
	if err != nil {
		return nil, err
	}

	for _, w := range watches {
		m := &Medium{}
		if err := c.Medium.FindByID(w.MediumId, m); err != nil {
			return nil, err
		}
		w.Medium = m
		if m.Type == "Episode" {
			s := &Series{}
			if err := c.Series.FindByID(m.SeriesId, s); err != nil {
				return nil, err
			}

			if isAnimeKind(string(s.Kind)) {
				m.Display = fmt.Sprintf("%02dx%02d #%03d %s", m.SeasonNumber, m.EpisodeNumber, m.AbsoluteNumber, m.Title)
			} else {
				m.Display = fmt.Sprintf("%02dx%02d %s", m.SeasonNumber, m.EpisodeNumber, m.Title)
			}
			m.Title = s.Display
		}
	}

	return watches, nil
}
