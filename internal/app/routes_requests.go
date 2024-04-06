package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GET /requests/
func (a *Application) RequestsIndex(c echo.Context, page int, limit int) error {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 25
	}
	skip := (page - 1) * limit
	list, err := a.DB.Request.Query().Desc("created_at").Limit(limit).Skip(skip).Run()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error loading Requests"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: list})
}

// POST /requests/
func (a *Application) RequestsCreate(c echo.Context, subject *Request) error {
	// TODO: process the subject
	if err := a.DB.Request.Save(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error saving Requests"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// GET /requests/:id
func (a *Application) RequestsShow(c echo.Context, id string) error {
	subject, err := a.DB.Request.Get(id, &Request{})
	if err != nil {
		return c.JSON(http.StatusNotFound, &Response{Error: true, Message: "not found"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// PUT /requests/:id
func (a *Application) RequestsUpdate(c echo.Context, id string, subject *Request) error {
	req := &Request{}
	err := app.DB.Request.Find(id, req)
	if err != nil {
		return err
	}

	req.Status = subject.Status
	if err := app.DB.Request.Update(req); err != nil {
		return err
	}

	if subject.Status == "approved" {
		if err := app.Workers.Enqueue(&CreateMediaFromRequests{}); err != nil {
			return err
		}
	}

	if err := a.DB.Request.Save(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error saving Requests"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// PATCH /requests/:id
func (a *Application) RequestsSettings(c echo.Context, id string, setting *Setting) error {
	subject, err := a.DB.Request.Get(id, &Request{})
	if err != nil {
		return c.JSON(http.StatusNotFound, &Response{Error: true, Message: "not found"})
	}

	// switch Setting.Name {
	// case "something":
	//    subject.Something = Setting.Value
	// }

	if err := a.DB.Request.Save(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error saving Requests"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// DELETE /requests/:id
func (a *Application) RequestsDelete(c echo.Context, id string) error {
	subject, err := a.DB.Request.Get(id, &Request{})
	if err != nil {
		return c.JSON(http.StatusNotFound, &Response{Error: true, Message: "not found"})
	}
	if err := a.DB.Request.Delete(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error deleting Requests"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}
