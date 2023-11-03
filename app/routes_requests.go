package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequestsIndex(c *gin.Context) {
	list, err := db.Request.Query().Desc("created_at").Run()
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func RequestsShow(c *gin.Context, id string) {
	c.JSON(http.StatusOK, gin.H{"message": "RequestsShow"})
}
