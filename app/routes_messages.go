package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *Application) MessagesIndex(c *gin.Context) {
	list, err := app.DB.Message.Query().Desc("created_at").Limit(250).Run()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}
