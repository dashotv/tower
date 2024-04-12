package app

import (
	"context"

	"github.com/samber/lo"

	"github.com/dashotv/fae"
	scry "github.com/dashotv/scry/client"
	"github.com/dashotv/scry/search"
)

func init() {
	initializers = append(initializers, setupScry)
}

func setupScry(app *Application) error {
	app.Scry = scry.New(app.Config.ScryURL)
	return nil
}

func (a *Application) ScrySearchEpisode(search *DownloadSearch) (*search.Release, error) {
	req := &scry.ReleasesIndexRequest{
		Type:       search.Type,
		Text:       search.Title,
		Group:      search.Group,
		Source:     search.Source,
		Uncensored: search.Uncensored,
		Bluray:     search.Bluray,
		Verified:   search.Verified,
		Exact:      search.Exact,
	}
	if search.Year > 0 {
		req.Year = search.Year
	}
	if search.Season > 0 {
		req.Season = search.Season
	}
	if search.Episode > 0 {
		req.Episode = search.Episode
	}
	if search.Resolution > 0 {
		req.Resolution = search.Resolution
	}

	resp, err := a.Scry.Releases.Index(context.Background(), req)
	if err != nil {
		return nil, fae.Wrap(err, "failed to search releases")
	}
	if len(resp.Result.Releases) == 0 {
		return nil, nil
	}

	return selectRelease(search, resp.Result.Releases)
}

func selectRelease(s *DownloadSearch, releases []*search.Release) (*search.Release, error) {
	nzbs := lo.Filter(releases, func(item *search.Release, i int) bool {
		return item.NZB
	})
	for _, r := range nzbs {
		if Preferred(s, r) {
			return r, nil
		}
	}
	for _, r := range nzbs {
		if Good(s, r) {
			return r, nil
		}
	}

	tors := lo.Filter(releases, func(item *search.Release, i int) bool {
		return !item.NZB
	})
	for _, r := range tors {
		if Preferred(s, r) {
			return r, nil
		}
	}
	for _, r := range tors {
		if Good(s, r) {
			return r, nil
		}
	}

	return nil, nil
}

func Preferred(s *DownloadSearch, r *search.Release) bool {
	for _, g := range app.Config.DownloadsPreferred {
		if s.Group == g {
			return true
		}
	}

	return false
}

func Good(s *DownloadSearch, r *search.Release) bool {
	group := false
	for _, g := range app.Config.DownloadsGroups {
		if r.Group == g {
			group = true
			break
		}
	}

	if !group && s.Group != r.Group {
		return false
	}

	if s.Exact && s.Title != r.Name {
		return false
	}

	return true
}
