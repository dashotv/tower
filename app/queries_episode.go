package app

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/maps"
)

const imagesBaseURL = "/media-images" // proxy this instead of dealing with CORS

func (c *Connector) Upcoming() ([]*Episode, error) {
	seriesMap := map[primitive.ObjectID]*Series{}
	utc := time.Now().UTC()
	today := time.Date(utc.Year(), utc.Month(), utc.Day(), 0, 0, 0, 0, time.UTC)
	later := today.Add(time.Hour * 24 * 90)
	list, err := c.Episode.Query().
		Where("_type", "Episode").
		Where("downloaded", false).
		Where("completed", false).
		Where("skipped", false).
		In("missing", []interface{}{false, nil}).
		GreaterThan("season_number", 0).
		GreaterThan("episode_number", 0).
		GreaterThanEqual("release_date", today).
		LessThanEqual("release_date", later).
		Asc("release_date").Asc("season_number").Asc("episode_number").
		Limit(1000).
		Run()
	if err != nil {
		return nil, err
	}

	list = groupEpisodes(list)

	// Create a slice of ids
	sids := make([]primitive.ObjectID, 0)
	for _, e := range list {
		sids = append(sids, e.SeriesId)
	}

	series, err := c.Series.Query().Where("_type", "Series").In("_id", sids).Limit(-1).Run()
	if err != nil {
		return nil, err
	}

	for _, s := range series {
		if seriesMap[s.ID] == nil {
			seriesMap[s.ID] = s
		}
	}

	// Copy the paths (images) from Series to Episode
	for _, e := range list {
		sid := e.SeriesId
		if seriesMap[sid] != nil {
			//if seriesMap[sid].Type == "Anime" {
			//	e.Display = fmt.Sprintf("#%d %s", e.AbsoluteNumber, e.Title)
			//} else {
			unwatched, err := c.SeriesAllUnwatched(seriesMap[sid])
			if err != nil {
				c.log.Errorf("getting unwatched %s: %s", sid, err)
			}
			e.Unwatched = unwatched
			e.Title = seriesMap[sid].Title
			for _, p := range seriesMap[sid].Paths {
				if p.Type == "cover" {
					e.Cover = fmt.Sprintf("%s/%s.%s", imagesBaseURL, p.Local, p.Extension)
					continue
				}
				if p.Type == "background" {
					e.Background = fmt.Sprintf("%s/%s.%s", imagesBaseURL, p.Local, p.Extension)
					continue
				}
			}
		}
	}

	c.log.Infof("episodes %d sids %d series %d seriesmap %d", len(list), len(sids), len(series), len(maps.Keys(seriesMap)))
	return list, nil
}

func (c *Connector) EpisodeSetting(id, setting string, value bool) error {
	e := &Episode{}
	err := c.Episode.Find(id, e)
	if err != nil {
		return err
	}

	switch setting {
	case "downloaded":
		e.Downloaded = value
	case "skipped":
		e.Skipped = value
	case "completed":
		e.Completed = value
	}

	return c.Episode.Update(e)
}

func (c *Connector) EpisodePaths(id string) ([]Path, error) {
	e := &Episode{}
	err := db.Episode.Find(id, e)
	if err != nil {
		return nil, err
	}

	return e.Paths, nil
}

func groupEpisodes(list []*Episode) []*Episode {
	track := map[string]bool{}
	out := []*Episode{}

	for _, e := range list {
		sid := e.SeriesId.Hex()
		if !track[sid] {
			out = append(out, e)
			track[sid] = true
		}
	}

	return out
}
