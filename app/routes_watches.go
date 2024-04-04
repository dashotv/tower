package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GET /watches/
func (a *Application) WatchesIndex(c echo.Context, medium_id string, username string) error {
	list, err := app.DB.Watches(medium_id, username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error loading Watches"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: list})
}
