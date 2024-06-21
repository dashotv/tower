package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GET /file/
func (a *Application) FileIndex(c echo.Context, page int, limit int) error {
	list, total, err := a.DB.FileList(page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error loading File"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: list, Total: total})
}

// POST /file/
func (a *Application) FileCreate(c echo.Context, subject *File) error {
	// TODO: process the subject
	if err := a.DB.File.Save(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error saving File"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// GET /file/:id
func (a *Application) FileShow(c echo.Context, id string) error {
	subject, err := a.DB.FileGet(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "not found"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// PUT /file/:id
func (a *Application) FileUpdate(c echo.Context, id string, subject *File) error {
	// TODO: process the subject

	// if you need to copy or compare to existing object...
	// data, err := a.DB.FileGet(id)
	// if err != nil {
	//     return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "not found"})
	// }
	// data.Name = subject.Name ...
	if err := a.DB.File.Save(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error saving File"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// PATCH /file/:id
func (a *Application) FileSettings(c echo.Context, id string, setting *Setting) error {
	subject, err := a.DB.FileGet(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "not found"})
	}

	// switch Setting.Name {
	// case "something":
	//    subject.Something = Setting.Value
	// }

	if err := a.DB.File.Save(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error saving File"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// DELETE /file/:id
func (a *Application) FileDelete(c echo.Context, id string) error {
	subject, err := a.DB.FileGet(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "not found"})
	}
	if err := a.DB.File.Delete(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error deleting File"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}
