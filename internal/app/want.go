package app

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"github.com/samber/lo"

	"github.com/dashotv/fae"
	"github.com/dashotv/grimoire"
	runic "github.com/dashotv/runic/client"
)

var seriesWantedBuffer = 3

func init() {
	starters = append(starters, startWant)
}

func startWant(ctx context.Context, a *Application) error {
	go func() {
		defer TickTock("want")()
		w := NewWant(a.DB, a.Log.Named("want"))
		w.groups = a.Config.DownloadsGroups
		w.preferred = a.Config.DownloadsPreferred
		if err := w.Build(); err != nil {
			a.Log.Errorf("want build error: %s", err)
		}

		a.Want = w
	}()
	return nil
}

func NewWant(db *Connector, log *zap.SugaredLogger) *Want {
	return &Want{
		db:               db,
		log:              log,
		series_ids:       []primitive.ObjectID{},
		movies:           map[string]*Medium{},
		series_unwatched: map[string]int{},
		series_titles:    map[string]string{},
		series_episodes:  map[string][]*Medium{},
	}
}

type Want struct {
	db  *Connector
	log *zap.SugaredLogger

	series_ids       []primitive.ObjectID
	preferred        []string
	groups           []string
	series_unwatched map[string]int
	series_downloads map[string]int
	series_titles    map[string]string
	series_episodes  map[string][]*Medium
	movies           map[string]*Medium
}

func (w *Want) Release(r *runic.Release) *Medium {
	if r == nil {
		return nil
	}

	if r.Title == "" {
		return nil
	}

	// if r.Source != "rift" && !lo.Contains(w.preferred, r.Group) && !lo.Contains(w.groups, r.Group) {
	// 	return nil
	// }
	if !r.Verified {
		return nil
	}

	switch r.Type {
	case "movies":
		return w.releaseMovie(r)
	case "tv", "anime":
		return w.releaseEpisode(r)
	default:
		return nil
	}
}

func (w *Want) releaseMovie(release *runic.Release) *Medium {
	title := path(release.Title)
	// w.log.Debugf("movie %s", title)
	m, ok := w.movies[title]
	if !ok {
		return nil
	}
	if release.Year == 0 || m.ReleaseDate.Year() != release.Year {
		return nil
	}
	r, err := strconv.Atoi(release.Resolution)
	if err != nil {
		return nil
	}
	if r == 0 || m.SearchParams.Resolution != r {
		return nil
	}
	return m
}

func (w *Want) releaseEpisode(release *runic.Release) *Medium {
	seriesTitle := path(release.Title)
	series, ok := w.series_titles[seriesTitle]
	if !ok {
		return nil
	}

	r, err := strconv.Atoi(release.Resolution)
	if err != nil {
		return nil
	}
	if r != 1080 { // HACK: fix this
		return nil
	}

	for _, e := range w.series_episodes[series] {
		if e.SeasonNumber == release.Season && e.EpisodeNumber == release.Episode {
			return e
		}
		if e.AbsoluteNumber == release.Episode {
			return e
		}
	}
	return nil
}

func (w *Want) SeriesWanted(seriesID string) (*Wanted, error) {
	names := []string{}
	// we loop so we don't have to load from db
	for t, id := range w.series_titles {
		if id == seriesID {
			names = append(names, t)
		}
	}

	eps := lo.Map(w.series_episodes[seriesID], func(e *Medium, _ int) string {
		return fmt.Sprintf("%02dx%02d #%03d %s", e.SeasonNumber, e.EpisodeNumber, e.AbsoluteNumber, e.Title)
	})

	wanted := &Wanted{
		Names:    names,
		Episodes: eps,
	}

	return wanted, nil
}

func (w *Want) MovieWanted(movieID string) (*Wanted, error) {
	names := []string{}
	// we loop so we don't have to load from db
	for t, m := range w.movies {
		if m.ID.Hex() == movieID {
			names = append(names, t)
		}
	}

	wanted := &Wanted{
		Names: names,
	}

	return wanted, nil
}

