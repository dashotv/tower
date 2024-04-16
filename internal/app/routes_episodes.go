package app

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PATCH /episodes/:id
func (a *Application) EpisodesSettings(c echo.Context, id string, data *Setting) error {
	err := app.DB.EpisodeSetting(id, data.Name, data.Value)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, &Response{Error: false, Result: data})
}

// PUT /episodes/:id
func (a *Application) EpisodesUpdate(c echo.Context, id string, subject *Episode) error {
	// TODO: process the subject

	// if you need to copy or compare to existing object...
	// data, err := a.DB.EpisodeGet(id)
	// if err != nil {
	//     return c.JSON(http.StatusNotFound, &Response{Error: true, Message: "not found"})
	// }
	// data.Name = subject.Name ...
	if err := a.DB.Episode.Save(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error saving Episodes"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

func (a *Application) EpisodesSettingsBatch(c echo.Context, settings *SettingsBatch) error {
	if settings.Name == "watched" {
		for _, sid := range settings.IDs {
			id, err := primitive.ObjectIDFromHex(sid)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: err.Error()})
			}
			// TODO: need CreateOrUpdate method
			if err := app.DB.Watch.Save(&Watch{MediumID: id, Username: "someone", WatchedAt: time.Now()}); err != nil {
				return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: err.Error()})
			}
		}
	} else {
		_, err := app.DB.Episode.Collection.UpdateMany(context.Background(), bson.M{"_id": bson.M{"$in": settings.IDs}}, bson.M{"$set": bson.M{settings.Name: settings.Value}})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: err.Error()})
		}
	}

	return c.JSON(http.StatusOK, &Response{Error: false, Result: settings})
}
