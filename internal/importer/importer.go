package importer

import (
	"github.com/dashotv/tmdb"
	"github.com/dashotv/tower/internal/fanart"
	"github.com/dashotv/tvdb"
)

func New(opts *Options) (*Importer, error) {
	c, err := tvdb.Login(opts.TvdbKey)
	if err != nil {
		return nil, err
	}

	if opts.Language == "" {
		opts.Language = DefaultOptions.Language
	}
	if opts.TmdbImageURL == "" {
		opts.TmdbImageURL = DefaultOptions.TmdbImageURL
	}

	i := &Importer{
		Opts:   opts,
		Tmdb:   tmdb.New(opts.TmdbToken),
		Tvdb:   c,
		Fanart: fanart.New(opts.FanartURL, opts.FanartKey),
	}

	return i, nil
}

type Importer struct {
	Opts   *Options
	Tmdb   *tmdb.Client
	Tvdb   *tvdb.Client
	Fanart *fanart.Fanart
}

func (i *Importer) Series(tvdbid int64) (*Series, error) {
	return i.loadSeries(tvdbid)
}

func (i *Importer) SeriesUpdated(since int64) ([]int64, error) {
	return i.loadSeriesUpdated(since)
}

func (i *Importer) SeriesEpisodes(tvdbid int64, episodeOrder int) ([]*Episode, error) {
	return i.loadEpisodes(tvdbid)
}

func (i *Importer) SeriesImages(tvdbid int64) ([]string, []string, error) {
	return i.loadSeriesImages(tvdbid)
}
