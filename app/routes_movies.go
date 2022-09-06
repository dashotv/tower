package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func MoviesIndex(c *gin.Context) {
	q := App().DB.Medium.Query()
	results, err := q.
		Where("_type", "Movie").
		Desc("created_at").Run()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
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
