package app

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/samber/lo"

	"github.com/dashotv/fae"
	"github.com/dashotv/scry/search"
)

func init() {
	initializers = append(initializers, setupScry)
}

func setupScry(app *Application) error {
	app.Scry = &Scry{
		URL: "https://scry:10080",
		c:   resty.New().SetBaseURL("http://scry:10080"),
	}
	return nil
}

type Scry struct {
	URL string
	c   *resty.Client
}

type SearchOptions struct {
	Text       string
	Source     string
	Type       string
	Author     string
	Group      string
	Year       int
	Season     int
	Episode    int
	Resolution int
	Verified   bool
	Uncensored bool
	Bluray     bool
	Exact      bool
}

func selectRelease(opt *SearchOptions, releases []*search.Release) (*search.Release, error) {
	nzbs := lo.Filter(releases, func(item *search.Release, i int) bool {
		return item.NZB
	})
	for _, r := range nzbs {
		if Preferred(opt, r) {
			return r, nil
		}
	}
	for _, r := range nzbs {
		if Good(opt, r) {
			return r, nil
		}
	}

	tors := lo.Filter(releases, func(item *search.Release, i int) bool {
		return !item.NZB
	})
	for _, r := range tors {
		if Preferred(opt, r) {
			return r, nil
		}
	}
	for _, r := range tors {
		if Good(opt, r) {
			return r, nil
		}
	}

	return nil, nil
}

func Preferred(opt *SearchOptions, r *search.Release) bool {
	for _, g := range app.Config.DownloadsPreferred {
		if r.Group == g {
			return true
		}
	}

	return false
}

func Good(opt *SearchOptions, r *search.Release) bool {
	group := false
	for _, g := range app.Config.DownloadsGroups {
		if r.Group == g {
			group = true
			break
		}
	}
	if !group && opt.Group != r.Group && opt.Author != r.Group && opt.Author != r.Author {
		return false
	}

	if opt.Exact && opt.Text != r.Name {
		return false
	}

	return true
}

func (s *Scry) Search(options *SearchOptions) (*search.ReleaseSearchResponse, error) {
	r := &search.ReleaseSearchResponse{}
	req := s.c.R().
		SetHeader("Content-Type", "application/json").
		SetResult(r)
	if options.Text != "" {
		req.SetQueryParam("text", options.Text)
	}
	if options.Source != "" {
		req.SetQueryParam("source", options.Source)
	}
	if options.Type != "" {
		req.SetQueryParam("type", options.Type)
	}
	if options.Author != "" {
		req.SetQueryParam("author", options.Author)
	}
	if options.Group != "" {
		req.SetQueryParam("group", options.Group)
	}
	if options.Year > 0 {
		req.SetQueryParam("year", fmt.Sprintf("%d", options.Year))
	}
	if options.Season > 0 {
		req.SetQueryParam("season", fmt.Sprintf("%d", options.Season))
	}
	if options.Episode > 0 {
		req.SetQueryParam("episode", fmt.Sprintf("%d", options.Episode))
	}
	if options.Resolution > 0 {
		req.SetQueryParam("resolution", fmt.Sprintf("%d", options.Resolution))
	}
	if options.Verified {
		req.SetQueryParam("verified", "true")
	}
	if options.Uncensored {
		req.SetQueryParam("uncensored", "true")
	}
	if options.Bluray {
		req.SetQueryParam("bluray", "true")
	}
	if options.Exact {
		req.SetQueryParam("exact", "true")
	}

	resp, err := req.
		Get("/releases")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fae.New(resp.Status())
	}
	return r, nil
}

func (s *Scry) ScrySearchEpisode(ep *Medium) (*search.Release, error) {
	series := &Series{}
	err := app.DB.Series.FindByID(ep.SeriesID, series)
	if err != nil {
		return nil, err
	}

	opt := &SearchOptions{
		Text:  series.Search,
		Exact: true,
	}

	params := series.SearchParams
	if params == nil {
		return nil, fae.New("no search params")
	}

	switch params.Type {
	case "tv":
		opt.Type = "tv"
		opt.Season = ep.SeasonNumber
		opt.Episode = ep.EpisodeNumber
	case "anime":
	case "ecchi":
		opt.Type = "anime"
		opt.Episode = ep.AbsoluteNumber
	}

	if params.Verified {
		opt.Verified = true
	}
	if params.Uncensored {
		opt.Uncensored = params.Uncensored
	}
	if params.Bluray {
		opt.Bluray = params.Bluray
	}
	if params.Group != "" {
		opt.Group = params.Group
	}
	if params.Author != "" {
		opt.Author = params.Author
	}
	if params.Resolution > 0 {
		opt.Resolution = params.Resolution
	}
	if params.Source != "" {
		opt.Source = params.Source
	}

	resp, err := app.Scry.Search(opt)
	if err != nil {
		return nil, fae.Wrap(err, "failed to search releases")
	}
	if len(resp.Releases) == 0 {
		return nil, nil
	}

	return selectRelease(opt, resp.Releases)
}
