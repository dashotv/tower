package app

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dashotv/tower/internal/plex"
)

func (a *Application) HooksPlex(c *gin.Context) {
	data := &plex.HookData{}
	if err := c.ShouldBindJSON(data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false})
}
