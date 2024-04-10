package app

import (
	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/dashotv/fae"
)

func (a *Application) DownloadsIndex(c echo.Context, page, limit int) error {
	results, err := app.DB.ActiveDownloads()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &Response{Error: false, Result: results})
}

func (a *Application) DownloadsLast(c echo.Context) error {
	var t int
	_, err := app.Cache.Get("seer_downloads", &t)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &Response{Error: false, Result: t})
}

type DownloadRequest struct {
	MediumID string `json:"medium_id"`
}

func (a *Application) DownloadsCreate(c echo.Context, data *Download) error {
	if data.MediumID == primitive.NilObjectID {
		return fae.New("medium_id is required")
	}

	data.Status = "searching"
	err := app.DB.Download.Save(data)
	if err != nil {
		return err
	}

	m := &Medium{}
	err = app.DB.Medium.FindByID(data.MediumID, m)
	if err != nil {
		return err
	}

	m.Downloaded = true
	err = app.DB.Medium.Update(m)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &Response{Error: false, Result: data})
}

func (a *Application) DownloadsShow(c echo.Context, id string) error {
	result := &Download{}
	err := app.DB.Download.Find(id, result)
	if err != nil {
		return err
	}

	list := []*Download{result}
	app.DB.processDownloads(list)

	return c.JSON(http.StatusOK, &Response{Error: false, Result: result})
}

func (a *Application) DownloadsUpdate(c echo.Context, id string, data *Download) error {
	if id != data.ID.Hex() || id == primitive.NilObjectID.Hex() || data.ID == primitive.NilObjectID {
		return fae.New("ID mismatch")
	}
	err := app.DB.Download.Save(data)
	if err != nil {
		return err
	}

	if data.Status == "deleted" {
		m := &Medium{}
		err = app.DB.Medium.Find(data.MediumID.Hex(), m)
		if err != nil {
			return err
		}

		m.Downloaded = false
		err = app.DB.Medium.Update(m)
		if err != nil {
			return err
		}
	} else if data.Status == "loading" && (data.URL != "" || data.ReleaseID != "") {
		if err := app.Workers.Enqueue(&DownloadsProcess{}); err != nil {
			return err
		}
	}

	return c.JSON(http.StatusOK, &Response{Error: false, Result: data})
}

func (a *Application) DownloadsSettings(c echo.Context, id string, data *Setting) error {
	err := app.DB.DownloadSetting(id, data.Name, data.Value)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &Response{Error: false, Result: data})
}

func (a *Application) DownloadsDelete(c echo.Context, id string) error {
	return c.JSON(http.StatusOK, &Response{Error: false})
}

func (a *Application) DownloadsRecent(c echo.Context, page int, mid string) error {
	if page < 1 {
		page = 1
	}

	results, total, err := app.DB.RecentDownloads(mid, page)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &Response{Total: total, Result: results})
}

type DownloadSelector struct {
	MediumID string
	Num      int
}

func (a *Application) DownloadsSelect(c echo.Context, id string, medium_id string, num int) error {
	err := app.DB.DownloadSelect(id, medium_id, num)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &Response{Error: false})
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
		return c.JSON(http.StatusOK, &Response{Error: false})
	}

	if download.Medium.Type == "Series" {
		return a.SeriesSeasonEpisodesAll(c, download.MediumID.Hex())
	}

	return c.JSON(http.StatusOK, &Response{Error: false, Result: []*Medium{download.Medium}})
}

var thashIsTorrent = regexp.MustCompile(`^[a-f0-9]{40}$`)

func (a *Application) DownloadsTorrent(c echo.Context, id string) error {
	download := &Download{}
	err := app.DB.Download.Find(id, download)
	if err != nil {
		return err
	}

	if download.Thash == "" || thashIsTorrent.MatchString(download.Thash) == false {
		return c.JSON(http.StatusOK, &Response{Error: false, Message: "No torrent hash available"})
	}

	torrent, err := app.Flame.Torrent(download.Thash)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, torrent)
}
