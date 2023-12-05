package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func JobsIndex(c *gin.Context, page, limit int) {
	if limit == 0 {
		limit = 25
	}
	skip := (page * limit) - limit
	list, err := db.Minion.Query().Skip(skip).Limit(limit).Desc("created_at").Run()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"jobs": list})
}

func JobsCreate(c *gin.Context, job string) {
	j := workersList[job]
	if j == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid job: " + job})
		return
	}

	log.Infof("Enqueuing job: %s", j.Kind())
	err := workers.Enqueue(j)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, j)
}

func MessagesIndex(c *gin.Context) {
	list, err := db.Message.Query().Desc("created_at").Limit(250).Run()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, list)
}
