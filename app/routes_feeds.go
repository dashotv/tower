package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *Application) FeedsIndex(c *gin.Context, page, limit int) {
	results, err := app.DB.Feed.Query().
		Desc("processed").
		Limit(1000).
		Run()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func (a *Application) FeedsCreate(c *gin.Context) {
	data := &Feed{}
	err := c.BindJSON(&data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = app.DB.Feed.Save(data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "id": data.ID.Hex(), "feed": data})
}

func (a *Application) FeedsShow(c *gin.Context, id string) {
	result := &Feed{}
	err := app.DB.Feed.Find(id, result)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (a *Application) FeedsUpdate(c *gin.Context, id string) {
	data := &Feed{}
	err := c.BindJSON(&data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = app.DB.FeedUpdate(id, data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false})
}

func (a *Application) FeedsSettings(c *gin.Context, id string) {
	data := &Setting{}
	err := c.BindJSON(&data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = app.DB.FeedSetting(id, data.Setting, data.Value)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false})
}

func (a *Application) FeedsDelete(c *gin.Context, id string) {
	c.JSON(http.StatusOK, gin.H{"error": false})
}
