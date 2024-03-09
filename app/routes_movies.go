package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *Application) MoviesIndex(c echo.Context, page, limit int) error {
	if page == 0 {
		page = 1
	}

	kind := QueryString(c, "kind")
	source := QueryString(c, "source")
	completed := QueryBool(c, "completed")
	downloaded := QueryBool(c, "downloaded")
	broken := QueryBool(c, "broken")

	q := app.DB.Movie.Query()
	if kind != "" {
		q = q.Where("kind", kind)
	}
	if source != "" {
		q = q.Where("source", source)
	}
	if completed {
		q = q.Where("completed", true)
	}
	if downloaded {
		q = q.Where("downloaded", true)
	}
	if broken {
		q = q.Where("broken", true)
	}

	count, err := q.Count()
	if err != nil {
		return err
	}

	results, err := q.
		Limit(pagesize).
		Skip((page - 1) * pagesize).
		Desc("created_at").Run()
	if err != nil {
		return err
	}

	// TODO: do this with custom unmarshaling?
	for _, m := range results {
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

	return c.JSON(http.StatusOK, gin.H{"count": count, "results": results})
}

func (a *Application) MoviesCreate(c echo.Context) error {
	r := &CreateRequest{}
	c.Bind(r)
	if r.ID == "" || r.Source == "" {
		return errors.New("id and source are required")
	}

	m := &Movie{
		Type:         "Movie",
		SourceId:     r.ID,
		Source:       r.Source,
		Title:        r.Title,
		Description:  r.Description,
		Kind:         primitive.Symbol(r.Kind),
		SearchParams: &SearchParams{Type: "movies", Resolution: 1080, Verified: true},
	}

	d, err := time.Parse("2006-01-02", r.Date)
	if err != nil {
		return err
	}
	m.ReleaseDate = d

	err = app.DB.Movie.Save(m)
	if err != nil {
		return err
	}

	if err := app.Workers.Enqueue(&TmdbUpdateMovie{ID: m.ID.Hex()}); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, gin.H{"error": false, "movie": m})
}

func (a *Application) MoviesShow(c echo.Context, id string) error {
	m := &Movie{}
	err := app.DB.Movie.Find(id, m)
	if err != nil {
		return err
	}

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

	return c.JSON(http.StatusOK, m)
}

func (a *Application) MoviesUpdate(c echo.Context, id string) error {
	data := &Movie{}
	err := c.Bind(&data)
	if err != nil {
		return err
	}

	err = app.DB.MovieUpdate(id, data)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, gin.H{"errors": false, "data": data})
}

func (a *Application) MoviesRefresh(c echo.Context, id string) error {
	if err := app.Workers.Enqueue(&TmdbUpdateMovie{ID: id}); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, gin.H{"error": false})
}

func (a *Application) MoviesSettings(c echo.Context, id string) error {
	data := &Setting{}
	err := c.Bind(&data)
	if err != nil {
		return err
	}

	err = app.DB.MovieSetting(id, data.Setting, data.Value)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, gin.H{"errors": false, "data": data})
}

func (a *Application) MoviesDelete(c echo.Context, id string) error {
	return c.JSON(http.StatusOK, gin.H{"error": false})
}

func (a *Application) MoviesPaths(c echo.Context, id string) error {
	results, err := app.DB.MoviePaths(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, results)
}
