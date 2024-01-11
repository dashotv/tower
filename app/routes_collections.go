package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *Application) CollectionsIndex(c *gin.Context, page int, limit int) {
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 25
	}

	list, err := a.DB.CollectionList(limit, (page-1)*limit)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "count": len(list), "results": list})
}

func (a *Application) CollectionsCreate(c *gin.Context) {
	col := &Collection{}
	err := c.BindJSON(col)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = a.DB.Collection.Save(col)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "id": col.ID, "collection": col})
}

func (a *Application) CollectionsShow(c *gin.Context, id string) {
	subject, err := a.DB.CollectionGet(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subject)
}

func (a *Application) CollectionsUpdate(c *gin.Context, id string) {
	subject := &Collection{}

	if err := c.BindJSON(subject); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := a.DB.Collection.Save(subject); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(subject.Media) > 0 {
		if err := a.Workers.Enqueue(&PlexCollectionUpdate{Id: subject.ID.Hex()}); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
	})
}

func (a *Application) CollectionsSettings(c *gin.Context, id string) {
	// asssuming this is a CRUD route, get the subject from the database
	// subject, err := a.DB.Collections.Get(id)
	c.JSON(http.StatusOK, gin.H{
		"error": false,
	})
}

func (a *Application) CollectionsDelete(c *gin.Context, id string) {
	col, err := a.DB.Collection.Get(id, &Collection{})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if col.RatingKey != "" {
		if err := app.Plex.DeleteCollection(col.RatingKey); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	err = a.DB.Collection.Delete(col)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false})
}
