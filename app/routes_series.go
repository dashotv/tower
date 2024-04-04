package app

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/dashotv/fae"
)

// GET /series/
func (a *Application) SeriesIndex(c echo.Context, page, limit int, kind, source string, active, favorite, broken bool) error {
	q := app.DB.Series.Query()
	if kind != "" {
		q = q.Where("kind", kind)
	}
	if source != "" {
		q = q.Where("source", source)
	}
	if active {
		q = q.Where("active", true)
	}
	if favorite {
		q = q.Where("favorite", true)
	}
	if broken {
		q = q.Where("broken", true)
	}

	count, err := q.Count()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: err.Error()})
	}

	results, err := q.
		Limit(limit).
		Skip((page - 1) * limit).
		Desc("created_at").Run()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: err.Error()})
	}

	for _, s := range results {
		unwatched, err := app.DB.SeriesUserUnwatched(s)
		if err != nil {
			return err
		}
		s.Unwatched = unwatched
		unwatchedall, err := app.DB.SeriesUnwatched(s, "")
		if err != nil {
			return err
		}
		s.UnwatchedAll = unwatchedall

		for _, p := range s.Paths {
			if p.Type == "cover" {
				s.Cover = fmt.Sprintf("%s/%s.%s", imagesBaseURL, p.Local, p.Extension)
				continue
			}
			if p.Type == "background" {
				s.Background = fmt.Sprintf("%s/%s.%s", imagesBaseURL, p.Local, p.Extension)
				continue
			}
		}
	}

	return c.JSON(http.StatusOK, &Response{Error: false, Result: results, Total: count})
}

// POST /series/
func (a *Application) SeriesCreate(c echo.Context, subject *Series) error {
	if subject.SourceId == "" || subject.Source == "" {
		return fae.New("id and source are required")
	}

	subject.Type = "Series"
	subject.SearchParams = &SearchParams{Resolution: 1080, Verified: true, Type: "tv"}
	if isAnimeKind(string(subject.Kind)) {
		subject.SearchParams.Type = "anime"
	}

	if err := a.DB.Series.Save(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error saving Series"})
	}
	if err := app.Workers.Enqueue(&SeriesUpdate{ID: subject.ID.Hex()}); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error queueing Series"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// GET /series/:id
func (a *Application) SeriesShow(c echo.Context, id string) error {
	// TODO: cache this? have to figure out how to handle breaking cache
	result := &Series{}
	err := app.DB.Series.Find(id, result)
	if err != nil {
		return c.JSON(http.StatusNotFound, &Response{Error: true, Message: err.Error()})
	}

	unwatched, err := app.DB.SeriesUserUnwatched(result)
	if err != nil {
		return c.JSON(http.StatusNotFound, &Response{Error: true, Message: err.Error()})
	}
	result.Unwatched = unwatched
	unwatchedall, err := app.DB.SeriesUnwatched(result, "")
	if err != nil {
		return c.JSON(http.StatusNotFound, &Response{Error: true, Message: err.Error()})
	}
	result.UnwatchedAll = unwatchedall

	for _, p := range result.Paths {
		if p.Type == "cover" {
			result.Cover = fmt.Sprintf("%s/%s.%s", imagesBaseURL, p.Local, p.Extension)
			continue
		}
		if p.Type == "background" {
			result.Background = fmt.Sprintf("%s/%s.%s", imagesBaseURL, p.Local, p.Extension)
			continue
		}
	}

	//Paths
	result.Paths, err = app.DB.SeriesPaths(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, &Response{Error: true, Message: err.Error()})
	}

	//Seasons
	result.Seasons, err = app.DB.SeriesSeasons(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, &Response{Error: true, Message: err.Error()})
	}

	//CurrentSeason
	result.CurrentSeason, err = app.DB.SeriesCurrentSeason(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, &Response{Error: true, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, &Response{Error: false, Result: result})
}

