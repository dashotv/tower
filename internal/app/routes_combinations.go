package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GET /combinations/
func (a *Application) CombinationsIndex(c echo.Context, page int, limit int) error {
	list, err := a.DB.CombinationList(page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: list})
}

// POST /combinations/
func (a *Application) CombinationsCreate(c echo.Context, subject *Combination) error {
	if err := a.DB.Combination.Save(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// GET /combinations/:id
func (a *Application) CombinationsShow(c echo.Context, name string) error {
	children, err := a.DB.CombinationChildren(name)
	if err != nil {
		return c.JSON(http.StatusNotFound, &Response{Error: true, Message: "not found"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: children})
}

// PUT /combinations/:id
func (a *Application) CombinationsUpdate(c echo.Context, id string, subject *Combination) error {
	if err := a.DB.Combination.Update(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// DELETE /combinations/:id
func (a *Application) CombinationsDelete(c echo.Context, id string) error {
	subject, err := a.DB.CombinationGet(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "getting subject failed:" + err.Error()})
	}
	if err := a.DB.Combination.Delete(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}
