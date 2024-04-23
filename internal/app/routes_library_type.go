package app

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GET /release_type/
func (a *Application) LibraryTypeIndex(c echo.Context, page int, limit int) error {
	list, err := a.DB.LibraryTypeList(page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error loading LibraryType"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: list})
}

// POST /release_type/
func (a *Application) LibraryTypeCreate(c echo.Context, subject *LibraryType) error {
	// TODO: process the subject
	if err := a.DB.LibraryType.Save(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error saving LibraryType"})
	}
	if err := startDestination(context.Background(), a); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error starting destination"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// GET /release_type/:id
func (a *Application) LibraryTypeShow(c echo.Context, id string) error {
	subject, err := a.DB.LibraryTypeGet(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "not found"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// PUT /release_type/:id
func (a *Application) LibraryTypeUpdate(c echo.Context, id string, subject *LibraryType) error {
	// TODO: process the subject

	// if you need to copy or compare to existing object...
	// data, err := a.DB.LibraryTypeGet(id)
	// if err != nil {
	//     return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "not found"})
	// }
	// data.Name = subject.Name ...
	if err := a.DB.LibraryType.Save(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error saving LibraryType"})
	}
	if err := startDestination(context.Background(), a); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error starting destination"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// PATCH /release_type/:id
func (a *Application) LibraryTypeSettings(c echo.Context, id string, setting *Setting) error {
	subject, err := a.DB.LibraryTypeGet(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "not found"})
	}

	// switch Setting.Name {
	// case "something":
	//    subject.Something = Setting.Value
	// }

	if err := a.DB.LibraryType.Save(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error saving LibraryType"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// DELETE /release_type/:id
func (a *Application) LibraryTypeDelete(c echo.Context, id string) error {
	subject, err := a.DB.LibraryTypeGet(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "not found"})
	}
	if err := a.DB.LibraryType.Delete(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error deleting LibraryType"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}
