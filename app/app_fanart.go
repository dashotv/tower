package app

import (
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

var fanart *Fanart

func setupFanart() error {
	fanart = NewFanart(cfg.FanartApiURL, cfg.FanartApiKey)
	return nil
}

type Fanart struct {
	ApiKey string
	c      *resty.Client
}

func NewFanart(url, apiKey string) *Fanart {
	return &Fanart{
		ApiKey: apiKey,
		c:      resty.New().SetBaseURL(url),
	}
}

type FanartImage struct {
	ID    string `json:"id"`
	URL   string `json:"url"`
	Lang  string `json:"lang"`
	Likes string `json:"likes"`
}

type FanartMovieImages struct {
	Name        string        `json:"name"`
	TmdbID      string        `json:"tmdb_id"`
	ImdbID      string        `json:"imdb_id"`
	Posters     []FanartImage `json:"movieposter"`
	Backgrounds []FanartImage `json:"moviebackground"`
	Status      string        `json:"status"`
	Erorr       string        `json:"error message"`
}

func (f *Fanart) GetMovieImages(id string) (*FanartShowImages, error) {
	res := &FanartShowImages{}
	resp, err := f.c.R().
		SetQueryParam("api_key", f.ApiKey).
		SetResult(res).
		Get("/movies/" + id)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.Errorf("fanart: %s", resp.Status())
	}
	return res, nil
}

type FanartShowImages struct {
	Name        string        `json:"name"`
	TvdbID      string        `json:"thetvdb_id"`
	Posters     []FanartImage `json:"tvposter"`
	Backgrounds []FanartImage `json:"showbackground"`
	Status      string        `json:"status"`
	Erorr       string        `json:"error message"`
}

func (f *Fanart) GetShowImages(id string) (*FanartShowImages, error) {
	res := &FanartShowImages{}
	resp, err := f.c.R().
		SetQueryParam("api_key", f.ApiKey).
		SetResult(res).
		Get("/tv/" + id)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.Errorf("fanart: %s", resp.Status())
	}
	if res.Status == "error" {
		return nil, errors.Errorf("fanart: %s", res.Erorr)
	}
	return res, nil
}
