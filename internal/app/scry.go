package app

import (
	"context"
	"strings"

	"github.com/samber/lo"

	"github.com/dashotv/fae"
	runic "github.com/dashotv/runic/client"
	scry "github.com/dashotv/scry/client"
)

func init() {
	initializers = append(initializers, setupScry)
}

func setupScry(app *Application) error {
	app.Scry = scry.New(app.Config.ScryURL)
	return nil
}

func (a *Application) ScrySearchMovie(search *DownloadSearch) (*runic.Release, error) {
	if search == nil {
		return nil, fae.New("search is nil")
	}
	req := &scry.RunicIndexRequest{
		Type:       search.Type,
		Text:       search.Title,
		Group:      search.Group,
		Source:     search.Source,
		Uncensored: search.Uncensored,
		Bluray:     search.Bluray,
		Exact:      search.Exact,
		Verified:   true,
		Year:       -1,
		Season:     -1,
		Episode:    -1,
		Resolution: -1,
	}
	if search.Year > 0 {
		req.Year = search.Year
	}
	if search.Resolution > 0 {
		req.Resolution = search.Resolution
	}

	resp, err := a.Scry.Runic.Index(context.Background(), req)
	if err != nil {
		return nil, fae.Wrap(err, "failed to search releases")
	}

	// app.Log.Named("search").Warnf("ScrySearchMovie(): %s (%d) %02dx%02d => %d search: %s\n", search.Title, search.Year, search.Season, search.Episode, len(resp.Result.Releases), resp.Result.Search)
	if len(resp.Result.Releases) == 0 {
		return nil, nil
	}

	return selectRunicRelease(search, resp.Result.Releases)
}

func (a *Application) ScrySearchEpisode(search *DownloadSearch) (*runic.Release, error) {
	if search == nil {
		return nil, fae.New("search is nil")
	}
	req := &scry.RunicIndexRequest{
		Type:       search.Type,
		Text:       search.Title,
		Group:      search.Group,
		Website:    search.Website,
		Source:     search.Source,
		Uncensored: search.Uncensored,
		Bluray:     search.Bluray,
		Verified:   search.Verified,
		Exact:      search.Exact,
		Year:       -1,
		Season:     -1,
		Episode:    -1,
		Resolution: -1,
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

	resp, err := a.Scry.Runic.Index(context.Background(), req)
	if err != nil {
		return nil, fae.Wrap(err, "failed to search releases")
	}

	// app.Log.Named("search").Warnf("ScrySearchEpisode(): %s (%d) %02dx%02d => %d search: %s\n", search.Title, search.Year, search.Season, search.Episode, len(resp.Result.Releases), resp.Result.Search)
	if len(resp.Result.Releases) == 0 {
		return nil, nil
	}

	return selectRunicRelease(search, resp.Result.Releases)
}

func selectRunicRelease(s *DownloadSearch, releases []*runic.Release) (*runic.Release, error) {
	c := &RunicChooser{
		Group: s.Group,
		Title: s.Title,
		Exact: s.Exact,
		data: map[string]map[string][]*runic.Release{
			"nzbs": {
				"preferred": {},
				"good":      {},
			},
			"tors": {
				"preferred": {},
				"good":      {},
			},
		},
	}

	for _, r := range releases {
		c.add(r)
	}

	return c.choose(), nil
}

type RunicChooser struct {
	Group string
	Title string
	Exact bool
	data  map[string]map[string][]*runic.Release
}

func (c *RunicChooser) add(r *runic.Release) {
	k := "tors"
	if r.Downloader == "nzb" {
		k = "nzbs"
	}

	if lo.Contains(app.Config.DownloadsPreferred, strings.ToLower(r.Website)) ||
		lo.Contains(app.Config.DownloadsPreferred, strings.ToLower(r.Group)) {
		c.data[k]["preferred"] = append(c.data[k]["preferred"], r)
	}
	if lo.Contains(app.Config.DownloadsGroups, strings.ToLower(r.Website)) ||
		lo.Contains(app.Config.DownloadsGroups, strings.ToLower(r.Group)) {
		c.data[k]["good"] = append(c.data[k]["good"], r)
	}
}

func (c *RunicChooser) choose() *runic.Release {
	// app.Log.Debugf("chooser: %+v", c.data)
	if len(c.data["nzbs"]["preferred"]) > 0 {
		return c.data["nzbs"]["preferred"][0]
	}
	if len(c.data["nzbs"]["good"]) > 0 {
		return c.data["nzbs"]["good"][0]
	}
	if len(c.data["tors"]["preferred"]) > 0 {
		return c.data["tors"]["preferred"][0]
	}
	if len(c.data["tors"]["good"]) > 0 {
		if !c.Exact {
			return c.data["tors"]["good"][0]
		}

		for _, r := range c.data["tors"]["good"] {
			if c.Title == r.Title {
				return r
			}
		}
	}

	// 	if r == nil {
	// 		return nil
	// 	}
	//
	// 	if c.Group == r.Group {
	// 		return r
	// 	}
	// 	if c.Exact && r.Name == c.Title {
	// 		return r
	// 	}

	return nil
}
