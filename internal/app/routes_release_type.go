package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GET /release_type/
func (a *Application) ReleaseTypeIndex(c echo.Context, page int, limit int) error {
	list, err := a.DB.ReleaseTypeList(page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error loading ReleaseType"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: list})
}

// POST /release_type/
func (a *Application) ReleaseTypeCreate(c echo.Context, subject *ReleaseType) error {
	// TODO: process the subject
	if err := a.DB.ReleaseType.Save(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error saving ReleaseType"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// GET /release_type/:id
func (a *Application) ReleaseTypeShow(c echo.Context, id string) error {
	subject, err := a.DB.ReleaseTypeGet(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "not found"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// PUT /release_type/:id
func (a *Application) ReleaseTypeUpdate(c echo.Context, id string, subject *ReleaseType) error {
	// TODO: process the subject

	// if you need to copy or compare to existing object...
	// data, err := a.DB.ReleaseTypeGet(id)
	// if err != nil {
	//     return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "not found"})
	// }
	// data.Name = subject.Name ...
	if err := a.DB.ReleaseType.Save(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error saving ReleaseType"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// PATCH /release_type/:id
func (a *Application) ReleaseTypeSettings(c echo.Context, id string, setting *Setting) error {
	subject, err := a.DB.ReleaseTypeGet(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "not found"})
	}

	// switch Setting.Name {
	// case "something":
	//    subject.Something = Setting.Value
	// }

	if err := a.DB.ReleaseType.Save(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error saving ReleaseType"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// DELETE /release_type/:id
func (a *Application) ReleaseTypeDelete(c echo.Context, id string) error {
	subject, err := a.DB.ReleaseTypeGet(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "not found"})
	}
	if err := a.DB.ReleaseType.Delete(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error deleting ReleaseType"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}
