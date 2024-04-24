package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GET /feeds/
func (a *Application) FeedsIndex(c echo.Context, page int, limit int) error {
	list, err := a.DB.FeedList(page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error loading Feeds"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: list})
}

// POST /feeds/
func (a *Application) FeedsCreate(c echo.Context, subject *Feed) error {
	if err := a.DB.Feed.Save(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error saving Feeds"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// GET /feeds/:id
func (a *Application) FeedsShow(c echo.Context, id string) error {
	subject, err := a.DB.FeedGet(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, &Response{Error: true, Message: "not found"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// PUT /feeds/:id
func (a *Application) FeedsUpdate(c echo.Context, id string, subject *Feed) error {
	// if you need to copy or compare to existing object...
	// data, err := a.DB.FeedGet(id)
	// if err != nil {
	//     return c.JSON(http.StatusNotFound, &Response{Error: true, Message: "not found"})
	// }
	// data.Name = subject.Name ...
	if err := a.DB.Feed.Save(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error saving Feeds"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// PATCH /feeds/:id
func (a *Application) FeedsSettings(c echo.Context, id string, setting *Setting) error {
	subject, err := a.DB.FeedGet(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, &Response{Error: true, Message: "not found"})
	}

	// switch Setting.Name {
	// case "something":
	//    subject.Something = Setting.Value
	// }

	if err := a.DB.Feed.Save(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error saving Feeds"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}

// DELETE /feeds/:id
func (a *Application) FeedsDelete(c echo.Context, id string) error {
	subject, err := a.DB.FeedGet(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, &Response{Error: true, Message: "not found"})
	}
	if err := a.DB.Feed.Delete(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error deleting Feeds"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}
