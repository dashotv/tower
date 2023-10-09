package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func FeedsIndex(c *gin.Context) {
	results, err := db.Feed.Query().
		Desc("processed").
		Limit(1000).
		Run()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func FeedsCreate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func FeedsShow(c *gin.Context, id string) {
	result := &Feed{}
	err := db.Feed.Find(id, result)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func FeedsUpdate(c *gin.Context, id string) {
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func FeedsDelete(c *gin.Context, id string) {
	c.JSON(http.StatusOK, gin.H{"error": false})
}
