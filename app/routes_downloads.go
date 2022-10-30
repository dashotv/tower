package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func DownloadsIndex(c *gin.Context) {
	results, err := App().DB.ActiveDownloads()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func DownloadsCreate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func DownloadsShow(c *gin.Context, id string) {
	result := &Download{}
	err := App().DB.Download.Find(id, result)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	m := &Medium{}
	err = App().DB.Medium.FindByID(result.MediumId, m)
	if err != nil {
		App().Log.Errorf("could not find medium: %s", result.MediumId)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	App().Log.Infof("found %s: %s", m.ID, m.Title)
	result.Medium = *m

	c.JSON(http.StatusOK, result)
}

func DownloadsUpdate(c *gin.Context, id string) {
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func DownloadsDelete(c *gin.Context, id string) {
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func DownloadsRecent(c *gin.Context) {
	results, err := App().DB.RecentDownloads()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
