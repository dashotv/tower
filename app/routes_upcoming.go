package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpcomingIndex(c *gin.Context) {
	//episodes, err := App().DB.()
	//if err != nil {
	//	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}

	c.JSON(http.StatusOK, []*Medium{})
}
