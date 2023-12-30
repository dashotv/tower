package app

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/dashotv/tmdb"
	"github.com/dashotv/tvdb"
)

func init() {
	initializers = append(initializers, setupTvdb)
}

func setupTvdb(app *Application) error {
	c, err := tvdb.Login(app.Config.TvdbKey)
	if err != nil {
		return errors.Wrap(err, "tvdb login")
	}
	app.Tvdb = c
	return nil
}

func (a *Application) TvdbSeriesCovers(id int64) ([]string, error) {
	out := make([]string, 0)

	if resp, err := a.Tvdb.GetSeriesArtworks(id, nil, tvdb.Int64(2)); err == nil {
		for _, v := range resp.Data.Artworks {
			out = append(out, tvdb.StringValue(v.Image))
		}
	} else {
		a.Log.Warnf("failed to get series artworks: %s", err)
	}

	if list, err := a.TvdbSeriesCoversFanart(fmt.Sprintf("%d", id)); err == nil {
		out = append(out, list...)
	} else {
		a.Log.Warnf("fanart series images: %s", err)
	}

	if list, err := a.TvdbSeriesCoversTmdb(id); err == nil {
		out = append(out, list...)
	} else {
		a.Log.Warnf("tmdb series images: %s", err)
	}

	return out, nil
}

func (a *Application) TvdbSeriesCoversTmdb(id int64) ([]string, error) {
	find, err := app.Tmdb.FindByID(fmt.Sprintf("%d", id), "tvdb_id", tmdb.String("en-US"))
	if err != nil {
		return nil, errors.Wrap(err, "getting series")
	}
	if find.TvResults == nil || len(find.TvResults) == 0 {
		return nil, errors.New("not found getting series by tvdbid")
	}

	res := find.TvResults[0].(map[string]interface{})
	found := int(res["id"].(float64))

	resp, err := a.Tmdb.TvSeriesImages(found, nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "getting series images")
	}

	out := make([]string, 0)
	for _, v := range resp.Posters {
		if v.FilePath == nil {
			continue
		}
		u := fmt.Sprintf("%s%s", a.Config.TmdbImages, tmdb.StringValue(v.FilePath))
		out = append(out, u)
	}
	return out, nil
}

func (a *Application) TvdbSeriesCoversFanart(id string) ([]string, error) {
	resp, err := a.Fanart.GetShowImages(id)
	if err != nil {
		return nil, errors.Wrap(err, "getting fanart")
	}

	out := make([]string, 0)
	for _, v := range resp.Posters {
		out = append(out, v.URL)
	}
	return out, nil
}

func (a *Application) TvdbSeriesBackgrounds(id int64) ([]string, error) {
	out := make([]string, 0)

	if resp, err := a.Tvdb.GetSeriesArtworks(id, nil, tvdb.Int64(3)); err == nil {
		for _, v := range resp.Data.Artworks {
			out = append(out, tvdb.StringValue(v.Image))
		}
	} else {
		a.Log.Warnf("failed to get series artworks: %s", err)
	}

	if list, err := a.TvdbSeriesBackgroundsFanart(fmt.Sprintf("%d", id)); err == nil {
		out = append(out, list...)
	} else {
		a.Log.Warnf("fanart series images: %s", err)
	}

	if list, err := a.TvdbSeriesBackgroundsTmdb(id); err == nil {
		out = append(out, list...)
	} else {
		a.Log.Warnf("tmdb series images: %s", err)
	}

	return out, nil
}

func (a *Application) TvdbSeriesBackgroundsTmdb(id int64) ([]string, error) {
	find, err := app.Tmdb.FindByID(fmt.Sprintf("%d", id), "tvdb_id", tmdb.String("en-US"))
	if err != nil {
		return nil, errors.Wrap(err, "getting series")
	}
	if find.TvResults == nil || len(find.TvResults) == 0 {
		return nil, errors.New("not found getting series by tvdbid")
	}

	res := find.TvResults[0].(map[string]interface{})
	found := int(res["id"].(float64))

	resp, err := a.Tmdb.TvSeriesImages(found, nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "getting series images")
	}

	out := make([]string, 0)
	for _, v := range resp.Backdrops {
		if v.FilePath == nil {
			continue
		}
		u := fmt.Sprintf("%s%s", a.Config.TmdbImages, tmdb.StringValue(v.FilePath))
		out = append(out, u)
	}
	return out, nil
}

func (a *Application) TvdbSeriesBackgroundsFanart(id string) ([]string, error) {
	resp, err := a.Fanart.GetShowImages(id)
	if err != nil {
		return nil, errors.Wrap(err, "getting fanart")
	}

	out := make([]string, 0)
	for _, v := range resp.Backgrounds {
		out = append(out, v.URL)
	}
	return out, nil
}
