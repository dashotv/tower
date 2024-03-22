package importer

import (
	"errors"
	"fmt"

	"github.com/dashotv/tmdb"
	"github.com/dashotv/tvdb"
)

func (i *Importer) loadSeriesImages(tvdbid int64) ([]string, []string, error) {
	tmdbid, err := i.TmdbID(tvdbid)
	if err != nil {
		return nil, nil, fmt.Errorf("images: %w", err)
	}

	covers, err := i.loadSeriesCovers(tvdbid, tmdbid)
	if err != nil {
		return nil, nil, fmt.Errorf("covers: %w", err)
	}

	backgrounds, err := i.loadSeriesBackgrounds(tvdbid, tmdbid)
	if err != nil {
		return nil, nil, fmt.Errorf("backgrounds: %w", err)
	}

	return covers, backgrounds, nil
}

func (i *Importer) loadSeriesCovers(tvdbid int64, tmdbid int) ([]string, error) {
	covers := []string{}

	tvdbCovers, err := i.loadSeriesCoversTvdb(tvdbid)
	if err != nil {
		return nil, fmt.Errorf("tvdb: %w", err)
	}
	fanartCovers, err := i.loadSeriesCoversFanart(tvdbid)
	if err != nil {
		return nil, fmt.Errorf("fanart: %w", err)
	}
	tmdbCovers, err := i.loadSeriesCoversTmdb(tmdbid)
	if err != nil {
		return nil, fmt.Errorf("tmdb: %w", err)
	}

	covers = append(covers, tvdbCovers...)
	covers = append(covers, fanartCovers...)
	covers = append(covers, tmdbCovers...)
	return covers, nil
}

func (i *Importer) loadSeriesCoversTvdb(tvdbid int64) ([]string, error) {
	covers := []string{}

	r, err := i.Tvdb.GetSeriesArtworks(tvdbid, tvdb.String("eng"), tvdb.Int64(int64(2)))
	if err != nil {
		return nil, fmt.Errorf("covers: %w", err)
	}

	if r.Data == nil || len(r.Data.Artworks) == 0 {
		return nil, errors.New("covers: no data")
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
		return nil, fmt.Errorf("covers: %w", err)
	}

	if len(ftv.Posters) == 0 {
		return nil, errors.New("covers: no data")
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
		return nil, fmt.Errorf("cover: %w", err)
	}

	if resp.Posters == nil || len(resp.Posters) == 0 {
		return nil, errors.New("cover: no data")
	}

	for _, cover := range resp.Posters {
		covers = append(covers, i.Opts.TmdbImageURL+tmdb.StringValue(cover.FilePath))
	}

	return covers, nil
}

func (i *Importer) loadSeriesBackgrounds(tvdbid int64, tmdbid int) ([]string, error) {
	backgrounds := []string{}

	tvdbBackgrounds, err := i.loadSeriesBackgroundsTvdb(tvdbid)
	if err != nil {
		return nil, fmt.Errorf("tvdb: %w", err)
	}
	fanartBackgrounds, err := i.loadSeriesBackgroundsFanart(tvdbid)
	if err != nil {
		return nil, fmt.Errorf("fanart: %w", err)
	}
	tmdbBackgrounds, err := i.loadSeriesBackgroundsTmdb(tmdbid)
	if err != nil {
		return nil, fmt.Errorf("tmdb: %w", err)
	}

	backgrounds = append(backgrounds, tvdbBackgrounds...)
	backgrounds = append(backgrounds, fanartBackgrounds...)
	backgrounds = append(backgrounds, tmdbBackgrounds...)
	return backgrounds, nil
}

func (i *Importer) loadSeriesBackgroundsTvdb(tvdbid int64) ([]string, error) {
	backgrounds := []string{}

	r, err := i.Tvdb.GetSeriesArtworks(tvdbid, tvdb.String("eng"), tvdb.Int64(int64(3)))
	if err != nil {
		return nil, fmt.Errorf("backgrounds: %w", err)
	}

	if r.Data == nil || len(r.Data.Artworks) == 0 {
		return nil, errors.New("backgrounds: no data")
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
		return nil, fmt.Errorf("backgrounds: %w", err)
	}

	if len(ftv.Posters) == 0 {
		return nil, errors.New("backgrounds: no data")
	}

	for _, fanart := range ftv.Posters {
		backgrounds = append(backgrounds, fanart.URL)
	}

	return backgrounds, nil
}

func (i *Importer) loadSeriesBackgroundsTmdb(tmdbid int) ([]string, error) {
	backgrounds := []string{}

	resp, err := i.Tmdb.TvSeriesImages(tmdbid, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("background: %w", err)
	}

	if resp.Backdrops == nil || len(resp.Backdrops) == 0 {
		return nil, errors.New("background: no data")
	}

	for _, background := range resp.Backdrops {
		backgrounds = append(backgrounds, i.Opts.TmdbImageURL+tmdb.StringValue(background.FilePath))
	}

	return backgrounds, nil
}
