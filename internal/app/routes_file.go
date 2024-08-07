package app

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/dashotv/fae"
)

// GET /file/
func (a *Application) FileIndex(c echo.Context, page int, limit int) error {
	list, total, err := a.DB.FileList(page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error loading File"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: list, Total: total})
}

// GET /file/missing
func (a *Application) FileMissing(c echo.Context, page int, limit int, medium_id string) error {
	list, total, err := a.DB.FileMissing(page, limit)
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
	if err := a.DB.File.Save(subject); err != nil {
		return fae.Wrap(err, "error saving File")
	}
	if err := a.Workers.Enqueue(&FilesMove{ID: subject.ID.Hex(), Title: subject.Path}); err != nil {
		return fae.Wrap(err, "error enqueuing move")
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

// GET /file/list
func (a *Application) FileList(c echo.Context, page int, limit int, medium_id string) error {
	if medium_id != "" {
		list, total, err := a.DB.DirectoryFiles(medium_id, page, limit)
		if err != nil {
			return fae.Wrap(err, "directory files")
		}
		return c.JSON(http.StatusOK, &Response{Error: false, Result: list, Total: total})
	}
	list, total, err := a.DB.FileList(page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error loading File"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: list, Total: total})
}
