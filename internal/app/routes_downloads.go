package app

import (
	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/dashotv/fae"
)

func (a *Application) DownloadsIndex(c echo.Context, page, limit int) error {
	results := []*Download{}
	if ok, err := a.Cache.Get("downloads", &results); err != nil || !ok {
		return fae.Errorf("getting downloads: %w", err)
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: results})
}

func (a *Application) DownloadsLast(c echo.Context) error {
	var t int
	_, err := a.Cache.Get("seer_downloads", &t)
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
	err := a.DB.Download.Save(data)
	if err != nil {
		return err
	}

	m := &Medium{}
	err = a.DB.Medium.FindByID(data.MediumID, m)
	if err != nil {
		return err
	}

	m.Downloaded = true
	err = a.DB.Medium.Update(m)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &Response{Error: false, Result: data})
}

func (a *Application) DownloadsShow(c echo.Context, id string) error {
	results := []*Download{}
	if ok, err := a.Cache.Get("downloads", &results); err != nil || !ok {
		return fae.Errorf("getting downloads: %w", err)
	}

	for _, d := range results {
		if d.ID.Hex() == id {
			return c.JSON(http.StatusOK, &Response{Error: false, Result: d})
		}
	}

	result := &Download{}
	if err := a.DB.Download.Find(id, result); err != nil {
		return err
	}

	a.DB.processDownload(result)
	return c.JSON(http.StatusOK, &Response{Error: false, Result: result})
}

func (a *Application) DownloadsUpdate(c echo.Context, id string, data *Download) error {
	if id != data.ID.Hex() || id == primitive.NilObjectID.Hex() || data.ID == primitive.NilObjectID {
		return fae.New("ID mismatch")
	}
	err := a.DB.Download.Save(data)
	if err != nil {
		return err
	}

	if data.Status == "deleted" {
		if err := a.DB.MediumSetting(data.MediumID.Hex(), "downloaded", false); err != nil {
			return err
		}
		if data.Thash != "" {
			if err := a.FlameTorrentRemove(data.Thash); err != nil {
				return err
			}
		}
	} else if data.Status == "done" {
		if data.Thash != "" {
			if err := a.FlameTorrentRemove(data.Thash); err != nil {
				return err
			}
		}
	} else if data.Status == "loading" && (data.URL != "" || data.ReleaseID != "") {
		if err := a.Workers.Enqueue(&DownloadsProcessLoad{}); err != nil {
			return err
		}
	}

	return c.JSON(http.StatusOK, &Response{Error: false, Result: data})
}

func (a *Application) DownloadsSettings(c echo.Context, id string, data *Setting) error {
	err := a.DB.DownloadSetting(id, data.Name, data.Value)
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

	results, total, err := a.DB.RecentDownloads(mid, page)
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
	err := a.DB.DownloadSelect(id, medium_id, num)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &Response{Error: false})
}

func (a *Application) DownloadsMedium(c echo.Context, id string) error {
	download := &Download{}
	err := a.DB.Download.Find(id, download)
	if err != nil {
		return err
	}

	list := []*Download{download}
	a.DB.processDownloads(list)

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
	err := a.DB.Download.Find(id, download)
	if err != nil {
		return err
	}

	if download.Thash == "" || thashIsTorrent.MatchString(download.Thash) == false {
		return c.JSON(http.StatusOK, &Response{Error: false, Message: "No torrent hash available"})
	}

	torrent, err := a.FlameTorrent(download.Thash)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, torrent)
}
