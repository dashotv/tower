package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
)

func (a *Application) CollectionsIndex(c echo.Context, page int, limit int) error {
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 25
	}

	list, err := a.DB.CollectionList(limit, (page-1)*limit)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, gin.H{"error": false, "count": len(list), "results": list})
}

func (a *Application) CollectionsCreate(c echo.Context) error {
	col := &Collection{}
	err := c.Bind(col)
	if err != nil {
		return err
	}

	err = a.DB.Collection.Save(col)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, gin.H{"error": false, "id": col.ID, "collection": col})
}

func (a *Application) CollectionsShow(c echo.Context, id string) error {
	subject, err := a.DB.CollectionGet(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, subject)
}

func (a *Application) CollectionsUpdate(c echo.Context, id string) error {
	subject := &Collection{}

	if err := c.Bind(subject); err != nil {
		return err
	}

	if err := a.DB.Collection.Save(subject); err != nil {
		return err
	}

	if len(subject.Media) > 0 {
		if err := a.Workers.Enqueue(&PlexCollectionUpdate{Id: subject.ID.Hex()}); err != nil {
			return err
		}
	}

	return nil
}

func (a *Application) CollectionsSettings(c echo.Context, id string) error {
	// asssuming this is a CRUD route, get the subject from the database
	// subject, err := a.DB.Collections.Get(id)
	return c.JSON(http.StatusNotImplemented, H{"message": "not implemented"})
}

func (a *Application) CollectionsDelete(c echo.Context, id string) error {
	col, err := a.DB.Collection.Get(id, &Collection{})
	if err != nil {
		return err
	}

	if col.RatingKey != "" {
		if err := app.Plex.DeleteCollection(col.RatingKey); err != nil {
			return err
		}
	}

	err = a.DB.Collection.Delete(col)
	if err != nil {
		return err
	}

	return nil
}
