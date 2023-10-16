package app

import (
	"net/http"

	"github.com/dashotv/golem/web"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func DownloadsIndex(c *gin.Context) {
	results, err := db.ActiveDownloads()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func DownloadsLast(c *gin.Context) {
	var t int
	_, err := cache.Get("seer_downloads", &t)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"last": t})
}

func DownloadsCreate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func DownloadsShow(c *gin.Context, id string) {
	result := &Download{}
	err := db.Download.Find(id, result)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	list := []*Download{result}
	db.processDownloads(list)

	c.JSON(http.StatusOK, result)
}

func DownloadsUpdate(c *gin.Context, id string) {
	data := &Download{}
	err := c.BindJSON(data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = db.Download.Update(data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func DownloadsSetting(c *gin.Context, id string) {
	data := &Setting{}
	err := c.BindJSON(data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = db.DownloadSetting(id, data.Setting, data.Value)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"errors": false, "data": data})
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

	count, err := db.Series.Count(bson.M{})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	q := db.Download.Query()
	results, err := q.Where("status", "done").
		Desc("updated_at").Desc("created_at").
		Skip((page - 1) * pagesize).
		Limit(pagesize).
		Run()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.processDownloads(results)

	c.JSON(http.StatusOK, gin.H{"count": count, "results": results})
}

type DownloadSelector struct {
	MediumId string
	Num      int
}

func DownloadsSelect(c *gin.Context, id string) {
	data := &DownloadSelector{}
	err := c.BindJSON(&data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = db.DownloadSelect(id, data.MediumId, data.Num)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"errors": false, "data": data})
}

func DownloadsMedium(c *gin.Context, id string) {
	download := &Download{}
	err := db.Download.Find(id, download)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	list := []*Download{download}
	db.processDownloads(list)

	if download.Medium == nil {
		c.JSON(http.StatusOK, gin.H{"errors": false})
		return
	}

	if download.Medium.Type == "Series" {
		SeriesSeasonEpisodesAll(c, download.MediumId.Hex())
		return
	}

	c.JSON(http.StatusOK, []*Medium{download.Medium})
}
