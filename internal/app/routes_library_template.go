package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GET /library_template/
func (a *Application) LibraryTemplateIndex(c echo.Context, page int, limit int) error {
	list, err := a.DB.LibraryTemplateList(page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error loading LibraryTemplate"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: list})
}

// POST /library_template/
func (a *Application) LibraryTemplateCreate(c echo.Context, subject *LibraryTemplate) error {
	// TODO: process the subject
	if err := a.DB.LibraryTemplate.Save(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error saving LibraryTemplate"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// GET /library_template/:id
func (a *Application) LibraryTemplateShow(c echo.Context, id string) error {
	subject, err := a.DB.LibraryTemplateGet(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "not found"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// PUT /library_template/:id
func (a *Application) LibraryTemplateUpdate(c echo.Context, id string, subject *LibraryTemplate) error {
	// TODO: process the subject

	// if you need to copy or compare to existing object...
	// data, err := a.DB.LibraryTemplateGet(id)
	// if err != nil {
	//     return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "not found"})
	// }
	// data.Name = subject.Name ...
	if err := a.DB.LibraryTemplate.Save(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error saving LibraryTemplate"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// PATCH /library_template/:id
func (a *Application) LibraryTemplateSettings(c echo.Context, id string, setting *Setting) error {
	subject, err := a.DB.LibraryTemplateGet(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "not found"})
	}

	// switch Setting.Name {
	// case "something":
	//    subject.Something = Setting.Value
	// }

	if err := a.DB.LibraryTemplate.Save(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error saving LibraryTemplate"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// DELETE /library_template/:id
func (a *Application) LibraryTemplateDelete(c echo.Context, id string) error {
	subject, err := a.DB.LibraryTemplateGet(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "not found"})
	}
	if err := a.DB.LibraryTemplate.Delete(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error deleting LibraryTemplate"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}
