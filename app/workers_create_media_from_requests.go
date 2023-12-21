package app

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/dashotv/minion"
)

// CreateMediaFromRequests creates media from approved requests
type CreateMediaFromRequests struct {
	minion.WorkerDefaults[*CreateMediaFromRequests]
}

func (j *CreateMediaFromRequests) Kind() string { return "CreateMediaFromRequests" }
func (j *CreateMediaFromRequests) Work(ctx context.Context, job *minion.Job[*CreateMediaFromRequests]) error {
	requests, err := app.DB.Request.Query().Where("status", "approved").Run()
	if err != nil {
		return errors.Wrap(err, "querying requests")
	}

	for _, r := range requests {
		app.Log.Infof("processing request: %s", r.Title)
		if r.Source == "tmdb" {
			err := createMovieFromRequest(r)
			if err != nil {
				app.Log.Errorf("creating movie from request: %s", err)
				r.Status = "failed"
			} else {
				app.Log.Infof("created movie: %s", r.Title)
				r.Status = "completed"
			}
		} else if r.Source == "tvdb" {
			err := createShowFromRequest(r)
			if err != nil {
				app.Log.Errorf("creating series from request: %s", err)
				r.Status = "failed"
			} else {
				app.Log.Infof("created series: %s", r.Title)
				r.Status = "completed"
			}
		}

		app.Log.Infof("request: [%s] %s", r.Status, r.Title)
		if err := app.DB.Request.Update(r); err != nil {
			return errors.Wrap(err, "updating request")
		}

		if err := app.Events.Send("tower.requests", &EventRequests{Event: "update", Id: r.ID.Hex(), Request: r}); err != nil {
			return errors.Wrap(err, "sending event")
		}
	}
	return nil
}

func createShowFromRequest(r *Request) error {
	count, err := app.DB.Series.Count(bson.M{"_type": "Series", "source": r.Source, "source_id": r.SourceId})
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

	err = app.DB.Series.Save(s)
	if err != nil {
		return errors.Wrap(err, "saving show")
	}

	if err := app.Workers.Enqueue(&TvdbUpdateSeries{ID: s.ID.Hex()}); err != nil {
		return errors.Wrap(err, "queueing update job")
	}
	return nil
}

func createMovieFromRequest(r *Request) error {
	count, err := app.DB.Series.Count(bson.M{"_type": "Movie", "source": r.Source, "source_id": r.SourceId})
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

	err = app.DB.Movie.Save(m)
	if err != nil {
		return errors.Wrap(err, "saving movie")
	}

	if err := app.Workers.Enqueue(&TmdbUpdateMovie{ID: m.ID.Hex()}); err != nil {
		return errors.Wrap(err, "queueing update job")
	}
	return nil
}
