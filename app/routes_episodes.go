package app

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *Application) EpisodesUpdate(c echo.Context, id string) error {
	data := &Setting{}
	err := c.Bind(&data)
	if err != nil {
		return err
	}

	err = app.DB.EpisodeSetting(id, data.Name, data.Value)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, gin.H{"errors": false, "data": data})
}

func (a *Application) EpisodesSetting(c echo.Context, id string) error {
	data := &Setting{}
	err := c.Bind(&data)
	if err != nil {
		return err
	}

	err = app.DB.EpisodeSetting(id, data.Name, data.Value)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, gin.H{"errors": false, "data": data})
}

type EpisodeSettingsBatch struct {
	IDs   []primitive.ObjectID `json:"ids"`
	Field string               `json:"field"`
	Value bool                 `json:"value"`
}

func (a *Application) EpisodesSettings(c echo.Context) error {
	data := &EpisodeSettingsBatch{}
	err := c.Bind(data)
	if err != nil {
		return err
	}

	_, err = app.DB.Episode.Collection.UpdateMany(context.Background(), bson.M{"_id": bson.M{"$in": data.IDs}}, bson.M{"$set": bson.M{data.Field: data.Value}})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, gin.H{"errors": false})
}
