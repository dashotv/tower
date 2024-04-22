package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GET /library/
func (a *Application) LibraryIndex(c echo.Context, page int, limit int) error {
	list, err := a.DB.LibraryList(page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error loading Library"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: list})
}

// POST /library/
func (a *Application) LibraryCreate(c echo.Context, subject *Library) error {
	// TODO: process the subject
	if err := a.DB.Library.Save(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error saving Library"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// GET /library/:id
func (a *Application) LibraryShow(c echo.Context, id string) error {
	subject, err := a.DB.LibraryGet(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "not found"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// PUT /library/:id
func (a *Application) LibraryUpdate(c echo.Context, id string, subject *Library) error {
	// TODO: process the subject

	// if you need to copy or compare to existing object...
	// data, err := a.DB.LibraryGet(id)
	// if err != nil {
	//     return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "not found"})
	// }
	// data.Name = subject.Name ...
	if err := a.DB.Library.Save(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error saving Library"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// PATCH /library/:id
func (a *Application) LibrarySettings(c echo.Context, id string, setting *Setting) error {
	subject, err := a.DB.LibraryGet(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "not found"})
	}

	// switch Setting.Name {
	// case "something":
	//    subject.Something = Setting.Value
	// }

	if err := a.DB.Library.Save(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error saving Library"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// DELETE /library/:id
func (a *Application) LibraryDelete(c echo.Context, id string) error {
	subject, err := a.DB.LibraryGet(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "not found"})
	}
	if err := a.DB.Library.Delete(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error deleting Library"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}
