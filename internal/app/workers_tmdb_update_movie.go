package app

import (
	"context"
	"strconv"
	"time"

	"github.com/dashotv/fae"
	"github.com/dashotv/minion"
	"github.com/dashotv/tmdb"
)

// TmdbUpdateMovie
type TmdbUpdateMovie struct {
	minion.WorkerDefaults[*TmdbUpdateMovie]
	ID        string
	JustMedia bool
}

func (j *TmdbUpdateMovie) Kind() string { return "TmdbUpdateMovie" }
func (j *TmdbUpdateMovie) Work(ctx context.Context, job *minion.Job[*TmdbUpdateMovie]) error {
	id := job.Args.ID

	movie := &Movie{}
	err := app.DB.Movie.Find(id, movie)
	if err != nil {
		return fae.Wrap(err, "finding movie")
	}
	app.DB.processMovies([]*Movie{movie})

	mid, err := strconv.Atoi(movie.SourceId)
	if err != nil {
		return fae.Wrap(err, "converting source id")
	}

	resp, err := app.Tmdb.MovieDetails(mid, nil, nil)
	if err != nil {
		return fae.Wrap(err, "getting movie details")
	}

	movie.Title = tmdb.StringValue(resp.Title)
	if movie.Display == "" {
		movie.Display = movie.Title
	}
	if movie.Search == "" {
		movie.Search = path(movie.Title)
	}
	if movie.Directory == "" {
		movie.Directory = path(movie.Title)
	}
	movie.ImdbId = tmdb.StringValue(resp.ImdbID)
	movie.Description = tmdb.StringValue(resp.Overview)
	d, err := time.Parse("2006-01-02", tmdb.StringValue(resp.ReleaseDate))
	if err != nil {
		return fae.Wrap(err, "parsing release date")
	}
	movie.ReleaseDate = d

	if err := app.Workers.Enqueue(&PathCleanup{ID: id}); err != nil {
		return fae.Wrap(err, "enqueuing media paths")
	}

	if !job.Args.JustMedia {
		if resp.PosterPath != nil {
			app.Workers.Enqueue(&TmdbUpdateMovieImage{ID: movie.ID.Hex(), Type: "cover", Path: tmdb.StringValue(resp.PosterPath), Ratio: posterRatio})
		}
		if resp.BackdropPath != nil {
			app.Workers.Enqueue(&TmdbUpdateMovieImage{ID: movie.ID.Hex(), Type: "background", Path: tmdb.StringValue(resp.BackdropPath), Ratio: backgroundRatio})
		}
	}

	err = app.DB.Movie.Update(movie)
	if err != nil {
		return fae.Wrap(err, "saving movie")
	}

	return nil
}
