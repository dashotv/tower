package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func WatchesIndex(c *gin.Context, mediumId string) {
	watches, err := db.Watch.Query().Desc("watched_at").Run()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, w := range watches {
		m := &Medium{}
		if err := db.Medium.FindByID(w.MediumId, m); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		w.Medium = m
	}

	c.JSON(http.StatusOK, watches)
}
