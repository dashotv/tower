package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (a *Application) UpcomingIndex(c echo.Context) error {
	episodes, err := app.DB.Upcoming()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, episodes)
}
