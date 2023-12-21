package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const pagesize = 42

func (a *Application) SeriesIndex(c *gin.Context, page, limit int) {
	if page == 0 {
		page = 1
	}

	count, err := app.DB.Series.Count(bson.M{"_type": "Series"})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	q := app.DB.Series.Query()
	results, err := q.
		Limit(pagesize).
		Skip((page - 1) * pagesize).
		Desc("created_at").Run()

	for _, s := range results {
		unwatched, err := app.DB.SeriesUserUnwatched(s)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		s.Unwatched = unwatched

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

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count, "results": results})
}

func (a *Application) SeriesCreate(c *gin.Context) {
	r := &CreateRequest{}
	c.BindJSON(r)
	if r.ID == "" || r.Source == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "id and source are required"})
		return
	}

	app.Log.Debugf("series create: %+v", r)
	s := &Series{
		Type:         "Series",
		SourceId:     r.ID,
		Source:       r.Source,
		Title:        r.Title,
		Description:  r.Description,
		Kind:         primitive.Symbol(r.Kind),
		SearchParams: &SearchParams{Resolution: 1080, Verified: true, Type: "tv"},
	}

	if r.Kind == "anime" {
		s.SearchParams.Type = "anime"
	}

	d, err := time.Parse("2006-01-02", r.Date)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	s.ReleaseDate = d

	err = app.DB.Series.Save(s)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := app.Workers.Enqueue(&TvdbUpdateSeries{ID: s.ID.Hex(), Images: true, Paths: true, Episodes: true}); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "series": s})
}

func (a *Application) SeriesShow(c *gin.Context, id string) {
	result := &Series{}
	app.Log.Infof("series.show id=%s", id)
	// cache this? have to figure out how to handle breaking cache
	err := app.DB.Series.Find(id, result)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	unwatched, err := app.DB.SeriesUserUnwatched(result)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result.Unwatched = unwatched

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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Seasons
	result.Seasons, err = app.DB.SeriesSeasons(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//CurrentSeason
	result.CurrentSeason, err = app.DB.SeriesCurrentSeason(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (a *Application) SeriesUpdate(c *gin.Context, id string) {
	data := &Series{}
	err := c.BindJSON(&data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = app.DB.SeriesUpdate(id, data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"errors": false, "data": data})
}

func (a *Application) SeriesSettings(c *gin.Context, id string) {
	data := &Setting{}
	err := c.BindJSON(&data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = app.DB.SeriesSetting(id, data.Setting, data.Value)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"errors": false, "data": data})
}

func (a *Application) SeriesDelete(c *gin.Context, id string) {
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func (a *Application) SeriesCurrentSeason(c *gin.Context, id string) {
	i, err := app.DB.SeriesCurrentSeason(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"current": i})
}

func (a *Application) SeriesSeasons(c *gin.Context, id string) {
	results, err := app.DB.SeriesSeasons(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func (a *Application) SeriesSeasonEpisodesAll(c *gin.Context, id string) {
	results, err := app.DB.SeriesSeasonEpisodesAll(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func (a *Application) SeriesSeasonEpisodes(c *gin.Context, id string, season string) {
	results, err := app.DB.SeriesSeasonEpisodes(id, season)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func (a *Application) SeriesPaths(c *gin.Context, id string) {
	results, err := app.DB.SeriesPaths(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func (a *Application) SeriesWatches(c *gin.Context, id string) {
	results, err := app.DB.SeriesWatches(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func (a *Application) SeriesRefresh(c *gin.Context, id string) {
	if err := app.Workers.Enqueue(&TvdbUpdateSeries{ID: id, Images: true, Paths: true, Episodes: true}); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}
