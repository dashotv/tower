package app

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"

	"github.com/dashotv/fae"
)

const releasePageSize = 25

func (a *Application) ReleasesIndex(c echo.Context, page, limit int) error {
	if page == 0 {
		page = 1
	}
	results, err := app.DB.Release.Query().
		Desc("published_at").
		Desc("created_at").
		Limit(releasePageSize).Skip((page - 1) * releasePageSize).
		Run()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, results)
}

func (a *Application) ReleasesCreate(c echo.Context, release *Release) error {
	return c.JSON(http.StatusNotImplemented, gin.H{"error": true})
}

func (a *Application) ReleasesShow(c echo.Context, id string) error {
	result := &Release{}
	err := app.DB.Release.Find(id, result)
	if err != nil {
		// if err.Error() == "mongo: no documents in result" {
		// 	c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		// } else {
		// 	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// }
		return err
	}

	return c.JSON(http.StatusOK, result)
}

func (a *Application) ReleasesUpdate(c echo.Context, id string, release *Release) error {
	return c.JSON(http.StatusNotImplemented, gin.H{"error": true})
}

func (a *Application) ReleasesDelete(c echo.Context, id string) error {
	return c.JSON(http.StatusOK, gin.H{"error": false})
}

func (a *Application) ReleasesSettings(c echo.Context, id string, s *Setting) error {
	err := app.DB.ReleaseSetting(id, s.Name, s.Value)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, gin.H{"errors": false, "result": s})
}

func (a *Application) ReleasesPopular(c echo.Context, interval string) error {
	app.Log.Infof("ReleasesPopular: interval: %s", interval)
	out := map[string][]*Popular{}

	for _, t := range releaseTypes {
		results := make([]*Popular, 25)
		ok, err := app.Cache.Get(fmt.Sprintf("releases_popular_%s_%s", interval, t), &results)
		if err != nil {
			return err
		}
		if !ok {
			return fae.New("http.StatusNotFound")
		}
		out[t] = results
	}

	return c.JSON(http.StatusOK, out)
}
