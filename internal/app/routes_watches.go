package app

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// GET /watches/
func (a *Application) WatchesIndex(c echo.Context, medium_id string, username string) error {
	list, err := app.DB.Watches(medium_id, username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error loading Watches"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: list})
}

// POST /watches/
func (a *Application) WatchesCreate(c echo.Context, medium_id string, username string) error {
	m, err := app.DB.Medium.Get(medium_id, &Medium{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "loading medium:" + err.Error()})
	}

	w, err := app.DB.WatchGet(m.ID, username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "loading watch:" + err.Error()})
	}
	if w != nil {
		return c.JSON(http.StatusOK, &Response{Error: false, Message: "watch exists", Result: w})
	}

	watch := &Watch{
		MediumID:  m.ID,
		Username:  username,
		WatchedAt: time.Now(),
	}

	if err := app.DB.Watch.Save(watch); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "saving:" + err.Error()})
	}

	return c.JSON(http.StatusOK, &Response{Error: false, Result: watch})
}
