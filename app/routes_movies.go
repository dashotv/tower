package app

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/dashotv/golem/web"
)

func MoviesIndex(c *gin.Context) {
	page, err := web.QueryDefaultInteger(c, "page", 1)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	count, err := App().DB.Series.Count(bson.M{"_type": "Movie"})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	q := App().DB.Medium.Query()
	results, err := q.
		Where("_type", "Movie").
		Skip((page - 1) * pagesize).
		Desc("created_at").Run()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, m := range results {
		for _, p := range m.Paths {
			if p.Type == "cover" {
				m.Cover = fmt.Sprintf("%s/%s.%s", imagesBaseURL, p.Local, p.Extension)
				continue
			}
			if p.Type == "background" {
				m.Background = fmt.Sprintf("%s/%s.%s", imagesBaseURL, p.Local, p.Extension)
				continue
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"count": count, "results": results})
}

func MoviesCreate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func MoviesShow(c *gin.Context, id string) {
	result := &Medium{}
	err := App().DB.Medium.Find(id, result)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func MoviesUpdate(c *gin.Context, id string) {
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func MoviesDelete(c *gin.Context, id string) {
	c.JSON(http.StatusOK, gin.H{"error": false})
}
