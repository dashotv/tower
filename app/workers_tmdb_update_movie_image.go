package app

import (
	"context"
	"fmt"
	"path/filepath"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/dashotv/fae"
	"github.com/dashotv/minion"
)

// TmdbUpdateMovieImage
type TmdbUpdateMovieImage struct {
	minion.WorkerDefaults[*TmdbUpdateMovieImage]
	ID    string
	Type  string
	Path  string
	Ratio float32
}

func (j *TmdbUpdateMovieImage) Kind() string { return "TmdbUpdateMovieImage" }
func (j *TmdbUpdateMovieImage) Work(ctx context.Context, job *minion.Job[*TmdbUpdateMovieImage]) error {
	input := job.Args
	remote := app.Config.TmdbImages + input.Path
	extension := filepath.Ext(input.Path)[1:]
	local := fmt.Sprintf("movie-%s/%s", input.ID, input.Type)
	dest := fmt.Sprintf("%s/%s.%s", app.Config.DirectoriesImages, local, extension)
	thumb := fmt.Sprintf("%s/%s_thumb.%s", app.Config.DirectoriesImages, local, extension)

	if err := imageDownload(remote, dest); err != nil {
		return fae.Wrap(err, "downloading image")
	}

	height := 400
	width := int(float32(height) * input.Ratio)
	if err := imageResize(dest, thumb, width, height); err != nil {
		return fae.Wrap(err, "resizing image")
	}

	movie := &Movie{}
	if err := app.DB.Movie.Find(input.ID, movie); err != nil {
		return fae.Wrap(err, "finding movie")
	}
	app.DB.processMovies([]*Movie{movie})

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

	if err := app.DB.Movie.Update(movie); err != nil {
		return fae.Wrap(err, "updating movie")
	}

	return nil
}
