package app

import (
	"context"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d *Download) Created(ctx context.Context) error {
	if d.Title == "" {
		app.DB.processDownloads([]*Download{d})
	}
	return app.Events.Send("tower.downloads", &EventDownloads{"created", d.ID.Hex(), d})
}
func (d *Download) Updated(ctx context.Context, result *mongo.UpdateResult) error {
	if d.Title == "" {
		app.DB.processDownloads([]*Download{d})
	}
	return app.Events.Send("tower.downloads", &EventDownloads{"updated", d.ID.Hex(), d})
}
func (d *Download) Saving() error {
	// Call the DefaultModel Saving hook
	if err := d.DefaultModel.Saving(); err != nil {
		return err
	}
	if d.Files == nil {
		d.Files = []*DownloadFile{}
	}

	return nil
}

func sendEpisodeEvent(event string, e *Episode) {
	go func() {
		app.DB.processEpisode(e)
		if err := app.Events.Send("tower.episodes", &EventEpisodes{event, e.ID.Hex(), e}); err != nil {
			app.DB.Log.Errorf("error updating movie: %s", err)
		}
	}()
}
func (e *Episode) Created(ctx context.Context) error {
	sendEpisodeEvent("created", e)
	return nil
}
func (e *Episode) Updated(ctx context.Context, result *mongo.UpdateResult) error {
	sendEpisodeEvent("updated", e)
	return nil
}
func (e *Episode) Deleted(ctx context.Context, result *mongo.DeleteResult) error {
	sendEpisodeEvent("deleted", e)
	return nil
}
func (e *Episode) Saving() error {
	// Call the DefaultModel Saving hook
	if err := e.DefaultModel.Saving(); err != nil {
		return err
	}

	if e.Paths == nil {
		e.Paths = []*Path{}
	}
	for _, p := range e.Paths {
		if p.ID.IsZero() {
			p.ID = primitive.NewObjectID()
		}
		if p.UpdatedAt.IsZero() {
			p.UpdatedAt = time.Now()
		}
	}

	return nil
}

func sendMovieEvent(event string, m *Movie) {
	go func() {
		app.DB.processMovies([]*Movie{m})
		if err := app.Events.Send("tower.movies", &EventMovies{event, m.ID.Hex(), m}); err != nil {
			app.DB.Log.Errorf("error updating movie: %s", err)
		}
	}()
}
func (m *Movie) Created(ctx context.Context) error {
	sendMovieEvent("created", m)
	return nil
}
func (m *Movie) Updated(ctx context.Context, result *mongo.UpdateResult) error {
	sendMovieEvent("updated", m)
	return nil
}
func (m *Movie) Deleted(ctx context.Context, result *mongo.DeleteResult) error {
	sendMovieEvent("deleted", m)
	return nil
}
func (m *Movie) Saving() error {
	// Call the DefaultModel Saving hook
	if err := m.DefaultModel.Saving(); err != nil {
		return err
	}

	if m.Paths == nil {
		m.Paths = []*Path{}
	}
	for _, p := range m.Paths {
		if p.ID.IsZero() {
			p.ID = primitive.NewObjectID()
		}
		if p.UpdatedAt.IsZero() {
			p.UpdatedAt = time.Now()
		}
	}

	if m.SearchParams == nil {
		m.SearchParams = &SearchParams{Type: "movies", Resolution: 1080, Verified: true}
	}

	if m.Display == "" {
		m.Display = m.Title
	}
	if m.Search == "" {
		m.Search = path(m.Title)
	} else {
		m.Search = path(m.Search)
	}
	if m.Directory == "" {
		m.Directory = directory(m.Title)
	} else {
		m.Directory = directory(m.Directory)
	}

	return nil
}