// PUT /series/:id
func (a *Application) SeriesUpdate(c echo.Context, id string, subject *Series) error {
	if !strings.HasPrefix(subject.Cover, "/media-images") {
		cover := subject.GetCover()
		if cover == nil || cover.Remote != subject.Cover {
			if err := app.Workers.Enqueue(&SeriesImage{ID: id, Type: "cover", Path: subject.Cover, Ratio: posterRatio}); err != nil {
				return err
			}
		}
	}

	if !strings.HasPrefix(subject.Background, "/media-images") {
		background := subject.GetBackground()
		if background == nil || background.Remote != subject.Background {
			if err := app.Workers.Enqueue(&SeriesImage{ID: id, Type: "background", Path: subject.Background, Ratio: backgroundRatio}); err != nil {
				return err
			}
		}
	}

	if err := a.DB.Series.Save(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error saving Series"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// PATCH /series/:id
func (a *Application) SeriesSettings(c echo.Context, id string, setting *Setting) error {
	err := app.DB.SeriesSetting(id, setting.Name, setting.Value)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error saving Series setting"})
	}

	// switch Setting.Name {
	// case "something":
	//    subject.Something = Setting.Value
	// }

	return c.JSON(http.StatusOK, &Response{Error: false, Result: setting})
}

// DELETE /series/:id
func (a *Application) SeriesDelete(c echo.Context, id string) error {
	subject := &Series{}
	err := app.DB.Series.Find(id, subject)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: err.Error()})
	}
	if err := a.DB.Series.Delete(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error deleting Series"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// GET /series/:id/currentseason
func (a *Application) SeriesCurrentSeason(c echo.Context, id string) error {
	i, err := app.DB.SeriesCurrentSeason(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: i})
}

// GET /series/:id/paths
func (a *Application) SeriesPaths(c echo.Context, id string) error {
	results, err := app.DB.SeriesPaths(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: results})
}

// PUT /series/:id/refresh
func (a *Application) SeriesRefresh(c echo.Context, id string) error {
	if err := a.Workers.Enqueue(&SeriesUpdate{ID: id}); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, &Response{Error: false})
}

// GET /series/:id/seasons/all
func (a *Application) SeriesSeasonEpisodesAll(c echo.Context, id string) error {
	results, err := app.DB.SeriesSeasonEpisodesAll(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: results})
}

// GET /series/:id/seasons/:season
func (a *Application) SeriesSeasonEpisodes(c echo.Context, id string, season string) error {
	results, err := app.DB.SeriesSeasonEpisodes(id, season)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: results})
}

// GET /series/:id/watches
func (a *Application) SeriesWatches(c echo.Context, id string) error {
	results, err := app.DB.SeriesWatches(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: results})
}

// GET /series/:id/covers
func (a *Application) SeriesCovers(c echo.Context, id string) error {
	series, err := a.DB.Series.Get(id, &Series{})
	if err != nil {
		return fae.Wrap(err, "getting series")
	}

	if series == nil {
		return fae.New("series not found")
	}

	if series.Source != "tvdb" {
		return fae.New("series not from tvdb")
	}

	tvdbid, err := strconv.Atoi(series.SourceId)
	if err != nil {
		return fae.Wrap(err, "converting tvdb id")
	}

	resp, err := app.TvdbSeriesCovers(int64(tvdbid))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &Response{Error: false, Result: resp})
}

// GET /series/:id/backgrounds
func (a *Application) SeriesBackgrounds(c echo.Context, id string) error {
	series, err := a.DB.Series.Get(id, &Series{})
	if err != nil {
		return fae.Wrap(err, "getting series")
	}

	if series == nil {
		return fae.New("series not found")
	}

	if series.Source != "tvdb" {
		return fae.New("series not from tvdb")
	}

	tvdbid, err := strconv.Atoi(series.SourceId)
	if err != nil {
		return fae.Wrap(err, "converting tvdb id")
	}

	resp, err := app.TvdbSeriesBackgrounds(int64(tvdbid))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &Response{Error: false, Result: resp})
}

func seriesJob(name string, id string) error {
	switch name {
	case "refresh":
		return app.Workers.Enqueue(&SeriesUpdate{ID: id})
	case "paths":
		return app.Workers.Enqueue(&PathCleanup{ID: id})
	case "files":
		return app.Workers.Enqueue(&FileMatchMedium{ID: id})
	default:
		return fae.Errorf("unknown job: %s", name)
	}
}

// POST /series/:id/jobs
func (a *Application) SeriesJobs(c echo.Context, id string, name string) error {
	if err := seriesJob(name, id); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, &Response{Error: false})
}
