package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *Application) UpcomingIndex(c *gin.Context) {
	episodes, err := app.DB.Upcoming()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, episodes)
}
