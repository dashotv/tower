package app

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dashotv/golem/web"
)

const releasePageSize = 25

func ReleasesIndex(c *gin.Context) {
	page, err := web.QueryDefaultInteger(c, "page", 1)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results, err := App().DB.Release.Query().
		Desc("created_at").
		Limit(releasePageSize).Skip((page - 1) * releasePageSize).
		Run()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func ReleasesCreate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func ReleasesShow(c *gin.Context, id string) {
	result := &Release{}
	err := App().DB.Release.Find(id, result)
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

func ReleasesUpdate(c *gin.Context, id string) {
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func ReleasesDelete(c *gin.Context, id string) {
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func ReleasesSetting(c *gin.Context, id string) {
	s := &Setting{}
	err := c.BindJSON(s)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = App().DB.ReleaseSetting(id, s.Setting, s.Value)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"errors": false, "data": s})
}
