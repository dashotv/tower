package importer

import (
	"fmt"

	"github.com/sourcegraph/conc"

	"github.com/dashotv/fae"
	"github.com/dashotv/tmdb"
	"github.com/dashotv/tvdb"
)

func (i *Importer) loadSeriesImages(tvdbid int64) ([]string, []string, error) {
	var covers []string
	var backgrounds []string

	tmdbid, err := i.TmdbID(tvdbid)
	if err != nil {
		i.Log.Warnf("images: %s", err)
	}

	wg := conc.NewWaitGroup()
	wg.Go(func() {
		covers, err = i.loadSeriesCovers(tvdbid, tmdbid)
		if err != nil {
			i.Log.Warnf("covers: %s", err)
		}
	})

	wg.Go(func() {
		backgrounds, err = i.loadSeriesBackgrounds(tvdbid, tmdbid)
		if err != nil {
			i.Log.Warnf("backgrounds: %s", err)
		}
	})

	wg.Wait()
	if err != nil {
		return nil, nil, err
	}

	return covers, backgrounds, nil
}

func (i *Importer) loadSeriesCovers(tvdbid int64, tmdbid int) ([]string, error) {
	covers := []string{}
	tvdbCovers := []string{}
	fanartCovers := []string{}
	tmdbCovers := []string{}

	wg := conc.NewWaitGroup()
	wg.Go(func() {
		var err error
		tvdbCovers, err = i.loadSeriesCoversTvdb(tvdbid)
		if err != nil {
			i.Log.Warnf("tvdb: %s", err)
		}
	})

	wg.Go(func() {
		var err error
		fanartCovers, err = i.loadSeriesCoversFanart(tvdbid)
		if err != nil {
			i.Log.Warnf("fanart: %s", err)
		}
	})

	if tmdbid != 0 {
		wg.Go(func() {
			var err error
			tmdbCovers, err = i.loadSeriesCoversTmdb(tmdbid)
			if err != nil {
				i.Log.Warnf("tmdb: %s", err)
			}
		})
	}

	wg.Wait()

	covers = append(covers, tvdbCovers...)
	covers = append(covers, fanartCovers...)
	covers = append(covers, tmdbCovers...)
	return covers, nil
}

func (i *Importer) loadSeriesCoversTvdb(tvdbid int64) ([]string, error) {
	covers := []string{}

	r, err := i.Tvdb.GetSeriesArtworks(float64(tvdbid), nil, tvdb.Int64(int64(2)))
	if err != nil {
		return nil, fae.Wrap(err, "covers")
	}

	if r.Data == nil || len(r.Data.Artworks) == 0 {
		return nil, fae.New("covers: no data")
	}

	for _, cover := range r.Data.Artworks {
		covers = append(covers, tvdb.StringValue(cover.Image))
	}

	return covers, nil
}

func (i *Importer) loadSeriesCoversFanart(tvdbid int64) ([]string, error) {
	covers := []string{}

	ftv, err := i.Fanart.GetShowImages(fmt.Sprintf("%d", tvdbid))
	if err != nil {
		return nil, fae.Wrap(err, "covers")
	}

	if len(ftv.Posters) == 0 {
		return nil, fae.New("covers: no data")
	}

	for _, poster := range ftv.Posters {
		covers = append(covers, poster.URL)
	}

	return covers, nil
}

func (i *Importer) loadSeriesCoversTmdb(tmdbid int) ([]string, error) {
	covers := []string{}

	resp, err := i.Tmdb.TvSeriesImages(tmdbid, nil, nil)
	if err != nil {
		return nil, fae.Wrap(err, "cover")
	}

	if resp.Posters == nil || len(resp.Posters) == 0 {
		return nil, fae.New("cover: no data")
	}

	for _, cover := range resp.Posters {
		covers = append(covers, i.Opts.TmdbImageURL+tmdb.StringValue(cover.FilePath))
	}

	return covers, nil
}

func (i *Importer) loadSeriesBackgrounds(tvdbid int64, tmdbid int) ([]string, error) {
	backgrounds := []string{}

	tvdbBackgrounds := []string{}
	fanartBackgrounds := []string{}
	tmdbBackgrounds := []string{}

	wg := conc.NewWaitGroup()
	wg.Go(func() {
		var err error
		tvdbBackgrounds, err = i.loadSeriesBackgroundsTvdb(tvdbid)
		if err != nil {
			i.Log.Warnf("tvdb: %s", err)
		}
	})

	wg.Go(func() {
		var err error
		fanartBackgrounds, err = i.loadSeriesBackgroundsFanart(tvdbid)
		if err != nil {
			i.Log.Warnf("fanart: %s", err)
		}
	})

	if tmdbid != 0 {
		wg.Go(func() {
			var err error
			tmdbBackgrounds, err = i.loadSeriesBackgroundsTmdb(tmdbid)
			if err != nil {
				i.Log.Warnf("tmdb: %s", err)
			}
		})
	}

	wg.Wait()

	backgrounds = append(backgrounds, tvdbBackgrounds...)
	backgrounds = append(backgrounds, fanartBackgrounds...)
	backgrounds = append(backgrounds, tmdbBackgrounds...)
	return backgrounds, nil
}

func (i *Importer) loadSeriesBackgroundsTvdb(tvdbid int64) ([]string, error) {
	backgrounds := []string{}

	r, err := i.Tvdb.GetSeriesArtworks(float64(tvdbid), nil, tvdb.Int64(int64(3)))
	if err != nil {
		return nil, fae.Wrap(err, "backgrounds")
	}

	if r.Data == nil || len(r.Data.Artworks) == 0 {
		return nil, fae.New("backgrounds: no data")
	}

	for _, background := range r.Data.Artworks {
		backgrounds = append(backgrounds, tvdb.StringValue(background.Image))
	}

	return backgrounds, nil
}

func (i *Importer) loadSeriesBackgroundsFanart(tvdbid int64) ([]string, error) {
	backgrounds := []string{}

	ftv, err := i.Fanart.GetShowImages(fmt.Sprintf("%d", tvdbid))
	if err != nil {
		return nil, fae.Wrap(err, "backgrounds")
	}

	if len(ftv.Backgrounds) == 0 {
		return nil, fae.New("backgrounds: no data")
	}

	for _, fanart := range ftv.Backgrounds {
		backgrounds = append(backgrounds, fanart.URL)
	}

	return backgrounds, nil
}

func (i *Importer) loadSeriesBackgroundsTmdb(tmdbid int) ([]string, error) {
	backgrounds := []string{}

	resp, err := i.Tmdb.TvSeriesImages(tmdbid, nil, nil)
	if err != nil {
		return nil, fae.Wrap(err, "background")
	}

	if resp.Backdrops == nil || len(resp.Backdrops) == 0 {
		return nil, fae.New("background: no data")
	}

	for _, background := range resp.Backdrops {
		backgrounds = append(backgrounds, i.Opts.TmdbImageURL+tmdb.StringValue(background.FilePath))
	}

	return backgrounds, nil
}
