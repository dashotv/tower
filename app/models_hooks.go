package app

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (e *Episode) Saving() error {
	// Call the DefaultModel Saving hook
	if err := e.DefaultModel.Saving(); err != nil {
		return err
	}

	for _, p := range e.Paths {
		if p.Id.IsZero() {
			p.Id = primitive.NewObjectID()
		}
		p.UpdatedAt = time.Now()
	}

	return events.Send("tower.episodes", &EventTowerEpisode{"updated", e.ID.Hex(), e})
}

func (m *Movie) Saving() error {
	// Call the DefaultModel Saving hook
	if err := m.DefaultModel.Saving(); err != nil {
		return err
	}

	for _, p := range m.Paths {
		if p.Id.IsZero() {
			p.Id = primitive.NewObjectID()
		}
		p.UpdatedAt = time.Now()
	}

	if m.SearchParams == nil {
		m.SearchParams = &SearchParams{Type: "movie", Resolution: 1080, Verified: true}
	}

	return events.Send("tower.movies", &EventTowerMovie{"updated", m.ID.Hex(), m})
}

func (s *Series) Saving() error {
	// Call the DefaultModel Saving hook
	if err := s.DefaultModel.Saving(); err != nil {
		return err
	}

	for _, p := range s.Paths {
		if p.Id.IsZero() {
			p.Id = primitive.NewObjectID()
		}
		p.UpdatedAt = time.Now()
	}

	if s.SearchParams == nil {
		s.SearchParams = &SearchParams{Type: "tv", Resolution: 1080, Verified: true}
	}

	return events.Send("tower.series", &EventTowerSeries{"updated", s.ID.Hex(), s})
}
