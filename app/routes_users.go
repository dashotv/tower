package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *Application) UsersIndex(c *gin.Context) {
	users, err := app.DB.User.Query().Run()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}
