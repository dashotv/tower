package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dashotv/golem/web"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const pagesize = 42

func SeriesIndex(c *gin.Context) {
	page, err := web.QueryDefaultInteger(c, "page", 1)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	count, err := db.Series.Count(bson.M{"_type": "Series"})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	q := db.Series.Query()
	results, err := q.
		Where("_type", "Series").
		Limit(pagesize).
		Skip((page - 1) * pagesize).
		Desc("created_at").Run()

	for _, s := range results {
		unwatched, err := db.SeriesUserUnwatched(s)
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

func SeriesCreate(c *gin.Context) {
	r := &CreateRequest{}
	c.BindJSON(r)
	if r.ID == "" || r.Source == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "id and source are required"})
		return
	}

	server.Log.Debugf("series create: %+v", r)
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

	err = db.Series.Save(s)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := workers.Enqueue(&TvdbUpdateSeries{s.ID.Hex()}); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "series": s})
}

func SeriesShow(c *gin.Context, id string) {
	result := &Series{}
	log.Infof("series.show id=%s", id)
	// cache this? have to figure out how to handle breaking cache
	err := db.Series.Find(id, result)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	unwatched, err := db.SeriesUserUnwatched(result)
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
	result.Paths, err = db.SeriesPaths(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Seasons
	result.Seasons, err = db.SeriesSeasons(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//CurrentSeason
	result.CurrentSeason, err = db.SeriesCurrentSeason(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func SeriesUpdate(c *gin.Context, id string) {
	data := &Series{}
	err := c.BindJSON(&data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = db.SeriesUpdate(id, data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"errors": false, "data": data})
}

func SeriesSetting(c *gin.Context, id string) {
	data := &Setting{}
	err := c.BindJSON(&data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = db.SeriesSetting(id, data.Setting, data.Value)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"errors": false, "data": data})
}

func SeriesDelete(c *gin.Context, id string) {
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func SeriesCurrentSeason(c *gin.Context, id string) {
	i, err := db.SeriesCurrentSeason(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"current": i})
}

func SeriesSeasons(c *gin.Context, id string) {
	results, err := db.SeriesSeasons(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func SeriesSeasonEpisodesAll(c *gin.Context, id string) {
	results, err := db.SeriesSeasonEpisodesAll(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func SeriesSeasonEpisodes(c *gin.Context, id string, season string) {
	results, err := db.SeriesSeasonEpisodes(id, season)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func SeriesPaths(c *gin.Context, id string) {
	results, err := db.SeriesPaths(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func SeriesWatches(c *gin.Context, id string) {
	results, err := db.SeriesWatches(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func SeriesRefresh(c *gin.Context, id string) {
	if err := workers.Enqueue(&TvdbUpdateSeries{id}); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}
