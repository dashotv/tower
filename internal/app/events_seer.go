package app

import (
	"time"

	"github.com/dashotv/fae"
)

type EventSeerNotice struct {
	Event   string `json:"event,omitempty"`
	Time    string `json:"time,omitempty"`
	Class   string `json:"class,omitempty"`
	Level   string `json:"level,omitempty"`
	Message string `json:"message,omitempty"`
}

type EventSeerDownload struct {
	Event    string    `json:"event,omitempty"`
	ID       string    `json:"id,omitempty"`
	Download *Download `json:"download,omitempty"`
}

type EventSeerEpisode struct {
	Event string `json:"event,omitempty"`
	ID    string `json:"id,omitempty"`
}

type EventSeerLog struct {
	Time     time.Time `json:"time,omitempty"`
	Message  string    `json:"message,omitempty"`
	Level    string    `json:"level,omitempty"`
	Facility string    `json:"facility,omitempty"`
}

func onSeerDownloads(app *Application, msg *EventSeerDownload) (*EventDownloads, error) {
	d, err := app.DB.DownloadGet(msg.ID)
	if err != nil {
		return nil, fae.Wrap(err, "loading download")
	}
	return &EventDownloads{Event: msg.Event, ID: d.ID.Hex(), Download: d}, nil
}

func onSeerEpisodes(app *Application, msg *EventSeerEpisode) (*EventEpisodes, error) {
	ep, err := app.DB.EpisodeGet(msg.ID)
	if err != nil {
		return nil, fae.Wrap(err, "loading episode")
	}
	return &EventEpisodes{Event: msg.Event, ID: ep.ID.Hex(), Episode: ep}, nil
}

func onSeerLogs(app *Application, msg *EventSeerLog) (*EventLogs, error) {
	level := msg.Level
	if msg.Level == "info" {
		level = "debug"
	}

	l, err := app.DB.MessageCreate(level, msg.Message, msg.Facility, msg.Time)
	if err != nil {
		app.Events.Log.Errorf("error saving log: %s", err)
	}
	return &EventLogs{Event: "new", ID: l.ID.Hex(), Log: l}, nil
}

func onSeerNotices(app *Application, msg *EventSeerNotice) (*EventNotices, error) {
	if msg.Message == "processing downloads" {
		app.Cache.Set("seer_downloads", time.Now().Unix())
	}
	level := msg.Level
	if msg.Level == "info" {
		level = "debug"
	}
	l, err := app.DB.MessageCreate(level, msg.Message, msg.Class, time.Now())
	if err != nil {
		app.Events.Log.Errorf("error saving log: %s", err)
	}
	n := &EventNotices{
		Event:   msg.Event,
		Time:    msg.Time,
		Class:   msg.Class,
		Level:   msg.Level,
		Message: msg.Message,
	}
	if err := app.Events.Send("tower.logs", &EventLogs{Event: "new", ID: l.ID.Hex(), Log: l}); err != nil {
		app.Events.Log.Errorf("error sending log: %s", err)
	}
	return n, nil
}
