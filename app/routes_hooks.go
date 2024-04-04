package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GET /hooks/plex
func (a *Application) HooksPlex(c echo.Context) error {
	// TODO: implement the route
	return c.JSON(http.StatusNotImplemented, &Response{Error: false, Message: "not implmented"})
}

// POST /hooks/nzbget
func (a *Application) HooksNzbget(c echo.Context, payload *NzbgetPayload) error {
	// TODO: implement the route
	return c.JSON(http.StatusNotImplemented, &Response{Error: false, Message: "not implmented"})
}
