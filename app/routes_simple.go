package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func JobsIndex(c *gin.Context) {
	list, err := db.MinionJob.Query().Desc("created_at").Run()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"jobs": list})
}

func MessagesIndex(c *gin.Context) {
	list, err := db.Message.Query().Desc("created_at").Limit(250).Run()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, list)
}
