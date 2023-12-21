package app

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const releasePageSize = 25

func (a *Application) ReleasesIndex(c *gin.Context, page, limit int) {
	if page == 0 {
		page = 1
	}
	results, err := app.DB.Release.Query().
		Desc("published_at").
		Desc("created_at").
		Limit(releasePageSize).Skip((page - 1) * releasePageSize).
		Run()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func (a *Application) ReleasesCreate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func (a *Application) ReleasesShow(c *gin.Context, id string) {
	result := &Release{}
	err := app.DB.Release.Find(id, result)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, result)
}

func (a *Application) ReleasesUpdate(c *gin.Context, id string) {
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func (a *Application) ReleasesDelete(c *gin.Context, id string) {
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func (a *Application) ReleasesSettings(c *gin.Context, id string) {
	s := &Setting{}
	err := c.BindJSON(s)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = app.DB.ReleaseSetting(id, s.Setting, s.Value)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"errors": false, "data": s})
}

func (a *Application) ReleasesPopular(c *gin.Context, interval string) {
	app.Log.Infof("ReleasesPopular: interval: %s", interval)
	out := map[string][]*Popular{}

	for _, t := range releaseTypes {
		results := make([]*Popular, 25)
		ok, err := app.Cache.Get(fmt.Sprintf("releases_popular_%s_%s", interval, t), &results)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if !ok {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		out[t] = results
	}

	c.JSON(http.StatusOK, out)
}
