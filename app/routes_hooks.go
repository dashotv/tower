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
	if err := a.Workers.Enqueue(&NzbgetProcess{Payload: payload}); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, &Response{Error: false})
}
