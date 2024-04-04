package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GET /collections/
func (a *Application) CollectionsIndex(c echo.Context, page int, limit int) error {
	list, err := a.DB.CollectionList(page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error loading Collections"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: list})
}

// POST /collections/
func (a *Application) CollectionsCreate(c echo.Context, subject *Collection) error {
	// TODO: process the subject
	if err := a.DB.Collection.Save(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error saving Collections"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// GET /collections/:id
func (a *Application) CollectionsShow(c echo.Context, id string) error {
	subject, err := a.DB.CollectionGet(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, &Response{Error: true, Message: "not found"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// PUT /collections/:id
func (a *Application) CollectionsUpdate(c echo.Context, id string, subject *Collection) error {
	// TODO: process the subject

	// if you need to copy or compare to existing object...
	// data, err := a.DB.CollectionGet(id)
	// if err != nil {
	//     return c.JSON(http.StatusNotFound, &Response{Error: true, Message: "not found"})
	// }
	// data.Name = subject.Name ...
	if err := a.DB.Collection.Save(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error saving Collections"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// PATCH /collections/:id
func (a *Application) CollectionsSettings(c echo.Context, id string, setting *Setting) error {
	subject, err := a.DB.CollectionGet(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, &Response{Error: true, Message: "not found"})
	}

	// switch Setting.Name {
	// case "something":
	//    subject.Something = Setting.Value
	// }

	if err := a.DB.Collection.Save(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error saving Collections"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// DELETE /collections/:id
func (a *Application) CollectionsDelete(c echo.Context, id string) error {
	subject, err := a.DB.CollectionGet(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, &Response{Error: true, Message: "not found"})
	}
	if err := a.DB.Collection.Delete(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error deleting Collections"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}
