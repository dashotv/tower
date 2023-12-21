package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
