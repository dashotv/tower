package importer

import "go.uber.org/zap"

type Options struct {
	Logger   *zap.SugaredLogger
	Language string

	TvdbKey string
	// TvdbToken string
	// TvdbURL      string

	TmdbToken string
	// TmdbURL      string
	TmdbImageURL string

	FanartKey string
	FanartURL string
}

var DefaultOptions = &Options{
	Language:     "eng",
	TmdbImageURL: "https://image.tmdb.org/t/p/original",
}
