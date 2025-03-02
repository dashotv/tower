package app

import (
	"fmt"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/samber/lo"

	"github.com/dashotv/fae"
	"github.com/dashotv/grimoire"
)

const imagesBaseURL = "/media-images" // proxy this instead of dealing with CORS

func (c *Connector) UpcomingQuery() *grimoire.QueryBuilder[*Episode] {
	series, err := c.Series.Query().Where("active", true).Limit(-1).Run()
	if err != nil {
		c.Log.Errorf("error getting series: %s", err)
		return nil
	}

	ids := lo.Map(series, func(s *Series, i int) primitive.ObjectID {
		return s.ID
	})

	return c.Episode.Query().
		In("series_id", ids).
		Where("downloaded", false).
		Where("completed", false).
		Where("skipped", false).
		In("missing", []interface{}{false, nil}).
		GreaterThan("season_number", 0).
		GreaterThan("episode_number", 0).
		Asc("release_date").Asc("season_number").Asc("episode_number").Asc("absolute_number")
}

func (c *Connector) Upcoming() ([]*Upcoming, error) {
	utc := time.Now().UTC()
	today := time.Date(utc.Year(), utc.Month(), utc.Day(), 0, 0, 0, 0, time.UTC)
	later := today.Add(time.Hour * 24 * 7)
	q := c.Episode.Query().
		Where("downloaded", false).
		Where("completed", false).
		Where("skipped", false).
		In("missing", []interface{}{false, nil}).
		GreaterThan("season_number", 0).
		GreaterThan("episode_number", 0).
		Asc("release_date").Asc("season_number").Asc("episode_number").Asc("absolute_number").
		GreaterThanEqual("release_date", today).
		LessThanEqual("release_date", later).
		Limit(-1)
	return c.UpcomingFrom(q)
}

func (c *Connector) UpcomingNow() ([]*Upcoming, error) {
	utc := time.Now().UTC()
	// null time breaks seer, so we set unknown time to unix epoch
	// we want to avoid including those in the upcoming list
	after := time.Date(1974, 1, 1, 0, 0, 0, 0, time.UTC)
	today := time.Date(utc.Year(), utc.Month(), utc.Day(), 0, 0, 0, 0, time.UTC)
	tomorrow := today.Add(time.Hour * 24)
	q := c.UpcomingQuery().
		GreaterThanEqual("release_date", after).
		LessThan("release_date", tomorrow).
		Limit(-1)
	return c.UpcomingFrom(q)
}

func (c *Connector) UpcomingLater() ([]*Upcoming, error) {
	utc := time.Now().UTC()
	today := time.Date(utc.Year(), utc.Month(), utc.Day(), 0, 0, 0, 0, time.UTC)
	// start := today.Add(time.Hour * 24 * 7)
	end := today.Add(time.Hour * 24 * 90)
	q := c.Episode.Query().
		Where("downloaded", false).
		Where("completed", false).
		Where("skipped", false).
		In("missing", []interface{}{false, nil}).
		GreaterThan("season_number", 0).
		GreaterThan("episode_number", 0).
		Asc("release_date").Asc("season_number").Asc("episode_number").Asc("absolute_number").
		GreaterThanEqual("release_date", today).
		LessThanEqual("release_date", end).
		Limit(-1)
	return c.UpcomingFrom(q)
}

func (c *Connector) UpcomingFrom(query *grimoire.QueryBuilder[*Episode]) ([]*Upcoming, error) {
	// defer TickTock("UpcomingFrom")()
	seriesMap := map[primitive.ObjectID]*Series{}

	list, err := query.Run()
	if err != nil {
		return nil, fae.Wrap(err, "running query")
	}

	// c.Log.Debugf("upcoming: %d", len(list))
	list = groupEpisodes(list)

	// Create a slice of series ids
	sids := lo.Map(list, func(e *Episode, i int) primitive.ObjectID {
		return e.SeriesID
	})
	sids = lo.Uniq(sids)
	// c.Log.Debugf("upcoming series: %d", len(sids))

	// create map of series id to series
	err = c.Series.Query().In("_id", sids).Batch(100, func(results []*Series) error {
		for _, s := range results {
			seriesMap[s.ID] = s
		}
		return nil
	})
	if err != nil {
		return nil, fae.Wrap(err, "getting series")
	}
	// c.Log.Debugf("upcoming series map: %d", len(seriesMap))

	upcoming := lo.Map(list, func(e *Episode, i int) *Upcoming {
		if seriesMap[e.SeriesID] == nil {
			c.Log.Errorf("series not found %s", e.SeriesID.Hex())
			notifier.log("error", "upcoming", fmt.Sprintf("series not found %s (e:%s) %s", e.SeriesID.Hex(), e.ID.Hex(), e.Title))
			return nil
		}
		return c.processSeriesUpcoming(seriesMap[e.SeriesID], e)
	})
	upcoming = lo.Compact(upcoming)

	// c.Log.Infof("episodes %d sids %d series %d seriesmap %d", len(list), len(sids), len(series), len(maps.Keys(seriesMap)))
	return upcoming, nil
}

