package app

import (
	"context"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/sync/errgroup"

	"github.com/dashotv/fae"
	"github.com/dashotv/minion"
)

type MovieDelete struct {
	minion.WorkerDefaults[*MovieDelete]
	ID string `bson:"id" json:"id"`
}

func (j *MovieDelete) Kind() string { return "movie_delete" }
func (j *MovieDelete) Work(ctx context.Context, job *minion.Job[*MovieDelete]) error {
	id := job.Args.ID

	movie := &Movie{}
	if err := app.DB.Movie.Find(id, movie); err != nil {
		return fae.Wrap(err, "finding movie")
	}

	// delete files
	if err := mediumIdDeletePaths(id); err != nil {
		return fae.Wrap(err, "deleting paths")
	}

	// remove downloads referencing movie
	_, err := app.DB.Download.Collection.DeleteMany(ctx, bson.M{"medium_id": movie.ID})
	if err != nil {
		return fae.Wrap(err, "deleting downloads")
	}

	// remove movie
	if err := app.DB.Movie.Delete(movie); err != nil {
		return fae.Wrap(err, "delete medium")
	}

	return nil
}

type MovieUpdateAll struct {
	minion.WorkerDefaults[*MovieUpdateAll]
}

func (j *MovieUpdateAll) Kind() string { return "movie_update_al" }
func (j *MovieUpdateAll) Work(ctx context.Context, job *minion.Job[*MovieUpdateAll]) error {
	a := ContextApp(ctx)

	movies, err := a.DB.Movie.Query().Limit(-1).Run()
	if err != nil {
		return fae.Wrap(err, "querying movies")
	}

	for _, m := range movies {
		a.Log.Infof("updating movie: %s", m.Title)
		a.Workers.Enqueue(&MovieUpdate{ID: m.ID.Hex(), Title: m.Display, SkipImages: true})
	}

	return nil
}

type MovieUpdate struct {
	minion.WorkerDefaults[*MovieUpdate]
	ID         string `bson:"id" json:"id"`
	Title      string `bson:"title" json:"title"`
	SkipImages bool   `bson:"skip_images" json:"skip_images"`
}

func (j *MovieUpdate) Kind() string { return "movie_update" }
func (j *MovieUpdate) Work(ctx context.Context, job *minion.Job[*MovieUpdate]) error {
	eg, ctx := errgroup.WithContext(ctx)
	a := ContextApp(ctx)
	id := job.Args.ID

	movie := &Medium{}
	err := app.DB.Medium.Find(id, movie)
	if err != nil {
		return err
	}

	if movie.SourceID == "" || movie.Source != "tmdb" {
		return fae.New("movie source not tmdb or missing id")
	}

	tmdbid, err := strconv.Atoi(movie.SourceID)
	if err != nil {
		return fae.Wrap(err, "converting source id")
	}

	eg.Go(func() error {
		m, err := a.Importer.Movie(tmdbid)
		if err != nil {
			return fae.Wrap(err, "loading movie")
		}

		movie.Title = m.Title
		movie.Description = m.Description
		movie.ReleaseDate = dateFromString(m.Airdate)
		if movie.Display == "" {
			movie.Display = m.Title
		}
		if movie.Search == "" {
			movie.Search = path(m.Title)
		}
		if movie.Directory == "" {
			movie.Directory = path(m.Title)
		}

		if !job.Args.SkipImages {
			if m.Poster != "" {
				eg.Go(func() error {
					err := mediumImage(movie, "cover", a.Config.TmdbImages+m.Poster, posterRatio)
					if err != nil {
						app.Log.Errorf("movie %s cover: %v", movie.ID.Hex(), err)
					}
					return nil
				})
			}
			if m.Backdrop != "" {
				eg.Go(func() error {
					err := mediumImage(movie, "background", a.Config.TmdbImages+m.Backdrop, backgroundRatio)
					if err != nil {
						app.Log.Errorf("movie %s background: %v", movie.ID.Hex(), err)
					}
					return nil
				})
			}
		}

		return nil
	})

	if err := eg.Wait(); err != nil {
		return fae.Wrapf(err, "movie: %s", movie.Title)
	}

	err = app.DB.Medium.Save(movie)
	if err != nil {
		return fae.Wrap(err, "saving movie")
	}

	return nil
}
