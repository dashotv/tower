package app

import (
	"context"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d *Download) Created(ctx context.Context) error {
	return app.Events.Send("tower.downloads", &EventDownloads{"created", d.ID.Hex(), d})
}
func (d *Download) Updated(ctx context.Context, result *mongo.UpdateResult) error {
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

func (e *Episode) Created(ctx context.Context) error {
	return app.Events.Send("tower.episodes", &EventEpisodes{"created", e.ID.Hex(), e})
}
func (e *Episode) Updated(ctx context.Context, result *mongo.UpdateResult) error {
	return app.Events.Send("tower.episodes", &EventEpisodes{"updated", e.ID.Hex(), e})
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
		if p.Id.IsZero() {
			p.Id = primitive.NewObjectID()
		}
		if p.UpdatedAt.IsZero() {
			p.UpdatedAt = time.Now()
		}
	}

	return nil
}

func (m *Movie) Created(ctx context.Context) error {
	go func() {
		app.DB.processMovies([]*Movie{m})
		if err := app.Events.Send("tower.movies", &EventMovies{"created", m.ID.Hex(), m}); err != nil {
			app.DB.Log.Errorf("error updating movie: %s", err)
		}
	}()
	return nil
}
func (m *Movie) Updated(ctx context.Context, result *mongo.UpdateResult) error {
	go func() {
		app.DB.processMovies([]*Movie{m})
		if err := app.Events.Send("tower.movies", &EventMovies{"updated", m.ID.Hex(), m}); err != nil {
			app.DB.Log.Errorf("error updating movie: %s", err)
		}
	}()
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
		if p.Id.IsZero() {
			p.Id = primitive.NewObjectID()
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
func (s *Series) Created(ctx context.Context) error {
	go func() {
		app.DB.processSeries(s)
		if err := app.Events.Send("tower.series", &EventSeries{"created", s.ID.Hex(), s}); err != nil {
			app.DB.Log.Errorf("error updating series: %s", err)
		}
	}()
	return nil
}
func (s *Series) Updated(ctx context.Context, result *mongo.UpdateResult) error {
	go func() {
		app.DB.processSeries(s)
		if err := app.Events.Send("tower.series", &EventSeries{"updated", s.ID.Hex(), s}); err != nil {
			app.DB.Log.Errorf("error updating series: %s", err)
		}
	}()
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
		if p.Id.IsZero() {
			p.Id = primitive.NewObjectID()
		}
		if p.UpdatedAt.IsZero() {
			p.UpdatedAt = time.Now()
		}
	}

	if s.SearchParams == nil {
		s.SearchParams = &SearchParams{Type: "tv", Resolution: 1080, Verified: true}
	}

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

func (t *Release) Updated(ctx context.Context, result *mongo.UpdateResult) error {
	return app.Events.Send("tower.index.releases", t)
}
