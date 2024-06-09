package importer

import (
	"github.com/dashotv/fae"
	"github.com/dashotv/tmdb"
)

func (i *Importer) loadMovie(id int) (*Movie, error) {
	movie, err := i.loadMovieTmdb(id)
	if err != nil {
		return nil, fae.Wrap(err, "base")
	}

	return movie, nil
}

func (i *Importer) loadMovieTmdb(id int) (*Movie, error) {
	movie, err := i.Tmdb.MovieDetails(id, nil, nil)
	if err != nil {
		return nil, err
	}

	return &Movie{
		ID:          tmdb.Int64Value(movie.ID),
		ImdbID:      tmdb.StringValue(movie.ImdbID),
		Title:       tmdb.StringValue(movie.Title),
		Description: tmdb.StringValue(movie.Overview),
		Airdate:     tmdb.StringValue(movie.ReleaseDate),
		Poster:      tmdb.StringValue(movie.PosterPath),
		Backdrop:    tmdb.StringValue(movie.BackdropPath),
	}, nil
}

func (i *Importer) loadMovieImages(id int) ([]string, []string, error) {
	images, err := i.Tmdb.MovieImages(id, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	backdrops := []string{}
	posters := []string{}
	for _, image := range images.Backdrops {
		backdrops = append(backdrops, i.Opts.TmdbImageURL+tmdb.StringValue(image.FilePath))
	}
	for _, image := range images.Posters {
		posters = append(posters, i.Opts.TmdbImageURL+tmdb.StringValue(image.FilePath))
	}

	return posters, backdrops, nil
}
