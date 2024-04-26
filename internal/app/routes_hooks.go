package app

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"

	"github.com/dashotv/fae"
	"github.com/dashotv/tower/internal/plex"
)

// GET /hooks/plex
func (a *Application) HooksPlex(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: err.Error()})
	}

	if form.Value == nil {
		// return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "missing payload"})
		return fae.New("missing value")
	}

	data, ok := form.Value["payload"]
	if !ok {
		return fae.New("missing payload")
	}

	payload := &plex.WebhookPayload{}
	if err := json.Unmarshal([]byte(data[0]), payload); err != nil {
		return fae.Wrap(err, "unmarshal payload")
	}

	notifier.Log.Debugf("plex", "webhook received: %s", payload.Event)
	if lo.Contains(plexSupportedHooks, payload.Event) {
		if err := a.Workers.Enqueue(&PlexWebhook{Payload: payload}); err != nil {
			return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: err.Error()})
		}
	}

	return c.JSON(http.StatusOK, &Response{Error: false, Message: "ok"})
}

// POST /hooks/nzbget
func (a *Application) HooksNzbget(c echo.Context, payload *NzbgetPayload) error {
	if err := a.Workers.Enqueue(&NzbgetProcess{Payload: payload}); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, &Response{Error: false})
}
