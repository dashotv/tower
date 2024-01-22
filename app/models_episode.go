package app

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/maps"

	"github.com/dashotv/grimoire"
)

const imagesBaseURL = "/media-images" // proxy this instead of dealing with CORS

func (c *Connector) UpcomingQuery() *grimoire.QueryBuilder[*Episode] {
	return c.Episode.Query().
		Where("downloaded", false).
		Where("completed", false).
		Where("skipped", false).
		In("missing", []interface{}{false, nil}).
		GreaterThan("season_number", 0).
		GreaterThan("episode_number", 0).
		Asc("release_date").Asc("season_number").Asc("episode_number")
}

func (c *Connector) Upcoming() ([]*Episode, error) {
	utc := time.Now().UTC()
	today := time.Date(utc.Year(), utc.Month(), utc.Day(), 0, 0, 0, 0, time.UTC)
	later := today.Add(time.Hour * 24 * 90)
	q := c.UpcomingQuery().
		GreaterThanEqual("release_date", today).
		LessThanEqual("release_date", later).
		Limit(-1)
	return c.UpcomingFrom(q)
}

func (c *Connector) UpcomingNow() ([]*Episode, error) {
	utc := time.Now().UTC()
	today := time.Date(utc.Year(), utc.Month(), utc.Day(), 0, 0, 0, 0, time.UTC)
	tomorrow := today.Add(time.Hour * 24)
	q := c.UpcomingQuery().
		GreaterThanEqual("release_date", today.Add(-30*time.Hour*24)).
		LessThan("release_date", tomorrow).
		Limit(-1)
	return c.UpcomingFrom(q)
}

func (c *Connector) UpcomingFrom(query *grimoire.QueryBuilder[*Episode]) ([]*Episode, error) {
	seriesMap := map[primitive.ObjectID]*Series{}
	list, err := query.Run()
	if err != nil {
		return nil, err
	}

	c.Log.Debugf("upcoming: %d", len(list))
	list = groupEpisodes(list)

	// Create a slice of ids
	sids := make([]primitive.ObjectID, 0)
	for _, e := range list {
		sids = append(sids, e.SeriesId)
	}

	series, err := c.Series.Query().In("_id", sids).Limit(-1).Run()
	if err != nil {
		return nil, err
	}

	for _, s := range series {
		if seriesMap[s.ID] == nil {
			seriesMap[s.ID] = s
		}
	}

	for _, e := range list {
		sid := e.SeriesId
		if seriesMap[sid] != nil {
			c.processSeriesEpisode(seriesMap[sid], e)
		}
	}

	c.Log.Infof("episodes %d sids %d series %d seriesmap %d", len(list), len(sids), len(series), len(maps.Keys(seriesMap)))
	return list, nil
}

func (c *Connector) SeriesDownloadCounts() (map[string]int, error) {
	counts := map[string]int{}

	list, err := c.ActiveDownloads()
	if err != nil {
		return nil, err
	}

	for _, d := range list {
		m := &Medium{}
		err := c.Medium.Find(d.MediumId.Hex(), m)
		if err != nil {
			return nil, err
		}
		if m.Type == "Episode" {
			counts[m.SeriesId.Hex()]++
		}
	}

	return counts, nil
}

func (c *Connector) EpisodeGet(id string) (*Episode, error) {
	e := &Episode{}
	err := c.Episode.Find(id, e)
	if err != nil {
		return nil, err
	}
	c.processEpisode(e)
	return e, nil
}

func (c *Connector) processEpisode(e *Episode) error {
	s := &Series{}
	err := c.Series.Find(e.SeriesId.Hex(), s)
	if err != nil {
		return errors.Wrap(err, "processEpisode")
	}

	c.processSeriesEpisode(s, e)
	return nil
}

func (c *Connector) processSeriesEpisode(s *Series, e *Episode) {
	if s.Kind == "anime" {
		e.Display = fmt.Sprintf("#%d %s", e.AbsoluteNumber, e.Title)
	} else {
		e.Display = fmt.Sprintf("%02dx%02d %s", e.SeasonNumber, e.EpisodeNumber, e.Title)
	}
	unwatched, err := c.SeriesUserUnwatched(s)
	if err != nil {
		c.Log.Errorf("getting unwatched %s: %s", s.ID.Hex(), err)
	}
	e.Unwatched = unwatched
	e.Directory = s.Directory
	e.Active = s.Active
	e.Favorite = s.Favorite
	e.Title = s.Title
	e.Kind = s.Kind
	e.Source = s.Source
	e.SourceId = s.SourceId
	for _, p := range s.Paths {
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

func (c *Connector) EpisodePaths(id string) ([]*Path, error) {
	e := &Episode{}
	err := c.Episode.Find(id, e)
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
