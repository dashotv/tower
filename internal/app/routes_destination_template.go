package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GET /destination_template/
func (a *Application) DestinationTemplateIndex(c echo.Context, page int, limit int) error {
	list, err := a.DB.DestinationTemplateList(page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error loading DestinationTemplate"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: list})
}

// POST /destination_template/
func (a *Application) DestinationTemplateCreate(c echo.Context, subject *DestinationTemplate) error {
	// TODO: process the subject
	if err := a.DB.DestinationTemplate.Save(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error saving DestinationTemplate"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// GET /destination_template/:id
func (a *Application) DestinationTemplateShow(c echo.Context, id string) error {
	subject, err := a.DB.DestinationTemplateGet(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "not found"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// PUT /destination_template/:id
func (a *Application) DestinationTemplateUpdate(c echo.Context, id string, subject *DestinationTemplate) error {
	// TODO: process the subject

	// if you need to copy or compare to existing object...
	// data, err := a.DB.DestinationTemplateGet(id)
	// if err != nil {
	//     return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "not found"})
	// }
	// data.Name = subject.Name ...
	if err := a.DB.DestinationTemplate.Save(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error saving DestinationTemplate"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// PATCH /destination_template/:id
func (a *Application) DestinationTemplateSettings(c echo.Context, id string, setting *Setting) error {
	subject, err := a.DB.DestinationTemplateGet(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "not found"})
	}

	// switch Setting.Name {
	// case "something":
	//    subject.Something = Setting.Value
	// }

	if err := a.DB.DestinationTemplate.Save(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error saving DestinationTemplate"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// DELETE /destination_template/:id
func (a *Application) DestinationTemplateDelete(c echo.Context, id string) error {
	subject, err := a.DB.DestinationTemplateGet(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "not found"})
	}
	if err := a.DB.DestinationTemplate.Delete(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error deleting DestinationTemplate"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}
