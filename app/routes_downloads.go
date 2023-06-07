package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/dashotv/golem/web"
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

	list := []*Download{result}
	processDownloads(list)

	c.JSON(http.StatusOK, result)
}

func DownloadsUpdate(c *gin.Context, id string) {
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func DownloadsDelete(c *gin.Context, id string) {
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func DownloadsRecent(c *gin.Context) {
	page, err := web.QueryDefaultInteger(c, "page", 1)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	count, err := App().DB.Series.Count(bson.M{})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	q := App().DB.Download.Query()
	results, err := q.Where("status", "done").
		Desc("updated_at").Desc("created_at").
		Skip((page - 1) * pagesize).
		Limit(pagesize).
		Run()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	processDownloads(results)

	c.JSON(http.StatusOK, gin.H{"count": count, "results": results})
}
