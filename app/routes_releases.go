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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
