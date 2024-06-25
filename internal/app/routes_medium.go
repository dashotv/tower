package app

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/dashotv/fae"
)

// GET /medium/medium/:id
func (a *Application) MediumShow(c echo.Context, id string) error {
	subject, err := a.DB.Medium.Get(id, &Medium{})
	if err != nil {
		return fae.Wrap(err, "finding medium")
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}
