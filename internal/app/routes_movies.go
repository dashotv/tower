package app

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/dashotv/fae"
)

// GET /movies/
func (a *Application) MoviesIndex(c echo.Context, page int, limit int, kind, source string, completed, downloaded, broken bool) error {
	q := app.DB.Movie.Query()
	if kind != "" {
		q = q.Where("kind", kind)
	}
	if source != "" {
		q = q.Where("source", source)
	}
	if broken {
		q = q.Where("broken", true)
	}
	if completed {
		q = q.Where("completed", true)
	}
	if downloaded {
		q = q.Where("downloaded", true)
	}

	count, err := q.Count()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: err.Error()})
	}

	list, err := q.
		Limit(limit).
		Skip((page - 1) * limit).
		Desc("created_at").Run()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: err.Error()})
	}

	for _, m := range list {
		for _, p := range m.Paths {
			if p.Type == "cover" {
				m.Cover = fmt.Sprintf("%s/%s.%s", imagesBaseURL, p.Local, p.Extension)
				continue
			}
			if p.Type == "background" {
				m.Background = fmt.Sprintf("%s/%s.%s", imagesBaseURL, p.Local, p.Extension)
				continue
			}
		}
	}

	return c.JSON(http.StatusOK, &Response{Error: false, Total: count, Result: list})
}

// POST /movies/
func (a *Application) MoviesCreate(c echo.Context, subject *Movie) error {
	if subject.SourceID == "" || subject.Source == "" {
		return fae.New("id and source are required")
	}

	subject.Type = "Movie"
	subject.SearchParams = &SearchParams{Resolution: 1080, Verified: true, Type: "movies"}

	if err := a.DB.Movie.Save(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error saving Movies"})
	}
	if err := app.Workers.Enqueue(&TmdbUpdateMovie{ID: subject.ID.Hex()}); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// GET /movies/:id
func (a *Application) MoviesShow(c echo.Context, id string) error {
	subject, err := a.DB.MovieGet(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, &Response{Error: true, Message: "not found"})
	}
	for _, p := range subject.Paths {
		if p.Type == "cover" {
			subject.Cover = fmt.Sprintf("%s/%s.%s", imagesBaseURL, p.Local, p.Extension)
			continue
		}
		if p.Type == "background" {
			subject.Background = fmt.Sprintf("%s/%s.%s", imagesBaseURL, p.Local, p.Extension)
			continue
		}
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// PUT /movies/:id
func (a *Application) MoviesUpdate(c echo.Context, id string, subject *Movie) error {
	// TODO: process the subject

	// if you need to copy or compare to existing object...
	// data, err := a.DB.MovieGet(id)
	// if err != nil {
	//     return c.JSON(http.StatusNotFound, &Response{Error: true, Message: "not found"})
	// }
	// data.Name = subject.Name ...
	if err := a.DB.Movie.Update(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error saving Movies"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// PATCH /movies/:id
func (a *Application) MoviesSettings(c echo.Context, id string, setting *Setting) error {
	err := a.DB.MovieSetting(id, setting.Name, setting.Value)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, &Response{Error: false, Result: setting})
}

// DELETE /movies/:id
func (a *Application) MoviesDelete(c echo.Context, id string) error {
	subject, err := a.DB.MovieGet(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "not found"})
	}
	if err := a.Workers.Enqueue(&MovieDelete{ID: id}); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// PUT /movies/:id/refresh
func (a *Application) MoviesRefresh(c echo.Context, id string) error {
	if err := app.Workers.Enqueue(&TmdbUpdateMovie{ID: id}); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &Response{Error: false})
}

// GET /movies/:id/paths
func (a *Application) MoviesPaths(c echo.Context, id string) error {
	results, err := app.DB.MoviePaths(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: results})
}
