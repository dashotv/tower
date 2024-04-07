package app

import (
	"context"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"github.com/samber/lo"

	"github.com/dashotv/fae"
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
		movies:           map[string]string{},
		series_unwatched: map[string]int{},
		series_titles:    map[string]string{},
		series_episodes:  map[string][]*Episode{},
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
	series_episodes  map[string][]*Episode
	movies           map[string]string
}

func (w *Want) Release(r *runic.Release) string {
	if r == nil {
		return ""
	}

	if !lo.Contains(w.preferred, r.Group) && !lo.Contains(w.groups, r.Group) {
		return ""
	}

	if r.Title == "" {
		return ""
	}

	switch r.Type {
	case "movies":
		return w.Movie(r.Title)
	case "tv", "anime":
		return w.Episode(r.Title, r.Season, r.Episode)
	default:
		return ""
	}
}

func (w *Want) Movie(title string) string {
	title = path(title)
	// w.log.Debugf("movie %s", title)
	f, ok := w.movies[title]
	if !ok {
		return ""
	}
	return f
}

func (w *Want) Episode(seriesTitle string, season int, episode int) string {
	seriesTitle = path(seriesTitle)
	w.log.Debugf("series %s", seriesTitle)
	series, ok := w.series_titles[seriesTitle]
	if !ok {
		return ""
	}
	// w.log.Debugf("series %s %dx%d", seriesTitle, season, episode)
	for _, e := range w.series_episodes[series] {
		if e.SeasonNumber == season && e.EpisodeNumber == episode {
			return e.ID.Hex()
		}
		if e.AbsoluteNumber == episode {
			return e.ID.Hex()
		}
	}
	return ""
}

func (w *Want) NextEpisode(seriesID string) *Episode {
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

func (w *Want) addMovie(m *Movie) {
	if m.Directory != "" {
		w.movies[m.Directory] = m.ID.Hex()
	}
	if m.Search != "" && m.Search != m.Directory {
		w.movies[m.Search] = m.ID.Hex()
	}
}

func (w *Want) addEpisode(e *Episode) {
	// if len(w.series_episodes[e.SeriesID.Hex()]) < seriesWantedBuffer {
	if len(w.series_episodes[e.SeriesID.Hex()]) >= (seriesWantedBuffer - w.series_unwatched[e.SeriesID.Hex()]) {
		return
	}
	w.series_episodes[e.SeriesID.Hex()] = append(w.series_episodes[e.SeriesID.Hex()], e)
}

func (w *Want) Build() error {
	utc := time.Now().UTC()
	after := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	today := time.Date(utc.Year(), utc.Month(), utc.Day(), 0, 0, 0, 0, time.UTC)
	tomorrow := today.Add(time.Hour * 24)

	seriesDownloads, err := app.DB.SeriesDownloadCounts()
	if err != nil {
		return fae.Wrap(err, "failed to get series download counts")
	}
	w.series_downloads = seriesDownloads

	err = w.db.Series.Query().Where("active", true).Batch(100, func(results []*Series) error {
		for _, s := range results {
			w.addSeries(s)
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = w.db.Movie.Query().Where("downloaded", false).Where("completed", false).Batch(100, func(results []*Movie) error {
		for _, m := range results {
			w.addMovie(m)
		}
		return nil
	})
	if err != nil {
		return err
	}

	q := w.db.Episode.Query().
		In("series_id", w.series_ids).
		GreaterThan("release_date", after).LessThan("release_date", tomorrow).
		Where("completed", false).
		Where("skipped", false).
		In("missing", []interface{}{false, nil}).
		GreaterThan("season_number", 0).
		GreaterThan("episode_number", 0).
		Asc("release_date").Asc("season_number").Asc("episode_number").Asc("absolute_number")
	err = q.Batch(100, func(results []*Episode) error {
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
