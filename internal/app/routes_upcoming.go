package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GET /upcoming/
func (a *Application) UpcomingIndex(c echo.Context) error {
	list, err := a.DB.Upcoming()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error loading Upcoming"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: list})
}
