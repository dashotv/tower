package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GET /upcoming/
func (a *Application) UpcomingIndex(c echo.Context) error {
	list, err := a.DB.Upcoming()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: list})
}
func (a *Application) UpcomingLater(c echo.Context) error {
	list, err := a.DB.UpcomingLater()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: list})
}
func (a *Application) UpcomingNow(c echo.Context) error {
	list, err := a.DB.UpcomingNow()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: list})
}
