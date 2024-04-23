package app

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// POST /paths/:id
func (a *Application) PathsUpdate(c echo.Context, id string, medium_id string, path *Path) error {
	m, err := a.DB.Medium.Get(medium_id, &Medium{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "not found"})
	}

	list := lo.Filter(m.Paths, func(p *Path, i int) bool {
		return p.ID.Hex() == id
	})
	if len(list) == 0 {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "not found"})
	}
	if len(list) > 1 {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "duplicate Paths"})
	}

	t := fileType(fmt.Sprintf("%s.%s", path.Local, path.Extension))
	if t == "" {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "unknown file type"})
	}

	subject := list[0]
	subject.Local = path.Local
	subject.Extension = path.Extension
	subject.Type = primitive.Symbol(t)

	if err := a.DB.Medium.Save(m); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error saving Paths"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// DELETE /paths/:id
func (a *Application) PathsDelete(c echo.Context, id string, medium_id string) error {
	if err := a.Workers.Enqueue(&PathDelete{MediumID: medium_id, PathID: id}); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error deleting Paths"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false})
}
