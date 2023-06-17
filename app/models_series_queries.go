package app

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *Connector) SeriesActive() ([]*Series, error) {
	return c.Series.Query().
		Where("_type", "Series").
		Where("active", true).
		Limit(1000).
		Run()
}

func (c *Connector) SeriesAll() ([]*Series, error) {
	return c.Series.Query().
		Where("_type", "Series").
		Limit(5000).
		Run()
}

func (c *Connector) SeriesAllUnwatched(s *Series) (int, error) {
	list, err := c.Episode.Query().Where("series_id", s.ID).GreaterThan("season_number", 0).Where("completed", true).Limit(10).Run()
	if err != nil {
		return 0, err
	}

	// get ids of all episodes from query above
	ids := []primitive.ObjectID{}
	for _, e := range list {
		ids = append(ids, e.ID)
	}
	// get distinct list of ids
	ids = lo.Uniq[primitive.ObjectID](ids)

	// get watches for those ids
	watches, err := c.Watch.Query().Where("username", "xenonsoul").In("medium_id", ids).Run()
	if err != nil {
		return 0, err
	}

	// return total episodes - total watches
	return len(ids) - len(watches), nil
}

func (c *Connector) SeriesAllUnwatchedCached(s *Series) (int, error) {
	unwatched := 0
	_, err := App().Cache.Fetch(fmt.Sprintf("series-unwatched-%s", s.ID.Hex()), &unwatched, func() (interface{}, error) {
		return c.SeriesAllUnwatched(s)
	})
	if err != nil {
		return 0, err
	}

	return unwatched, nil
}

func (c *Connector) SeriesSeasons(id string) ([]int, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	App().Log.Infof("seasons: oid=%s", oid)

	col := c.Episode.Collection
	results, err := col.Distinct(context.TODO(), "season_number", bson.D{{"series_id", oid}})
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
		Where("_type", "Episode").
		Where("series_id", oid).
		Where("season_number", s).
		Asc("episode_number").
		Limit(1000).
		Run()

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
		Where("_type", "Episode").
		Where("series_id", oid).
		Asc("season_number").
		Asc("episode_number").
		Asc("absolute_number").
		Limit(1000).
		Run()

	return eps, nil
}

func (c *Connector) SeriesSetting(id, setting string, value bool) error {
	s := &Series{}
	err := c.Series.Find(id, s)
	if err != nil {
		return err
	}

	App().Log.Infof("series setting: %s %t", setting, value)
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

func (c *Connector) SeriesCurrentSeason(id string) (int, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return -1, err
	}

	eps, err := c.Episode.Query().
		Where("_type", "Episode").
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

func (c *Connector) SeriesPaths(id string) ([]Path, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	s := &Series{}
	err = App().DB.Series.FindByID(oid, s)
	if err != nil {
		return nil, err
	}

	var out []Path
	out = append(out, s.Paths...)

	eps, err := App().DB.Episode.Query().Where("_type", "Episode").
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
	err = App().DB.Series.FindByID(oid, s)
	if err != nil {
		return nil, err
	}

	eps, err := App().DB.Episode.Query().Where("_type", "Episode").
		Where("series_id", oid).
		Desc("season_number").Desc("episode_number").Desc("absolute_number").
		Limit(5000).
		Run()

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
		//App().Log.Infof("medium %s: %#v", m.ID.Hex(), m)
		mmap[m.ID] = m
	}

	for _, w := range watches {
		//App().Log.Infof("watch %s: %#v", w.MediumId.Hex(), w.MediumId)
		w.Medium = mmap[w.MediumId]
	}

	return watches, nil
}
