package app

import "github.com/gin-gonic/gin"

func RequestsIndex(c *gin.Context) {
	list, err := db.Request.Query().Run()
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, list)
}

func RequestsShow(c *gin.Context, id string) {
	c.JSON(200, gin.H{"message": "RequestsShow"})
}
