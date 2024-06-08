package app

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"

	"github.com/dashotv/fae"
	"github.com/dashotv/tower/internal/plex"
)

func plexPinToPin(pin *plex.Pin) *Pin {
	return &Pin{
		Pin:        pin.ID,
		Code:       pin.Code,
		Product:    pin.Product,
		Identifier: pin.Identifier,
		Token:      pin.Token,
	}
}

func pinToPlexPin(pin *Pin) *plex.Pin {
	return &plex.Pin{
		ID:         pin.Pin,
		Code:       pin.Code,
		Product:    pin.Product,
		Identifier: pin.Identifier,
		Token:      pin.Token,
	}
}

func (a *Application) PlexIndex(c echo.Context) error {
	// get pin
	plexPin, err := app.Plex.CreatePin()
	if err != nil {
		return err
	}

	pin := plexPinToPin(plexPin)

	app.Log.Debugf("PlexIndex: saving pin %+v", pin)
	err = app.DB.Pin.Save(pin)
	if err != nil {
		return err
	}

	authUrl := app.Plex.GetAuthUrl(app.Config.Plex, plexPin)
	return c.Redirect(302, authUrl)
}

func (a *Application) PlexAuth(c echo.Context) error {
	id := c.QueryParam("pin")
	pinID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	list, err := app.DB.Pin.Query().Where("pin", int64(pinID)).Run()
	if err != nil {
		return err
	}
	if len(list) != 1 {
		return fae.New("pin not found")
	}

	plexPin := pinToPlexPin(list[0])
	ok, err := app.Plex.CheckPin(plexPin)
	if err != nil {
		return err
	}
	if !ok {
		return fae.New("something went wrong...")
	}

	list[0].Token = plexPin.Token
	if err := app.DB.Pin.Save(list[0]); err != nil {
		return err
	}

	if err := app.Workers.Enqueue(&PlexPinToUsers{}); err != nil {
		return err
	}
	return c.String(http.StatusOK, "Authorization complete!")
}

func (a *Application) PlexUpdate(c echo.Context) error {
	if err := app.Workers.Enqueue(&PlexPinToUsers{}); err != nil {
		return err
	}
	return c.String(http.StatusOK, "Updating users...")
}

func (a *Application) PlexLibraries(c echo.Context) error {
	list, err := a.Plex.GetLibraries()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: list})
}

func (a *Application) PlexSearch(c echo.Context, query, section string) error {
	list, err := a.Plex.Search(query, section)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &Response{Error: false, Result: list})
}

func (a *Application) PlexCollectionsIndex(c echo.Context, section string) error {
	list, err := a.Plex.ListCollections(section)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: list})
}
func (a *Application) PlexCollectionsShow(c echo.Context, section, ratingKey string) error {
	list, err := a.Plex.GetCollection(ratingKey)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: list})
}

func (a *Application) PlexMetadata(c echo.Context, key string) error {
	a.Log.Debugf("PlexMetadata: key=%s", key)
	resp, err := a.Plex.GetMetadataByKey(key)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: resp})
}

func (a *Application) PlexClients(c echo.Context) error {
	list, err := a.Plex.GetClients()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: list})
}

func (a *Application) PlexDevices(c echo.Context) error {
	list, err := a.Plex.GetDevices()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: list})
}

func (a *Application) PlexResources(c echo.Context) error {
	provides := QueryString(c, "provides")
	if provides == "" {
		provides = "player"
	}

	list, err := a.Plex.GetResources()
	if err != nil {
		return err
	}

	// filter by provides
	filtered := lo.Filter(list, func(r *plex.Resource, i int) bool {
		provided := strings.Split(r.Provides, ",")
		return lo.Contains(provided, provides) && r.Name != "iPhone"
	})

	return c.JSON(http.StatusOK, &Response{Error: false, Result: filtered})
}
func (a *Application) PlexPlay(c echo.Context, ratingKey, player string) error {
	err := a.Plex.Play(ratingKey, player)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &Response{Error: false})
}
func (a *Application) PlexStop(c echo.Context, session string) error {
	err := a.Plex.Stop(session)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &Response{Error: false})
}
func (a *Application) PlexSessions(c echo.Context) error {
	list, err := a.Plex.GetSessions()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: list})
}
func (a *Application) PlexFiles(c echo.Context) error {
	return c.JSON(http.StatusOK, &Response{Error: false, Result: a.PlexFileCache.files})
}