func (w *Want) NextEpisode(seriesID string) *Medium {
	list, ok := w.series_episodes[seriesID]
	if !ok {
		return nil
	}
	if len(list) == 0 {
		return nil
	}

	return list[0]
}

func (w *Want) addSeries(s *Series) error {
	// w.log.Debugf("addSeries: %s: %s\n", s.ID.Hex(), s.Title)
	unwatched, err := w.db.SeriesUnwatchedByID(s.ID.Hex())
	if err != nil {
		return err
	}
	if unwatched >= seriesWantedBuffer {
		return nil
	}

	sid := s.ID.Hex()
	w.series_ids = append(w.series_ids, s.ID)
	w.series_unwatched[sid] = unwatched

	if s.Directory != "" {
		w.series_titles[s.Directory] = sid
	}
	if s.Search != "" {
		parts := strings.Split(s.Search, ":")
		if len(parts) > 0 && parts[0] != s.Directory {
			w.series_titles[parts[0]] = sid
		}
	}
	return nil
}

func (w *Want) addMovie(m *Medium) {
	if m.Kind != "movies" {
		return
	}
	if m.Directory != "" {
		w.movies[m.Directory] = m
	}
	if m.Search != "" && m.Search != m.Directory {
		w.movies[m.Search] = m
	}
}

func (w *Want) addEpisode(e *Medium) {
	// if len(w.series_episodes[e.SeriesID.Hex()]) < seriesWantedBuffer {
	if len(w.series_episodes[e.SeriesID.Hex()]) >= (seriesWantedBuffer - w.series_unwatched[e.SeriesID.Hex()]) {
		return
	}
	w.series_episodes[e.SeriesID.Hex()] = append(w.series_episodes[e.SeriesID.Hex()], e)
}

func (w *Want) Build() error {
	utc := time.Now().UTC()
	start := time.Date(1974, 1, 1, 0, 0, 0, 0, time.UTC)
	today := time.Date(utc.Year(), utc.Month(), utc.Day(), 0, 0, 0, 0, time.UTC)
	later := today.Add(3 * time.Hour * 24)

	seriesDownloads, err := app.DB.SeriesDownloadCounts()
	if err != nil {
		return fae.Wrap(err, "failed to get series download counts")
	}
	w.series_downloads = seriesDownloads

	err = w.db.Series.Query().ComplexOr(func(qq, qr *grimoire.QueryBuilder[*Series]) {
		qq.Where("active", true)
		qr.Where("kind", "donghua")
	}).Batch(100, func(results []*Series) error {
		for _, s := range results {
			w.addSeries(s)
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = w.db.Medium.Query().Where("_type", "Movie").Where("downloaded", false).Where("completed", false).Batch(100, func(results []*Medium) error {
		for _, m := range results {
			w.addMovie(m)
		}
		return nil
	})
	if err != nil {
		return err
	}

	q := w.db.Medium.Query().
		Where("_type", "Episode").
		In("series_id", w.series_ids).
		GreaterThan("release_date", start).LessThan("release_date", later).
		Where("completed", false).
		Where("skipped", false).
		In("missing", []interface{}{false, nil}).
		GreaterThan("season_number", 0).
		GreaterThan("episode_number", 0).
		Asc("release_date").Asc("season_number").Asc("episode_number").Asc("absolute_number")
	err = q.Batch(100, func(results []*Medium) error {
		for _, e := range results {
			w.addEpisode(e)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// w.log.Debugf("want: %+v", w)
	// w.log.Debugf("series ids: %d", len(w.series_ids))
	// for t, id := range w.series_titles {
	// 	for _, e := range w.series_episodes[id] {
	// 		w.log.Debugf("series: %s %dx%d", t, e.SeasonNumber, e.EpisodeNumber)
	// 	}
	// }
	// for t := range w.movies {
	// 	w.log.Debugf("movie: %s", t)
	// }

	return nil
}
