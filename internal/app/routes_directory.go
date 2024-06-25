package app

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/dashotv/fae"
)

// GET /directory/index
func (a *Application) DirectoryIndex(c echo.Context, library string, page, limit int) error {
	// parts := strings.Split(path, string(filepath.Separator))
	// if path != "" && len(parts) == 2 {
	// 	list, total, err := a.DB.DirectoryFiles(parts[0], parts[1], page, limit)
	// 	if err != nil {
	// 		return fae.Wrap(err, "directory files")
	// 	}
	// 	return c.JSON(http.StatusOK, &Response{Error: false, Result: list, Total: total})
	// } else
	if library != "" { // && len(parts) == 1 {
		list, total, err := a.DB.DirectoryMedia(library, page, limit)
		if err != nil {
			return fae.Wrap(err, "directory media")
		}
		return c.JSON(http.StatusOK, &Response{Error: false, Result: list, Total: total})
	} else {
		list, total, err := a.DB.DirectoryLibraries(page, limit)
		if err != nil {
			return fae.Wrap(err, "directory libraries")
		}
		return c.JSON(http.StatusOK, &Response{Error: false, Result: list, Total: total})
	}
}
