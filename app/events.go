package app

import (
	"time"

	"github.com/dashotv/mercury"
	"go.uber.org/zap"
)

var events *Events

type EventsChannel string
type EventsTopic string

type Events struct {
	Merc *mercury.Mercury
	Log  *zap.SugaredLogger
}

type EventSeerNotice struct {
	Event   string
	Time    string
	Class   string
	Level   string
	Message string
}

type EventSeerDownload struct {
	Event    string
	ID       string
	Download *Download
}

func NewEvents() (*Events, error) {
	m, err := mercury.New("tower", cfg.Nats.URL)
	if err != nil {
		return nil, err
	}

	e := &Events{
		Merc: m,
		Log:  log.Named("events"),
	}

	return e, nil
}

func (e *Events) Start() error {
	seer_notices := make(chan *EventSeerNotice, 5)
	if err := e.Merc.Receiver("seer.notices", seer_notices); err != nil {
		return err
	}
	seer_downloads := make(chan *EventSeerDownload, 5)
	if err := e.Merc.Receiver("seer.downloads", seer_downloads); err != nil {
		return err
	}

	for {
		select {
		case r := <-seer_notices:
			if r.Message == "processing downloads" {
				cache.Set("seer_downloads", time.Now().Unix())
			}
		case r := <-seer_downloads:
			e.Log.Infof("download: %s %s", r.ID, r.Event)
		}
	}
}

func setupEvents() (err error) {
	events, err = NewEvents()
	if err != nil {
		return err
	}
	return nil
}
