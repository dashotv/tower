package app

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d *Download) Created(ctx context.Context) error {
	db.log.Debugf("Download Created: %s", d.ID.Hex())
	return events.Send("tower.downloads", &EventTowerDownload{"created", d.ID.Hex(), d})
}
func (d *Download) Updated(ctx context.Context, result *mongo.UpdateResult) error {
	db.log.Debugf("Download Updated: %s", d.ID.Hex())
	return events.Send("tower.downloads", &EventTowerDownload{"updated", d.ID.Hex(), d})
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
	return events.Send("tower.episodes", &EventTowerEpisode{"created", e.ID.Hex(), e})
}
func (e *Episode) Updated(ctx context.Context, result *mongo.UpdateResult) error {
	return events.Send("tower.episodes", &EventTowerEpisode{"updated", e.ID.Hex(), e})
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
	if err := events.Send("tower.index.movies", m); err != nil {
		return errors.Wrap(err, "sending elastic search index")
	}
	return events.Send("tower.movies", &EventTowerMovie{"created", m.ID.Hex(), m})
}
func (m *Movie) Updated(ctx context.Context, result *mongo.UpdateResult) error {
	if err := events.Send("tower.index.movies", m); err != nil {
		return errors.Wrap(err, "sending elastic search index")
	}
	return events.Send("tower.movies", &EventTowerMovie{"updated", m.ID.Hex(), m})
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
		m.Directory = path(m.Title)
	} else {
		m.Directory = path(m.Directory)
	}

	return nil
}
func (s *Series) Created(ctx context.Context) error {
	if err := events.Send("tower.index.series", s); err != nil {
		return errors.Wrap(err, "sending elastic search index")
	}
	return events.Send("tower.series", &EventTowerSeries{"created", s.ID.Hex(), s})
}
func (s *Series) Updated(ctx context.Context, result *mongo.UpdateResult) error {
	if err := events.Send("tower.index.series", s); err != nil {
		return errors.Wrap(err, "sending elastic search index")
	}
	return events.Send("tower.series", &EventTowerSeries{"updated", s.ID.Hex(), s})
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
		s.Search = path(s.Search)
	}
	if s.Directory == "" {
		s.Directory = path(s.Title)
	} else {
		s.Directory = path(s.Directory)
	}

	return nil
}

func (t *Release) Updated(ctx context.Context, result *mongo.UpdateResult) error {
	return events.Send("tower.index.releases", t)
}
