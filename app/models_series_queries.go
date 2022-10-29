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

func (c *Connector) SeriesSeasons(id string) ([]string, error) {
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
		return []string{"1"}, nil
	}

	var out []string
	for _, r := range results {
		App().Log.Infof("seasons: result=%v", r)
		out = append(out, fmt.Sprintf("%v", r))
	}
	sort.Strings(out)

	return out, nil
}

func (c *Connector) SeriesSeasonEpisodes(id string, season string) ([]*Episode, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	s, err := strconv.Atoi(season)

	q := c.Episode.Query()
	return q.
		Where("_type", "Episode").
		Where("series_id", oid).
		Where("season_number", s).
		Asc("episode_number").
		Limit(1000).
		Run()
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

	i, err := strconv.Atoi(seasons[len(seasons)-1])
	if err != nil {
		return -1, err
	}

	return i, nil
}