func (c *Connector) SeriesDownloadCounts() (map[string]int, error) {
	counts := map[string]int{}

	list, err := c.ActiveDownloads()
	if err != nil {
		return nil, err
	}

	for _, d := range list {
		if d.Medium == nil {
			c.Log.Warnf("missing medium %s", d.ID.Hex())
			continue
		}
		if d.Medium.Type == "Episode" {
			counts[d.Medium.SeriesID.Hex()]++
		}
		if d.Medium.Type == "Series" {
			counts[d.Medium.ID.Hex()]++
		}
	}

	return counts, nil
}
func (c *Connector) SeriesMultiDownloads() (map[string]bool, error) {
	out := map[string]bool{}

	list, err := c.ActiveDownloads()
	if err != nil {
		return nil, err
	}

	for _, d := range list {
		if d.Medium == nil {
			c.Log.Warnf("missing medium %s", d.ID.Hex())
			continue
		}
		if d.Medium.Type == "Episode" {
			out[d.Medium.SeriesID.Hex()] = true
		}
		if d.Medium.Type == "Series" {
			out[d.Medium.ID.Hex()] = true
		}
	}

	return out, nil
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
	err := c.Series.Find(e.SeriesID.Hex(), s)
	if err != nil {
		return fae.Wrap(err, "processEpisode")
	}

	c.processSeriesEpisode(s, e)
	return nil
}

func (c *Connector) processSeriesUpcoming(s *Series, e *Episode) *Upcoming {
	u := &Upcoming{
		ID:             e.ID,
		Type:           "Episode",
		Title:          s.Title,
		SourceID:       e.SourceID,
		Description:    e.Description,
		SeasonNumber:   e.SeasonNumber,
		EpisodeNumber:  e.EpisodeNumber,
		AbsoluteNumber: e.AbsoluteNumber,
		ReleaseDate:    e.ReleaseDate,
		Skipped:        e.Skipped,
		Downloaded:     e.Downloaded,
		Completed:      e.Completed,

		SeriesID:       s.ID,
		SeriesKind:     s.Kind,
		SeriesSource:   s.Source,
		SeriesTitle:    s.Title,
		Directory:      s.Directory,
		SeriesActive:   s.Active,
		SeriesFavorite: s.Favorite,
	}

	if isAnimeKind(string(s.Kind)) {
		u.Display = fmt.Sprintf("#%d %s", e.AbsoluteNumber, e.Title)
	} else {
		u.Display = fmt.Sprintf("%02dx%02d %s", e.SeasonNumber, e.EpisodeNumber, e.Title)
	}

	unwatched, err := c.SeriesUnwatched(s, "")
	if err != nil {
		c.Log.Errorf("getting unwatched %s: %s", s.ID.Hex(), err)
	}
	u.SeriesUnwatched = unwatched

	for _, p := range s.Paths {
		if p.Type == "cover" {
			u.SeriesCover = fmt.Sprintf("%s/%s.%s", imagesBaseURL, p.Local, p.Extension)
			u.SeriesCoverUpdated = p.UpdatedAt
			continue
		}
		if p.Type == "background" {
			u.SeriesBackground = fmt.Sprintf("%s/%s.%s", imagesBaseURL, p.Local, p.Extension)
			u.SeriesBackgroundUpdated = p.UpdatedAt
			continue
		}
	}

	return u
}
func (c *Connector) processSeriesEpisode(s *Series, e *Episode) {
	e.ApplyOverrides()
	if isAnimeKind(string(s.Kind)) {
		e.SeriesDisplay = fmt.Sprintf("#%d %s", e.AbsoluteNumber, e.Title)
	} else {
		e.SeriesDisplay = fmt.Sprintf("%02dx%02d %s", e.SeasonNumber, e.EpisodeNumber, e.Title)
	}

	unwatched, err := c.SeriesUnwatched(s, "")
	if err != nil {
		c.Log.Errorf("getting unwatched %s: %s", s.ID.Hex(), err)
	}
	e.SeriesUnwatched = unwatched

	e.SeriesActive = s.Active
	e.SeriesFavorite = s.Favorite
	e.Title = s.Display
	if e.Title == "" {
		e.Title = s.Title
	}
	e.SeriesKind = s.Kind
	e.SeriesSource = s.Source
	e.SourceID = s.SourceID
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

func (e *Episode) ApplyOverrides() {
	if e.Overrides == nil {
		return
	}
	a := e.Overrides.Absolute()
	if a >= 0 {
		e.HasOverrides = true
		e.AbsoluteNumber = a
	}
	s := e.Overrides.Season()
	if s >= 0 {
		e.HasOverrides = true
		e.SeasonNumber = s
	}
	ep := e.Overrides.Episode()
	if ep >= 0 {
		e.HasOverrides = true
		e.EpisodeNumber = ep
	}
}

func (o *Overrides) Episode() int {
	if o == nil {
		return -1
	}
	if o.EpisodeNumber != "" {
		num, _ := strconv.Atoi(o.EpisodeNumber)
		return num
	}
	return -1
}
func (o *Overrides) Season() int {
	if o == nil {
		return -1
	}
	if o.SeasonNumber != "" {
		num, _ := strconv.Atoi(o.SeasonNumber)
		return num
	}
	return -1
}
func (o *Overrides) Absolute() int {
	if o == nil {
		return -1
	}
	if o.AbsoluteNumber != "" {
		num, _ := strconv.Atoi(o.AbsoluteNumber)
		return num
	}
	return -1
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
		sid := e.SeriesID.Hex()
		if !track[sid] {
			out = append(out, e)
			track[sid] = true
		}
	}

	return out
}
