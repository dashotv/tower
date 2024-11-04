package app

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/dashotv/fae"
	"github.com/dashotv/minion"
)

// CreateMediaFromRequests creates media from approved requests
type CreateMediaFromRequests struct {
	minion.WorkerDefaults[*CreateMediaFromRequests]
}

func (j *CreateMediaFromRequests) Kind() string { return "CreateMediaFromRequests" }
func (j *CreateMediaFromRequests) Work(ctx context.Context, job *minion.Job[*CreateMediaFromRequests]) error {
	a := ContextApp(ctx)
	if a == nil {
		return fae.New("CreateMediaFromRequests: no app in context")
	}

	requests, err := a.DB.Request.Query().Where("status", "approved").Run()
	if err != nil {
		return fae.Wrap(err, "querying requests")
	}

	for _, r := range requests {
		a.Log.Infof("processing request: %s", r.Title)
		if r.Source == "tmdb" {
			err := a.createMovieFromRequest(r)
			if err != nil {
				a.Log.Errorf("creating movie from request: %s", err)
				r.Status = "failed"
			} else {
				a.Log.Infof("created movie: %s", r.Title)
				r.Status = "completed"
			}
		} else if r.Source == "imdb" {
			err := a.createMovieFromImdbRequest(r)
			if err != nil {
				a.Log.Errorf("creating movie from request: %s", err)
				r.Status = "failed"
			} else {
				a.Log.Infof("created series: %s", r.Title)
				r.Status = "completed"
			}
		} else if r.Source == "tvdb" {
			err := a.createShowFromRequest(r)
			if err != nil {
				a.Log.Errorf("creating series from request: %s", err)
				r.Status = "failed"
			} else {
				a.Log.Infof("created series: %s", r.Title)
				r.Status = "completed"
			}
		}

		a.Log.Infof("request: [%s] %s", r.Status, r.Title)
		if err := a.DB.Request.Update(r); err != nil {
			return fae.Wrap(err, "updating request")
		}

		if err := a.Events.Send("tower.requests", &EventRequests{Event: "update", ID: r.ID.Hex(), Request: r}); err != nil {
			return fae.Wrap(err, "sending event")
		}
	}
	return nil
}

func (a *Application) createShowFromRequest(r *Request) error {
	count, err := a.DB.Series.Count(bson.M{"_type": "Series", "source": r.Source, "source_id": r.SourceID})
	if err != nil {
		return fae.Wrap(err, "counting series")
	}
	if count > 0 {
		return nil
	}

	s := &Series{
		Type:     "Series",
		Source:   r.Source,
		SourceID: r.SourceID,
		Title:    r.Title,
		Kind:     "tv",
	}

	err = a.DB.Series.Save(s)
	if err != nil {
		return fae.Wrap(err, "saving show")
	}

	if err := a.Workers.Enqueue(&SeriesUpdate{ID: s.ID.Hex()}); err != nil {
		return fae.Wrap(err, "queueing update job")
	}
	return nil
}

func (a *Application) createMovieFromRequest(r *Request) error {
	count, err := a.DB.Movie.Count(bson.M{"source": r.Source, "source_id": r.SourceID})
	if err != nil {
		return fae.Wrap(err, "counting movies")
	}
	if count > 0 {
		return nil
	}

	m := &Movie{
		Type:     "Movie",
		Source:   r.Source,
		SourceID: r.SourceID,
		Title:    r.Title,
		Kind:     "movies",
	}

	err = a.DB.Movie.Save(m)
	if err != nil {
		return fae.Wrap(err, "saving movie")
	}

	if err := a.Workers.Enqueue(&MovieUpdate{ID: m.ID.Hex()}); err != nil {
		return fae.Wrap(err, "queueing update job")
	}
	return nil
}

func (a *Application) createMovieFromImdbRequest(r *Request) error {
	id, err := a.Importer.ImdbToTmdb(r.SourceID)
	if err != nil {
		return fae.Wrap(err, "converting imdb to tmdb")
	}

	if id == 0 {
		return fae.Errorf("no movie found for imdb %s", r.SourceID)
	}

	r.Source = "tmdb"
	r.SourceID = fmt.Sprintf("%d", id)

	return a.createMovieFromRequest(r)
}
