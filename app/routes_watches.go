package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func WatchesIndex(c *gin.Context, mediumId, username string) {
	watches, err := db.Watches(mediumId, username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, watches)
}
