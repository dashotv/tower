package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/dashotv/tmdb"
)

var tmdbClient *tmdb.Client
var posterRatio float32 = 0.6666666666666666
var backgroundRatio float32 = 1.7777777777777777

func setupTmdb() error {
	tmdbClient = tmdb.New(os.Getenv("TMDB_API_TOKEN"))
	return nil
}

func TmdbUpdateMovie(payload any) error {
	id := payload.(string)

	movie := &Movie{}
	err := db.Movie.Find(id, movie)
	if err != nil {
		return errors.Wrap(err, "finding movie")
	}

	mid, err := strconv.Atoi(movie.SourceId)
	if err != nil {
		return errors.Wrap(err, "converting source id")
	}

	resp, err := tmdbClient.MovieDetails(mid, nil, nil)
	if err != nil {
		return errors.Wrap(err, "getting movie details")
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
		return errors.Wrap(err, "parsing release date")
	}
	movie.ReleaseDate = d
	if resp.PosterPath != nil {
		workers.EnqueueWithPayload("TmdbUpdateMovieImage", &ImagePayload{movie.ID.Hex(), "cover", tmdb.StringValue(resp.PosterPath), posterRatio})
	}
	if resp.BackdropPath != nil {
		workers.EnqueueWithPayload("TmdbUpdateMovieImage", &ImagePayload{movie.ID.Hex(), "background", tmdb.StringValue(resp.BackdropPath), backgroundRatio})
	}

	err = db.Movie.Update(movie)
	if err != nil {
		return errors.Wrap(err, "saving movie")
	}

	return nil
}

type ImagePayload struct {
	ID    string
	Type  string
	Path  string
	Ratio float32
}

func TmdbUpdateMovieImage(payload any) error {
	input := payload.(*ImagePayload)
	remote := cfg.Tmdb.Images + input.Path
	extension := filepath.Ext(input.Path)[1:]
	local := fmt.Sprintf("movie-%s/%s", input.ID, input.Type)
	dest := fmt.Sprintf("%s/%s.%s", cfg.Directories.Images, local, extension)
	thumb := fmt.Sprintf("%s/%s_thumb.%s", cfg.Directories.Images, local, extension)

	if err := imageDownload(remote, dest); err != nil {
		return errors.Wrap(err, "downloading image")
	}

	height := 400
	width := int(float32(height) * input.Ratio)
	if err := imageResize(dest, thumb, width, height); err != nil {
		return errors.Wrap(err, "resizing image")
	}

	movie := &Movie{}
	if err := db.Movie.Find(input.ID, movie); err != nil {
		return errors.Wrap(err, "finding movie")
	}

	var img *Path
	for _, p := range movie.Paths {
		if string(p.Type) == input.Type {
			img = p
			break
		}
	}

	if img == nil {
		img = &Path{}
	}

	img.Type = primitive.Symbol(input.Type)
	img.Remote = remote
	img.Local = local
	img.Extension = extension

	movie.Paths = append(movie.Paths, img)

	if err := db.Movie.Update(movie); err != nil {
		return errors.Wrap(err, "updating movie")
	}

	return nil
}
