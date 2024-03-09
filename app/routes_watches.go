package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (a *Application) WatchesIndex(c echo.Context, mediumId, username string) error {
	watches, err := app.DB.Watches(mediumId, username)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, watches)
}
