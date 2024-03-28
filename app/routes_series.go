package app

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const pagesize = 42

func (a *Application) SeriesIndex(c echo.Context, page, limit int) error {
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = pagesize
	}

	kind := QueryString(c, "type")
	source := QueryString(c, "source")
	active := QueryBool(c, "active")
	favorite := QueryBool(c, "favorite")
	broken := QueryBool(c, "broken")

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
		return err
	}

	results, err := q.
		Limit(limit).
		Skip((page - 1) * limit).
		Desc("created_at").Run()
	if err != nil {
		return err
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

	return c.JSON(http.StatusOK, gin.H{"count": count, "results": results})
}

func (a *Application) SeriesCreate(c echo.Context) error {
	r := &CreateRequest{}
	c.Bind(r)
	if r.ID == "" || r.Source == "" {
		return errors.New("id and source are required")
	}

	a.Log.Debugf("series create: %+v", r)
	s := &Series{
		Type:         "Series",
		SourceId:     r.ID,
		Source:       r.Source,
		Title:        r.Title,
		Description:  r.Description,
		Kind:         primitive.Symbol(r.Kind),
		SearchParams: &SearchParams{Resolution: 1080, Verified: true, Type: "tv"},
	}

	if isAnimeKind(string(r.Kind)) {
		s.SearchParams.Type = "anime"
	}

	d, err := time.Parse("2006-01-02", r.Date)
	if err != nil {
		a.Log.Debugf("error parsing date: %s", err.Error())
		s.ReleaseDate = time.Unix(0, 0)
	} else {
		s.ReleaseDate = d
	}

	err = app.DB.Series.Save(s)
	if err != nil {
		return err
	}

	if err := app.Workers.Enqueue(&SeriesUpdate{ID: s.ID.Hex()}); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, gin.H{"error": false, "series": s})
}

func (a *Application) SeriesShow(c echo.Context, id string) error {
	result := &Series{}
	app.Log.Infof("series.show id=%s", id)
	// cache this? have to figure out how to handle breaking cache
	err := app.DB.Series.Find(id, result)
	if err != nil {
		return err
	}

	unwatched, err := app.DB.SeriesUserUnwatched(result)
	if err != nil {
		return err
	}
	result.Unwatched = unwatched
	unwatchedall, err := app.DB.SeriesUnwatched(result, "")
	if err != nil {
		return err
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
		return err
	}

	//Seasons
	result.Seasons, err = app.DB.SeriesSeasons(id)
	if err != nil {
		return err
	}

	//CurrentSeason
	result.CurrentSeason, err = app.DB.SeriesCurrentSeason(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

func (a *Application) SeriesUpdate(c echo.Context, id string) error {
	data := &Series{}
	err := c.Bind(&data)
	if err != nil {
		return err
	}

	if !strings.HasPrefix(data.Cover, "/media-images") {
		cover := data.GetCover()
		if cover == nil || cover.Remote != data.Cover {
			if err := app.Workers.Enqueue(&SeriesImage{ID: id, Type: "cover", Path: data.Cover, Ratio: posterRatio}); err != nil {
				return err
			}
		}
	}

	if !strings.HasPrefix(data.Background, "/media-images") {
		background := data.GetBackground()
		if background == nil || background.Remote != data.Background {
			if err := app.Workers.Enqueue(&SeriesImage{ID: id, Type: "background", Path: data.Background, Ratio: backgroundRatio}); err != nil {
				return err
			}
		}
	}

	err = app.DB.SeriesUpdate(id, data)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, gin.H{"errors": false, "data": data})
}

func (a *Application) SeriesSettings(c echo.Context, id string) error {
	data := &Setting{}
	err := c.Bind(&data)
	if err != nil {
		return err
	}

	err = app.DB.SeriesSetting(id, data.Setting, data.Value)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, gin.H{"errors": false, "data": data})
}

func (a *Application) SeriesDelete(c echo.Context, id string) error {
	return c.JSON(http.StatusOK, gin.H{"error": false})
}

func (a *Application) SeriesCurrentSeason(c echo.Context, id string) error {
	i, err := app.DB.SeriesCurrentSeason(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, gin.H{"current": i})
}

func (a *Application) SeriesSeasons(c echo.Context, id string) error {
	results, err := app.DB.SeriesSeasons(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, results)
}

func (a *Application) SeriesSeasonEpisodesAll(c echo.Context, id string) error {
	results, err := app.DB.SeriesSeasonEpisodesAll(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, results)
}

func (a *Application) SeriesSeasonEpisodes(c echo.Context, id string, season string) error {
	results, err := app.DB.SeriesSeasonEpisodes(id, season)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, results)
}

func (a *Application) SeriesPaths(c echo.Context, id string) error {
	results, err := app.DB.SeriesPaths(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, results)
}

func (a *Application) SeriesWatches(c echo.Context, id string) error {
	results, err := app.DB.SeriesWatches(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, results)
}

func (a *Application) SeriesRefresh(c echo.Context, id string) error {
	if err := app.Workers.Enqueue(&SeriesUpdate{ID: id}); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, gin.H{"error": false})
}

func (a *Application) SeriesCovers(c echo.Context, id string) error {
	series, err := a.DB.Series.Get(id, &Series{})
	if err != nil {
		return errors.Wrap(err, "getting series")
	}

	if series == nil {
		return errors.New("series not found")
	}

	if series.Source != "tvdb" {
		return errors.New("series not from tvdb")
	}

	tvdbid, err := strconv.Atoi(series.SourceId)
	if err != nil {
		return errors.Wrap(err, "converting tvdb id")
	}

	resp, err := app.TvdbSeriesCovers(int64(tvdbid))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, gin.H{"error": false, "covers": resp})
}

func (a *Application) SeriesBackgrounds(c echo.Context, id string) error {
	series, err := a.DB.Series.Get(id, &Series{})
	if err != nil {
		errors.Wrap(err, "getting series")
	}

	if series == nil {
		return errors.New("series not found")
	}

	if series.Source != "tvdb" {
		return errors.New("series not from tvdb")
	}

	tvdbid, err := strconv.Atoi(series.SourceId)
	if err != nil {
		errors.Wrap(err, "converting tvdb id")
	}

	resp, err := app.TvdbSeriesBackgrounds(int64(tvdbid))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, gin.H{"error": false, "backgrounds": resp})
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
		return fmt.Errorf("unknown job: %s", name)
	}
}

func (a *Application) SeriesJobs(c echo.Context, id string, name string) error {
	if err := seriesJob(name, id); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, gin.H{"error": false})
}
