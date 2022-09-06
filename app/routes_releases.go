package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ReleasesIndex(c *gin.Context) {
	results, err := App().DB.Release.Query().Run()
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
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func ReleasesUpdate(c *gin.Context, id string) {
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func ReleasesDelete(c *gin.Context, id string) {
	c.JSON(http.StatusOK, gin.H{"error": false})
}
