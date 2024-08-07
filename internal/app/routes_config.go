package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// PATCH /config/:id
func (a *Application) ConfigSettings(c echo.Context, id string, data *Setting) error {
	switch data.Name {
	case "runic":
		a.Config.ProcessRunicEvents = data.Value
	default:
		return c.JSON(http.StatusBadRequest, &Response{Error: true, Message: "invalid setting: " + data.Name})
	}

	return c.JSON(http.StatusOK, &Response{Error: false, Result: data})
}
