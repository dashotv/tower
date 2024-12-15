package app

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/dashotv/fae"
)

// GET /watches/
func (a *Application) WatchesIndex(c echo.Context, medium_id string, username string) error {
	list, err := app.DB.Watches(medium_id, username)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: list})
}

// POST /watches/
func (a *Application) WatchesCreate(c echo.Context, medium_id string, username string) error {
	m, err := app.DB.Medium.Get(medium_id, &Medium{})
	if err != nil {
		return err
	}

	w, err := app.DB.WatchGet(m.ID, username)
	if err != nil {
		return err
	}
	if w != nil {
		return c.JSON(http.StatusOK, &Response{Error: false, Message: "watch exists", Result: w})
	}

	watch := &Watch{}
	if username == app.Config.PlexUsername {
		sw, err := app.DB.WatchGet(m.ID, "someone") // generic user from the UI
		if err != nil {
			return err
		}
		if sw != nil {
			watch = sw
		}
	}

	watch.MediumID = m.ID
	watch.Username = username
	watch.WatchedAt = time.Now()

	if err := app.DB.Watch.Save(watch); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &Response{Error: false, Result: watch})
}

// DELETE /watches/:id
func (a *Application) WatchesDelete(c echo.Context, id string) error {
	w, err := app.DB.Watch.Get(id, &Watch{})
	if err != nil {
		return err
	}
	if w == nil {
		return err
	}

	if err := app.DB.Watch.Delete(w); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &Response{Error: false, Message: "watch deleted"})
}

// DELETE /watches/medium
func (a *Application) WatchesDeleteMedium(c echo.Context, medium_id string) error {
	if medium_id == "" {
		return fae.New("medium_id is required")
	}

	mid, err := primitive.ObjectIDFromHex(medium_id)
	if err != nil {
		return fae.New("invalid medium_id")
	}

	if _, err := app.DB.Watch.Collection.DeleteMany(context.Background(), bson.M{"medium_id": mid}); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &Response{Error: false, Message: "watch deleted"})
}
