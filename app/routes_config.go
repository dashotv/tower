package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GET /config/
func (a *Application) ConfigIndex(c echo.Context, page int, limit int) error {
	return c.JSON(http.StatusOK, H{"error": false, "config": a.Config})
}

// POST /config/
func (a *Application) ConfigCreate(c echo.Context) error {
	// TODO: implement the route
	return c.JSON(http.StatusNotImplemented, H{"error": "not implmented"})
	// return c.JSON(http.StatusOK, H{"error": false})
}

// GET /config/:id
func (a *Application) ConfigShow(c echo.Context, id string) error {
	// subject, err := a.DB.Config.Get(id)
	// if err != nil {
	//     return c.JSON(http.StatusNotFound, H{"error": true, "message": "not found"})
	// }

	// TODO: implement the route
	return c.JSON(http.StatusNotImplemented, H{"error": "not implmented"})
	// return c.JSON(http.StatusOK, H{"error": false})
}

// PUT /config/:id
func (a *Application) ConfigUpdate(c echo.Context, id string) error {
	// subject, err := a.DB.Config.Get(id)
	// if err != nil {
	//     return c.JSON(http.StatusNotFound, H{"error": true, "message": "not found"})
	// }

	// TODO: implement the route
	return c.JSON(http.StatusNotImplemented, H{"error": "not implmented"})
	// return c.JSON(http.StatusOK, H{"error": false})
}

// PATCH /config/:id
func (a *Application) ConfigSettings(c echo.Context, id string) error {
	data := &Setting{}
	err := c.Bind(data)
	if err != nil {
		return err
	}

	switch data.Setting {
	case "runic":
		a.Config.ProcessRunicEvents = data.Value
	default:
		return c.JSON(http.StatusBadRequest, H{"error": true, "message": "invalid setting: " + data.Setting})
	}

	return c.JSON(http.StatusOK, H{"error": false})
}

// DELETE /config/:id
func (a *Application) ConfigDelete(c echo.Context, id string) error {
	// subject, err := a.DB.Config.Get(id)
	// if err != nil {
	//     return c.JSON(http.StatusNotFound, H{"error": true, "message": "not found"})
	// }

	// TODO: implement the route
	return c.JSON(http.StatusNotImplemented, H{"error": "not implmented"})
	// return c.JSON(http.StatusOK, H{"error": false})
}
