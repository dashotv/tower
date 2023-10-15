package app

import "github.com/gin-gonic/gin"

func JobsIndex(c *gin.Context) {
	list, err := db.MinionJob.Query().Desc("created_at").Run()
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"jobs": list})
}
