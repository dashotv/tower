package app

import (
	"context"

	"github.com/dashotv/minion"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

// CreateMediaFromRequests creates media from approved requests
type CreateMediaFromRequests struct{}

func (j *CreateMediaFromRequests) Kind() string { return "CreateMediaFromRequests" }
func (j *CreateMediaFromRequests) Work(ctx context.Context, job *minion.Job[*CreateMediaFromRequests]) error {
	log := log.Named("job.CreateMediaFromRequests")

	requests, err := db.Request.Query().Where("status", "approved").Run()
	if err != nil {
		return errors.Wrap(err, "querying requests")
	}

	for _, r := range requests {
		log.Infof("processing request: %s", r.Title)
		if r.Source == "tmdb" {
			err := createMovieFromRequest(r)
			if err != nil {
				log.Errorf("creating movie from request: %s", err)
				r.Status = "failed"
			} else {
				log.Infof("created movie: %s", r.Title)
				r.Status = "completed"
			}
		} else if r.Source == "tvdb" {
			err := createShowFromRequest(r)
			if err != nil {
				log.Errorf("creating series from request: %s", err)
				r.Status = "failed"
			} else {
				log.Infof("created series: %s", r.Title)
				r.Status = "completed"
			}
		}

		log.Infof("request: [%s] %s", r.Status, r.Title)
		if err := db.Request.Update(r); err != nil {
			return errors.Wrap(err, "updating request")
		}

		if err := events.Send("tower.requests", &EventTowerRequest{Event: "update", ID: r.ID.Hex(), Request: r}); err != nil {
			return errors.Wrap(err, "sending event")
		}
	}
	return nil
}

func createShowFromRequest(r *Request) error {
	count, err := db.Series.Count(bson.M{"_type": "Series", "source": r.Source, "source_id": r.SourceId})
	if err != nil {
		return errors.Wrap(err, "counting series")
	}
	if count > 0 {
		return nil
	}

	s := &Series{
		Type:     "Series",
		Source:   r.Source,
		SourceId: r.SourceId,
		Title:    r.Title,
		Kind:     "tv",
	}

	err = db.Series.Save(s)
	if err != nil {
		return errors.Wrap(err, "saving show")
	}

	if err := workers.Enqueue(&TvdbUpdateSeries{s.ID.Hex()}); err != nil {
		return errors.Wrap(err, "queueing update job")
	}
	return nil
}

func createMovieFromRequest(r *Request) error {
	count, err := db.Series.Count(bson.M{"_type": "Movie", "source": r.Source, "source_id": r.SourceId})
	if err != nil {
		return errors.Wrap(err, "counting series")
	}
	if count > 0 {
		return nil
	}

	m := &Movie{
		Type:     "Movie",
		Source:   r.Source,
		SourceId: r.SourceId,
		Title:    r.Title,
		Kind:     "movies",
	}

	err = db.Movie.Save(m)
	if err != nil {
		return errors.Wrap(err, "saving movie")
	}

	if err := workers.Enqueue(&TmdbUpdateMovie{m.ID.Hex()}); err != nil {
		return errors.Wrap(err, "queueing update job")
	}
	return nil
}
