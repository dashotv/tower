package app

import (
	"context"
	"sort"
	"strconv"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *Connector) SeriesActive() ([]*Series, error) {
	return c.Series.Query().
		Where("active", true).
		Limit(-1).
		Run()
}

func (c *Connector) SeriesAll() ([]*Series, error) {
	return c.Series.Query().
		Limit(-1).
		Run()
}

func (c *Connector) SeriesUnwatchedByID(id string) (int, error) {
	s := &Series{}
	err := c.Series.Find(id, s)
	if err != nil {
		return 0, err
	}

	return c.SeriesUnwatched(s, "")
}

func (c *Connector) SeriesUnwatched(s *Series, user string) (int, error) {
	list, err := c.Episode.Query().
		Where("series_id", s.ID).
		GreaterThan("season_number", 0).
		Where("completed", true).
		Limit(-1).
		Run()
	if err != nil {
		return 0, err
	}

	grouped := lo.GroupBy(list, func(e *Episode) primitive.ObjectID {
		return e.ID
	})
	ids := lo.Keys(grouped)
	ids = lo.Uniq[primitive.ObjectID](ids)

	// get watches for those ids
	q := c.Watch.Query()
	if user != "" {
		q = q.Where("username", user)
	}
	watches, err := q.In("medium_id", ids).Limit(-1).Run()
	if err != nil {
		return 0, err
	}

	grpwatches := lo.GroupBy(watches, func(e *Watch) primitive.ObjectID {
		return e.ID
	})
	wids := lo.Keys(grpwatches)
	wids = lo.Uniq[primitive.ObjectID](wids)

	// return total episodes - total watches
	return len(ids) - len(wids), nil
}

func (c *Connector) SeriesUserUnwatched(s *Series) (int, error) {
	return c.SeriesUnwatched(s, "xenonsoul")
}

// func (c *Connector) SeriesUserUnwatchedCached(s *Series) (int, error) {
// 	unwatched := 0
// 	_, err := app.Cache.Fetch(fmt.Sprintf("series-unwatched-%s", s.ID.Hex()), &unwatched, func() (interface{}, error) {
// 		return c.SeriesUserUnwatched(s)
// 	})
// 	if err != nil {
// 		return 0, err
// 	}
//
// 	return unwatched, nil
// }

func (c *Connector) SeriesSeasons(id string) ([]int, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	c.Log.Infof("seasons: oid=%s", oid)

	col := c.Episode.Collection
	results, err := col.Distinct(context.TODO(), "season_number", bson.D{{Key: "series_id", Value: oid}})
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return []int{1}, nil
	}

	var out []int
	for _, r := range results {
		out = append(out, int(r.(int32)))
	}
	sort.Ints(out)

	return out, nil
}

func (c *Connector) SeriesSeasonEpisodes(id string, season string) ([]*Episode, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	s, err := strconv.Atoi(season)
	if err != nil {
		return nil, err
	}

	q := c.Episode.Query()
	eps, err := q.
		Where("series_id", oid).
		Where("season_number", s).
		Asc("episode_number").
		Limit(-1).
		Run()
	if err != nil {
		return nil, err
	}

	for _, e := range eps {
		e.Watched = c.MediumWatched(e.ID)
	}

	return eps, nil
}

func (c *Connector) SeriesSeasonEpisodesAll(id string) ([]*Episode, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	q := c.Episode.Query()
	eps, err := q.
		Where("series_id", oid).
		Asc("season_number").
		Asc("episode_number").
		Asc("absolute_number").
		Limit(1000).
		Run()
	if err != nil {
		return nil, err
	}

	return eps, nil
}

func (c *Connector) SeriesSetting(id, setting string, value bool) error {
	s := &Series{}
	err := c.Series.Find(id, s)
	if err != nil {
		return err
	}

	c.Log.Infof("series setting: %s %t", setting, value)
	switch setting {
	case "active":
		s.Active = value
	case "favorite":
		s.Favorite = value
	case "broken":
		s.Broken = value
	}

	return c.Series.Update(s)
}

func (c *Connector) SeriesUpdate(id string, data *Series) error {
	s := &Series{}
	err := c.Series.Find(id, s)
	if err != nil {
		return err
	}

	s.Display = data.Display
	s.Directory = data.Directory
	s.Kind = data.Kind
	s.Source = data.Source
	s.SourceId = data.SourceId
	s.Search = data.Search
	s.SearchParams = data.SearchParams

	return c.Series.Update(s)
}

func (c *Connector) SeriesCurrentSeason(id string) (int, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return -1, err
	}

	eps, err := c.Episode.Query().
		Where("series_id", oid).
		GreaterThan("season_number", 0).
		Asc("season_number").Asc("episode_number").Asc("absolute_number").
		Where("completed", false).Where("skipped", false).
		Limit(1).
		Run()
	if err != nil {
		return -1, err
	}

	if len(eps) > 0 && eps[0].SeasonNumber != 0 {
		return eps[0].SeasonNumber, nil
	}

	seasons, err := c.SeriesSeasons(id)
	if err != nil {
		return -1, err
	}

	return seasons[len(seasons)-1], nil
}

func (c *Connector) SeriesPaths(id string) ([]*Path, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	s := &Series{}
	err = c.Series.FindByID(oid, s)
	if err != nil {
		return nil, err
	}

	var out []*Path
	out = append(out, s.Paths...)

	eps, err := c.Episode.Query().
		Where("series_id", oid).
		Desc("season_number").Desc("episode_number").Desc("absolute_number").
		Limit(5000).
		Run()
	if err != nil {
		return nil, err
	}

	for _, e := range eps {
		if len(e.Paths) > 0 {
			out = append(out, e.Paths...)
		}
	}

	return out, nil
}

func (c *Connector) SeriesWatches(id string) ([]*Watch, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	s := &Series{}
	err = c.Series.FindByID(oid, s)
	if err != nil {
		return nil, err
	}

	eps, err := c.Episode.Query().
		Where("series_id", oid).
		Desc("season_number").Desc("episode_number").Desc("absolute_number").
		Limit(5000).
		Run()
	if err != nil {
		return nil, err
	}

	// get ids of all episodes from query above
	ids := []primitive.ObjectID{}
	for _, e := range eps {
		ids = append(ids, e.ID)
	}
	// get distinct list of ids
	ids = lo.Uniq[primitive.ObjectID](ids)

	// get watches for those ids
	watches, err := c.Watch.Query().Desc("watched_at").In("medium_id", ids).Limit(len(ids)).Run()
	if err != nil {
		return nil, err
	}

	media, err := c.Medium.Query().In("_id", ids).Limit(len(ids)).Run()
	if err != nil {
		return nil, err
	}

	mmap := map[primitive.ObjectID]*Medium{}
	for _, m := range media {
		//c.log.Infof("medium %s: %#v", m.ID.Hex(), m)
		mmap[m.ID] = m
	}

	for _, w := range watches {
		//c.log.Infof("watch %s: %#v", w.MediumId.Hex(), w.MediumId)
		w.Medium = mmap[w.MediumId]
	}

	return watches, nil
}
