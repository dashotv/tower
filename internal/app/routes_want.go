package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GET /want/series
func (a *Application) WantSeries(c echo.Context, id string) error {
	wanted, err := a.Want.SeriesWanted(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, &Response{Error: false, Result: wanted})
}

// GET /want/movie
func (a *Application) WantMovie(c echo.Context, id string) error {
	wanted, err := a.Want.MovieWanted(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, &Response{Error: false, Result: wanted})
}
