package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"

	"github.com/dashotv/tower/internal/plex"
)

func (a *Application) HooksPlex(c echo.Context) error {
	data := &plex.HookData{}
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
	}

	return c.JSON(http.StatusOK, gin.H{"error": false})
}

func (a *Application) HooksNzbget(c echo.Context, p *NzbgetPayload) error {
	if p.Status != "SUCCESS" {
		return c.JSON(http.StatusOK, gin.H{"error": false})
	}

	if err := a.Workers.Enqueue(&NzbgetProcess{Payload: p}); err != nil {
		return c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": err.Error()})
	}

	return c.JSON(http.StatusOK, gin.H{"error": false})
}
