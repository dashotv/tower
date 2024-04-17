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

func (s *Series) GetCover() *Path {
	for _, p := range s.Paths {
		if p.Type == "cover" {
			return p
		}
	}
	return nil
}

func (s *Series) GetBackground() *Path {
	for _, p := range s.Paths {
		if p.Type == "background" {
			return p
		}
	}
	return nil
}

func (c *Connector) processSeries(s *Series) {
	for _, p := range s.Paths {
		if p.Type == "cover" {
			s.Cover = fmt.Sprintf("%s/%s.%s", imagesBaseURL, p.Local, p.Extension)
			continue
		}
		if p.Type == "background" {
			s.Background = fmt.Sprintf("%s/%s.%s", imagesBaseURL, p.Local, p.Extension)
			continue
		}
	}
}

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
	// get all episodes for series
	eps, err := c.Episode.Query().
		Where("series_id", s.ID).
		Where("skipped", false).
		Where("completed", true).
		GreaterThan("season_number", 0).
		Limit(-1).
		Run()
	if err != nil {
		return 0, err
	}

	// get ids of all episodes who have videos
	grouped := lo.GroupBy(eps, func(e *Episode) primitive.ObjectID {
		for _, p := range e.Paths {
			if p.Type == "video" {
				return e.ID
			}
		}
		return primitive.NilObjectID
	})
	ids := lo.Keys(grouped)
	ids = lo.Filter(ids, func(id primitive.ObjectID, i int) bool {
		return id != primitive.NilObjectID
	})
	ids = lo.Uniq[primitive.ObjectID](ids)

	// get unique watches for those ids
	q := c.Watch.Query()
	if user != "" {
		q = q.Where("username", user)
	}
	watches, err := q.In("medium_id", ids).Limit(-1).Run()
	if err != nil {
		return 0, err
	}

	grpwatches := lo.GroupBy(watches, func(e *Watch) primitive.ObjectID {
		return e.MediumID
	})
	wids := lo.Keys(grpwatches)
	wids = lo.Uniq[primitive.ObjectID](wids)

	// return total episodes - total watches
	return len(ids) - len(wids), nil
}

func (c *Connector) SeriesUserUnwatched(s *Series) (int, error) {
	return c.SeriesUnwatched(s, app.Config.PlexUsername)
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
	series, err := c.Series.Get(id, &Series{})
	if err != nil {
		return nil, err
	}

	s, err := strconv.Atoi(season)
	if err != nil {
		return nil, err
	}

	q := c.Episode.Query().
		Asc("season_number").Asc("episode_number").Asc("absolute_number").
		Where("series_id", series.ID)
	if !isAnimeKind(string(series.Kind)) {
		q = q.Where("season_number", s)
	}
	eps, err := q.
		Limit(-1).
		Run()
	if err != nil {
		return nil, err
	}

	eids := map[primitive.ObjectID]*Episode{}
	for _, e := range eps {
		eids[e.ID] = e
	}

	watches, err := c.Watch.Query().In("medium_id", lo.Keys(eids)).Limit(-1).Run()
	if err != nil {
		return nil, err
	}

	for _, w := range watches {
		e, ok := eids[w.MediumID]
		if ok {
			if w.Username == app.Config.PlexUsername {
				e.Watched = true
			} else {
				e.WatchedAny = true
			}
		}
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
		Limit(-1).
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
	s.SourceID = data.SourceID
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
		//c.log.Infof("watch %s: %#v", w.MediumID.Hex(), w.MediumID)
		w.Medium = mmap[w.MediumID]
	}

	return watches, nil
}

func (c *Connector) SeriesBySearch(title string) (*Series, error) {
	title = path(title)
	{
		list, err := c.Series.Query().Where("kind", "donghua").Where("directory", title).Run()
		if err != nil {
			return nil, err
		}
		if len(list) == 1 {
			return list[0], nil
		}
	}

	{
		list, err := c.Series.Query().Where("kind", "donghua").Where("search", title).Run()
		if err != nil {
			return nil, err
		}
		if len(list) == 1 {
			return list[0], nil
		}
	}

	return nil, nil
}

func (c *Connector) SeriesEpisodeBy(s *Series, season, episode int) (*Episode, error) {
	{
		list, err := c.Episode.Query().Where("series_id", s.ID).Where("completed", false).Where("downloaded", false).Where("skipped", false).Where("season_number", season).Where("episode_number", episode).Run()
		if err != nil {
			return nil, err
		}
		// c.Log.Debugf("season/episode: %d/%d: %d", season, episode, len(list))
		if len(list) == 1 {
			return list[0], nil
		}
	}

	{
		list, err := c.Episode.Query().Where("series_id", s.ID).Where("completed", false).Where("downloaded", false).Where("skipped", false).Where("absolute_number", episode).Run()
		if err != nil {
			return nil, err
		}
		// c.Log.Debugf("season/episode (abs): %d/%d: %d", season, episode, len(list))
		if len(list) == 1 {
			return list[0], nil
		}
	}

	return nil, nil
}
