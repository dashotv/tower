package app

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *Application) DownloadsIndex(c *gin.Context, page, limit int) {
	results, err := app.DB.ActiveDownloads()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func (a *Application) DownloadsLast(c *gin.Context) {
	var t int
	_, err := app.Cache.Get("seer_downloads", &t)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "last": t})
}

type DownloadRequest struct {
	MediumId string `json:"medium_id"`
}

func (a *Application) DownloadsCreate(c *gin.Context) {
	data := &DownloadRequest{}
	err := c.BindJSON(data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if data.MediumId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "medium_id is required"})
		return
	}

	id, err := primitive.ObjectIDFromHex(data.MediumId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	d := &Download{MediumId: id, Status: "searching"}
	err = app.DB.Download.Save(d)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	m := &Medium{}
	err = app.DB.Medium.Find(data.MediumId, m)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	m.Downloaded = true
	err = app.DB.Medium.Update(m)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "id": d.ID.Hex()})
}

func (a *Application) DownloadsShow(c *gin.Context, id string) {
	result := &Download{}
	err := app.DB.Download.Find(id, result)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	list := []*Download{result}
	app.DB.processDownloads(list)

	c.JSON(http.StatusOK, result)
}

func (a *Application) DownloadsUpdate(c *gin.Context, id string) {
	data := &Download{}
	err := c.BindJSON(data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = app.DB.Download.Update(data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if data.Status == "deleted" {
		m := &Medium{}
		err = app.DB.Medium.Find(data.MediumId.Hex(), m)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		m.Downloaded = false
		err = app.DB.Medium.Update(m)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, data)
}

func (a *Application) DownloadsSettings(c *gin.Context, id string) {
	data := &Setting{}
	err := c.BindJSON(data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = app.DB.DownloadSetting(id, data.Setting, data.Value)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"errors": false, "data": data})
}

func (a *Application) DownloadsDelete(c *gin.Context, id string) {
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func (a *Application) DownloadsRecent(c *gin.Context) {
	mid := QueryString(c, "medium_id")
	page, err := QueryDefaultInteger(c, "page", 1)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results, total, err := app.DB.RecentDownloads(mid, page)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": total, "results": results})
}

type DownloadSelector struct {
	MediumId string
	Num      int
}

func (a *Application) DownloadsSelect(c *gin.Context, id string) {
	data := &DownloadSelector{}
	err := c.BindJSON(&data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = app.DB.DownloadSelect(id, data.MediumId, data.Num)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"errors": false, "data": data})
}

func (a *Application) DownloadsMedium(c *gin.Context, id string) {
	download := &Download{}
	err := app.DB.Download.Find(id, download)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	list := []*Download{download}
	app.DB.processDownloads(list)

	if download.Medium == nil {
		c.JSON(http.StatusOK, gin.H{"errors": false})
		return
	}

	if download.Medium.Type == "Series" {
		a.SeriesSeasonEpisodesAll(c, download.MediumId.Hex())
		return
	}

	c.JSON(http.StatusOK, []*Medium{download.Medium})
}

var thashIsTorrent = regexp.MustCompile(`^[a-f0-9]{40}$`)

func (a *Application) DownloadsTorrent(c *gin.Context, id string) {
	download := &Download{}
	err := app.DB.Download.Find(id, download)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if download.Thash == "" || thashIsTorrent.MatchString(download.Thash) == false {
		c.JSON(http.StatusOK, gin.H{"errors": false, "message": "No torrent hash available"})
		return
	}

	torrent, err := app.Flame.Torrent(download.Thash)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, torrent)
}
