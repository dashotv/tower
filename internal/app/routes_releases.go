package app

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/dashotv/fae"
)

// GET /releases/
func (a *Application) ReleasesIndex(c echo.Context, page int, limit int) error {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 25
	}
	skip := (page - 1) * limit
	list, err := a.DB.Release.Query().
		Desc("published_at").
		Desc("created_at").
		Limit(limit).Skip(skip).
		Run()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error loading Releases"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: list})
}

// POST /releases/
func (a *Application) ReleasesCreate(c echo.Context, subject *Release) error {
	if err := a.DB.Release.Save(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error saving Releases"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// GET /releases/:id
func (a *Application) ReleasesShow(c echo.Context, id string) error {
	subject := &Release{}
	err := app.DB.Release.Find(id, subject)
	if err != nil {
		return c.JSON(http.StatusNotFound, &Response{Error: true, Message: "not found"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// PUT /releases/:id
func (a *Application) ReleasesUpdate(c echo.Context, id string, subject *Release) error {
	// if you need to copy or compare to existing object...
	// data, err := a.DB.ReleaseGet(id)
	// if err != nil {
	//     return c.JSON(http.StatusNotFound, &Response{Error: true, Message: "not found"})
	// }
	// data.Name = subject.Name ...
	if err := a.DB.Release.Save(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error saving Releases"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// PATCH /releases/:id
func (a *Application) ReleasesSettings(c echo.Context, id string, setting *Setting) error {
	err := app.DB.ReleaseSetting(id, setting.Name, setting.Value)
	if err != nil {
		return c.JSON(http.StatusNotFound, &Response{Error: true, Message: "not found"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: setting})
}

// DELETE /releases/:id
func (a *Application) ReleasesDelete(c echo.Context, id string) error {
	subject := &Release{}
	err := a.DB.Release.Find(id, subject)
	if err != nil {
		return c.JSON(http.StatusNotFound, &Response{Error: true, Message: "not found"})
	}
	if err := a.DB.Release.Delete(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error deleting Releases"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// GET /releases/popular/:interval
func (a *Application) ReleasesPopular(c echo.Context, interval string) error {
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

	return c.JSON(http.StatusOK, &Response{Error: false, Result: out})
}
