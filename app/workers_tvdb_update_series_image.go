package app

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/dashotv/minion"
	"github.com/dashotv/tmdb"
	"github.com/dashotv/tvdb"
)

func TvdbUpdateSeriesCover(id string, sid int64) error {
	if err := TvdbUpdateSeriesCoverTvdb(id, sid); err != nil {
		app.Log.Warnf("failed to update cover from tvdb: %s", err)
	} else {
		return nil
	}
	if err := TvdbUpdateSeriesCoverTmdb(id, sid); err != nil {
		app.Log.Warnf("failed to update cover from tmdb: %s", err)
	} else {
		return nil
	}
	if err := TvdbUpdateSeriesCoverFanart(id, sid); err != nil {
		app.Log.Warnf("failed to update cover from fanart: %s", err)
	} else {
		return nil
	}

	return nil
}

func TvdbUpdateSeriesCoverTmdb(id string, sid int64) error {
	find, err := app.Tmdb.FindByID(fmt.Sprintf("%d", sid), "tvdb_id", tmdb.String("en-US"))
	if err != nil {
		return errors.Wrap(err, "getting series artworks")
	}
	if find.TvResults == nil || len(find.TvResults) == 0 {
		return errors.New("can't find id")
	}

	res := find.TvResults[0].(map[string]interface{})
	found := int(res["id"].(float64))

	app.Log.Named("TvdbUpdateSeriesCoverTmdb").Info("found:", found)
	resp, err := app.Tmdb.TvSeriesImages(found, nil, nil)
	if err != nil {
		return errors.Wrap(err, "getting series artworks")
	}

	if resp.Posters == nil || len(resp.Posters) == 0 {
		return errors.New("no data")
	}

	url := app.Config.TmdbImages + *resp.Posters[0].FilePath
	if err := app.Workers.Enqueue(&TvdbUpdateSeriesImage{ID: id, Type: "cover", Path: url, Ratio: posterRatio}); err != nil {
		return errors.Wrap(err, "enqueuing series episodes")
	}

	return nil
}

func TvdbUpdateSeriesCoverFanart(id string, sid int64) error {
	ftv, err := app.Fanart.GetShowImages(fmt.Sprintf("%d", sid))
	if err != nil {
		return errors.Wrap(err, "getting fanart")
	}
	if len(ftv.Posters) == 0 {
		return errors.New("no posters")
	}

	if err := app.Workers.Enqueue(&TvdbUpdateSeriesImage{ID: id, Type: "cover", Path: ftv.Posters[0].URL, Ratio: posterRatio}); err != nil {
		return errors.Wrap(err, "enqueuing series episodes")
	}

	return nil
}

func TvdbUpdateSeriesCoverTvdb(id string, sid int64) error {
	app.Log.Named("TvdbUpdateSeriesCover").Info("updating series images: cover")
	r, err := app.Tvdb.GetSeriesArtworks(sid, tvdb.String("eng"), tvdb.Int64(int64(2)))
	if err != nil {
		return errors.Wrap(err, "getting series artworks")
	}

	if r.Data == nil || len(r.Data.Artworks) == 0 {
		return errors.New("no data")
	}

	cover := r.Data.Artworks[0]
	if err := app.Workers.Enqueue(&TvdbUpdateSeriesImage{ID: id, Type: "cover", Path: tvdb.StringValue(cover.Image), Ratio: posterRatio}); err != nil {
		return errors.Wrap(err, "enqueuing series episodes")
	}

	return nil
}

func TvdbUpdateSeriesBackground(id string, sid int64) error {
	if err := TvdbUpdateSeriesBackgroundTvdb(id, sid); err != nil {
		app.Log.Warnf("failed to update background from tvdb: %s", err)
	} else {
		return nil
	}
	if err := TvdbUpdateSeriesBackgroundTmdb(id, sid); err != nil {
		app.Log.Warnf("failed to update background from fanart: %s", err)
	} else {
		return nil
	}
	if err := TvdbUpdateSeriesBackgroundFanart(id, sid); err != nil {
		app.Log.Warnf("failed to update background from fanart: %s", err)
	} else {
		return nil
	}

	return nil
}

