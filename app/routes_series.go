package app

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/dashotv/golem/web"
)

const pagesize = 25

func SeriesIndex(c *gin.Context) {
	page, err := web.QueryDefaultInteger(c, "page", 1)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	count, err := App().DB.Series.Count(bson.M{"_type": "Series"})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	q := App().DB.Series.Query()
	results, err := q.
		Where("_type", "Series").
		Skip((page - 1) * pagesize).
		Desc("created_at").Run()

	for _, s := range results {
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
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func SeriesShow(c *gin.Context, id string) {
	result := &Medium{}
	App().Log.Infof("series.show id=%s", id)
	err := App().DB.Medium.Find(id, result)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

	c.JSON(http.StatusOK, result)
}

func SeriesUpdate(c *gin.Context, id string) {
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func SeriesDelete(c *gin.Context, id string) {
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func SeriesSeasons(c *gin.Context, id string) {
	results, err := App().DB.SeriesSeasons(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func SeriesSeasonEpisodes(c *gin.Context, id string, season string) {
	results, err := App().DB.SeriesSeasonEpisodes(id, season)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
