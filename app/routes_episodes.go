package app

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Setting struct {
	Setting string
	Value   bool
}

func (a *Application) EpisodesUpdate(c *gin.Context, id string) {
	data := &Setting{}
	err := c.BindJSON(&data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = app.DB.EpisodeSetting(id, data.Setting, data.Value)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"errors": false, "data": data})
}

func (a *Application) EpisodesSetting(c *gin.Context, id string) {
	data := &Setting{}
	err := c.BindJSON(&data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = app.DB.EpisodeSetting(id, data.Setting, data.Value)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"errors": false, "data": data})
}

type EpisodeSettingsBatch struct {
	IDs   []primitive.ObjectID `json:"ids"`
	Field string               `json:"field"`
	Value bool                 `json:"value"`
}

func (a *Application) EpisodesSettings(c *gin.Context) {
	data := &EpisodeSettingsBatch{}
	err := c.BindJSON(data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = app.DB.Episode.Collection.UpdateMany(context.Background(), bson.M{"_id": bson.M{"$in": data.IDs}}, bson.M{"$set": bson.M{data.Field: data.Value}})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"errors": false})
}