func TvdbUpdateSeriesBackgroundTmdb(id string, sid int64) error {
	find, err := app.Tmdb.FindByID(fmt.Sprintf("%d", sid), "tvdb_id", tmdb.String("en-US"))
	if err != nil {
		return errors.Wrap(err, "getting series artworks")
	}
	if find.TvResults == nil || len(find.TvResults) == 0 {
		return errors.New("can't find id")
	}

	res := find.TvResults[0].(map[string]interface{})
	found := int(res["id"].(float64))

	app.Log.Named("TvdbUpdateSeriesBackground").Info("found:", found)
	resp, err := app.Tmdb.TvSeriesImages(found, nil, nil)
	if err != nil {
		return errors.Wrap(err, "getting series artworks")
	}

	if resp.Backdrops == nil || len(resp.Backdrops) == 0 {
		return errors.New("no data")
	}

	url := app.Config.TmdbImages + *resp.Backdrops[0].FilePath
	if err := app.Workers.Enqueue(&TvdbUpdateSeriesImage{ID: id, Type: "background", Path: url, Ratio: backgroundRatio}); err != nil {
		return errors.Wrap(err, "enqueuing series episodes")
	}

	return nil
}

func TvdbUpdateSeriesBackgroundFanart(id string, sid int64) error {
	ftv, err := app.Fanart.GetMovieImages(fmt.Sprintf("%d", sid))
	if err != nil {
		return errors.Wrap(err, "getting fanart")
	}
	if len(ftv.Posters) == 0 {
		return errors.New("no posters")
	}

	if err := app.Workers.Enqueue(&TvdbUpdateSeriesImage{ID: id, Type: "background", Path: ftv.Posters[0].URL, Ratio: backgroundRatio}); err != nil {
		return errors.Wrap(err, "enqueuing series episodes")
	}

	return nil
}

func TvdbUpdateSeriesBackgroundTvdb(id string, sid int64) error {
	app.Log.Named("TvdbUpdateSeriesBackground").Info("updating series images: background")
	r, err := app.Tvdb.GetSeriesArtworks(sid, tvdb.String("eng"), tvdb.Int64(int64(3)))
	if err != nil {
		return errors.Wrap(err, "getting series artworks")
	}

	if r.Data == nil || len(r.Data.Artworks) == 0 {
		return errors.New("no data")
	}
	if len(r.Data.Artworks) == 0 {
		return errors.New("no artworks")
	}

	background := r.Data.Artworks[0]
	if err := app.Workers.Enqueue(&TvdbUpdateSeriesImage{ID: id, Type: "background", Path: tvdb.StringValue(background.Image), Ratio: backgroundRatio}); err != nil {
		return errors.Wrap(err, "enqueuing series episodes")
	}

	return nil
}

// TvdbUpdateSeriesImage
type TvdbUpdateSeriesImage struct {
	minion.WorkerDefaults[*TvdbUpdateSeriesImage]
	ID    string
	Type  string
	Path  string
	Ratio float32
}

func (j *TvdbUpdateSeriesImage) Kind() string { return "TvdbUpdateSeriesImage" }
func (j *TvdbUpdateSeriesImage) Work(ctx context.Context, job *minion.Job[*TvdbUpdateSeriesImage]) error {
	app.Log.Info("updating series image")

	input := job.Args
	remote := input.Path // tvdb images are full urls
	extension := filepath.Ext(input.Path)[1:]
	local := fmt.Sprintf("series-%s/%s", input.ID, input.Type)
	dest := fmt.Sprintf("%s/%s.%s", app.Config.DirectoriesImages, local, extension)
	thumb := fmt.Sprintf("%s/%s_thumb.%s", app.Config.DirectoriesImages, local, extension)

	if err := imageDownload(remote, dest); err != nil {
		return errors.Wrap(err, "downloading image")
	}

	height := 400
	width := int(float32(height) * input.Ratio)
	if err := imageResize(dest, thumb, width, height); err != nil {
		return errors.Wrap(err, "resizing image")
	}

	series := &Series{}
	if err := app.DB.Series.Find(input.ID, series); err != nil {
		return errors.Wrap(err, "finding movie")
	}

	var img *Path
	for _, p := range series.Paths {
		if string(p.Type) == input.Type {
			img = p
			break
		}
	}

	if img == nil {
		app.Log.Info("path not found")
		img = &Path{}
		series.Paths = append(series.Paths, img)
	}

	img.Type = primitive.Symbol(input.Type)
	img.Remote = remote
	img.Local = local
	img.Extension = extension

	if err := app.DB.Series.Update(series); err != nil {
		return errors.Wrap(err, "updating series")
	}

	return nil
}