func sendSeriesEvent(event string, s *Series) {
	go func() {
		app.DB.processSeries(s)
		if err := app.Events.Send("tower.series", &EventSeries{event, s.ID.Hex(), s}); err != nil {
			app.DB.Log.Errorf("error updating series: %s", err)
		}
	}()
}
func (s *Series) Created(ctx context.Context) error {
	sendSeriesEvent("created", s)
	return nil
}
func (s *Series) Updated(ctx context.Context, result *mongo.UpdateResult) error {
	sendSeriesEvent("updated", s)
	return nil
}
func (s *Series) Deleted(ctx context.Context, result *mongo.DeleteResult) error {
	sendSeriesEvent("updated", s)
	return nil
}
func (s *Series) Saving() error {
	// Call the DefaultModel Saving hook
	if err := s.DefaultModel.Saving(); err != nil {
		return err
	}

	if s.Paths == nil {
		s.Paths = []*Path{}
	}
	for _, p := range s.Paths {
		if p.ID.IsZero() {
			p.ID = primitive.NewObjectID()
		}
		if p.UpdatedAt.IsZero() {
			p.UpdatedAt = time.Now()
		}
	}

	if s.SearchParams == nil {
		s.SearchParams = &SearchParams{Type: "tv", Resolution: 1080, Verified: true}
	}
	if s.SearchParams != nil && s.SearchParams.Type == "" {
		s.SearchParams.Type = "tv"
	}
	// if s.SearchParams != nil && s.SearchParams.Resolution == 0 {
	// 	s.SearchParams.Resolution = 1080
	// }

	if s.Display == "" {
		s.Display = s.Title
	}
	if s.Search == "" {
		s.Search = path(s.Title)
	} else {
		e := strings.Split(s.Search, ":")
		p := path(e[0])
		if len(e) > 1 {
			p = p + ":" + e[1]
		}
		s.Search = p
	}
	if s.Directory == "" {
		s.Directory = directory(s.Title)
	} else {
		s.Directory = directory(s.Directory)
	}

	return nil
}

func sendMediumEvent(event string, m *Medium) {
	go func() {
		switch m.Type {
		case "Movie":
			n := &Movie{}
			if err := app.DB.Movie.FindByID(m.ID, n); err != nil {
				app.DB.Log.Errorf("error finding movie: %s", err)
				return
			}
			app.DB.processMovies([]*Movie{n})
			if err := app.Events.Send("tower.movies", &EventMovies{event, n.ID.Hex(), n}); err != nil {
				app.DB.Log.Errorf("error updating movie: %s", err)
			}
		case "Series":
			s := &Series{}
			if err := app.DB.Series.FindByID(m.ID, s); err != nil {
				app.DB.Log.Errorf("error finding series: %s", err)
				return
			}
			app.DB.processSeries(s)
			if err := app.Events.Send("tower.series", &EventSeries{event, s.ID.Hex(), s}); err != nil {
				app.DB.Log.Errorf("error updating series: %s", err)
			}
		}
	}()
}

func (m *Medium) Created(ctx context.Context) error {
	sendMediumEvent("created", m)
	return nil
}
func (m *Medium) Updated(ctx context.Context, result *mongo.UpdateResult) error {
	sendMediumEvent("updated", m)
	return nil
}
func (m *Medium) Deleted(ctx context.Context, result *mongo.DeleteResult) error {
	sendMediumEvent("deleted", m)
	return nil
}
func (m *Medium) Saving() error {
	// Call the DefaultModel Saving hook
	if err := m.DefaultModel.Saving(); err != nil {
		return err
	}

	if m.Paths == nil {
		m.Paths = []*Path{}
	}
	for _, p := range m.Paths {
		if p.ID.IsZero() {
			p.ID = primitive.NewObjectID()
		}
		if p.UpdatedAt.IsZero() {
			p.UpdatedAt = time.Now()
		}
	}

	return nil
}

func (t *Release) Updated(ctx context.Context, result *mongo.UpdateResult) error {
	return app.Events.Send("tower.index.releases", t)
}
