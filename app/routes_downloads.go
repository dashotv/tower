package app

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *Application) DownloadsIndex(c echo.Context, page, limit int) error {
	results, err := app.DB.ActiveDownloads()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, results)
}

func (a *Application) DownloadsLast(c echo.Context) error {
	var t int
	_, err := app.Cache.Get("seer_downloads", &t)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, gin.H{"error": false, "last": t})
}

type DownloadRequest struct {
	MediumId string `json:"medium_id"`
}

func (a *Application) DownloadsCreate(c echo.Context) error {
	data := &DownloadRequest{}
	err := c.Bind(data)
	if err != nil {
		return err
	}

	if data.MediumId == "" {
		return errors.New("medium_id is required")
	}

	id, err := primitive.ObjectIDFromHex(data.MediumId)
	if err != nil {
		return err
	}

	d := &Download{MediumId: id, Status: "searching"}
	err = app.DB.Download.Save(d)
	if err != nil {
		return err
	}

	m := &Medium{}
	err = app.DB.Medium.Find(data.MediumId, m)
	if err != nil {
		return err
	}

	m.Downloaded = true
	err = app.DB.Medium.Update(m)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, gin.H{"error": false, "id": d.ID.Hex()})
}

func (a *Application) DownloadsShow(c echo.Context, id string) error {
	result := &Download{}
	err := app.DB.Download.Find(id, result)
	if err != nil {
		return err
	}

	list := []*Download{result}
	app.DB.processDownloads(list)

	return c.JSON(http.StatusOK, result)
}

func (a *Application) DownloadsUpdate(c echo.Context, id string) error {
	data := &Download{}
	err := c.Bind(data)
	if err != nil {
		return err
	}

	err = app.DB.Download.Update(data)
	if err != nil {
		return err
	}

	if data.Status == "deleted" {
		m := &Medium{}
		err = app.DB.Medium.Find(data.MediumId.Hex(), m)
		if err != nil {
			return err
		}

		m.Downloaded = false
		err = app.DB.Medium.Update(m)
		if err != nil {
			return err
		}
	} else if data.Status == "loading" && (data.Url != "" || data.ReleaseId != "") {
		if err := app.Workers.Enqueue(&DownloadsProcess{}); err != nil {
			return err
		}
	}

	return c.JSON(http.StatusOK, data)
}

func (a *Application) DownloadsSettings(c echo.Context, id string) error {
	data := &Setting{}
	err := c.Bind(data)
	if err != nil {
		return err
	}

	err = app.DB.DownloadSetting(id, data.Setting, data.Value)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, gin.H{"errors": false, "data": data})
}

func (a *Application) DownloadsDelete(c echo.Context, id string) error {
	return c.JSON(http.StatusOK, gin.H{"error": false})
}

func (a *Application) DownloadsRecent(c echo.Context) error {
	mid := QueryString(c, "medium_id")
	page, err := QueryDefaultInteger(c, "page", 1)
	if err != nil {
		return err
	}

	results, total, err := app.DB.RecentDownloads(mid, page)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, gin.H{"count": total, "results": results})
}

type DownloadSelector struct {
	MediumId string
	Num      int
}

func (a *Application) DownloadsSelect(c echo.Context, id string) error {
	data := &DownloadSelector{}
	err := c.Bind(&data)
	if err != nil {
		return err
	}

	err = app.DB.DownloadSelect(id, data.MediumId, data.Num)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, gin.H{"errors": false, "data": data})
}

func (a *Application) DownloadsMedium(c echo.Context, id string) error {
	download := &Download{}
	err := app.DB.Download.Find(id, download)
	if err != nil {
		return err
	}

	list := []*Download{download}
	app.DB.processDownloads(list)

	if download.Medium == nil {
		return c.JSON(http.StatusOK, gin.H{"errors": false})
	}

	if download.Medium.Type == "Series" {
		return a.SeriesSeasonEpisodesAll(c, download.MediumId.Hex())
	}

	return c.JSON(http.StatusOK, []*Medium{download.Medium})
}

var thashIsTorrent = regexp.MustCompile(`^[a-f0-9]{40}$`)

func (a *Application) DownloadsTorrent(c echo.Context, id string) error {
	download := &Download{}
	err := app.DB.Download.Find(id, download)
	if err != nil {
		return err
	}

	if download.Thash == "" || thashIsTorrent.MatchString(download.Thash) == false {
		return c.JSON(http.StatusOK, gin.H{"errors": false, "message": "No torrent hash available"})
	}

	torrent, err := app.Flame.Torrent(download.Thash)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, torrent)
}
