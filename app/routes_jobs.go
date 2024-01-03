package app

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/dashotv/minion"
)

func (a *Application) JobsIndex(c *gin.Context, page int, limit int) {
	if limit == 0 {
		limit = 25
	}
	skip := (page * limit) - limit
	list, err := app.DB.Minion.Query().Skip(skip).Limit(limit).Desc("created_at").Run()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"jobs": list})
}

func (a *Application) JobsCreate(c *gin.Context, job string) {
	j := workersList[job]
	if j == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid job: " + job})
		return
	}

	app.Log.Infof("Enqueuing job: %s", j.Kind())
	err := app.Workers.Enqueue(j)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, j)
}

func (a *Application) JobsDelete(c *gin.Context, id string, hard bool) {
	if id == string(minion.StatusPending) && !hard {
		filter := bson.M{"status": minion.StatusPending}
		if _, err := app.DB.Minion.Collection.UpdateMany(context.Background(), filter, bson.M{"$set": bson.M{"status": minion.StatusCancelled}}); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"error": false})
		return
	} else if id == string(minion.StatusFailed) && hard {
		filter := bson.M{"status": minion.StatusFailed}
		if _, err := app.DB.Minion.Collection.DeleteMany(context.Background(), filter); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"error": false})
		return
	} else if id == string(minion.StatusCancelled) && hard {
		filter := bson.M{"status": minion.StatusCancelled}
		if _, err := app.DB.Minion.Collection.DeleteMany(context.Background(), filter); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"error": false})
		return
	}

	j, err := app.DB.Minion.Get(id, &Minion{})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	j.Status = string(minion.StatusCancelled)
	if err := app.DB.Minion.Save(j); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false})
}
