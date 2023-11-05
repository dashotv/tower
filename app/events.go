package app

import (
	"time"

	"github.com/dashotv/mercury"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var events *Events

type EventsChannel string
type EventsTopic string

type Events struct {
	Merc      *mercury.Mercury
	Log       *zap.SugaredLogger
	Receivers map[EventsTopic]chan any
	Senders   map[EventsTopic]chan any
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

type EventTowerEpisode struct {
	Event   string
	ID      string
	Episode *Episode
}
type EventTowerSeries struct {
	Event  string
	ID     string
	Series *Series
}
type EventTowerMovie struct {
	Event string
	ID    string
	Movie *Movie
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
	e.Senders = map[EventsTopic]chan any{
		"tower.episodes": make(chan any),
		"tower.series":   make(chan any),
		"tower.movies":   make(chan any),
	}
	e.Receivers = map[EventsTopic]chan any{
		"seer.notices":   make(chan any, 5),
		"seer.downloads": make(chan any, 5),
	}

	for topic, channel := range e.Senders {
		if err := e.Merc.Sender(string(topic), channel); err != nil {
			return nil, err
		}
	}
	for topic, channel := range e.Receivers {
		if err := e.Merc.Receiver(string(topic), channel); err != nil {
			return nil, err
		}
	}

	return e, nil
}

func (e *Events) Start() error {
	e.Log.Infof("starting events...")

	for {
		select {
		case r := <-e.Receivers["seer.notices"]:
			m := r.(EventSeerNotice)
			if m.Message == "processing downloads" {
				cache.Set("seer_downloads", time.Now().Unix())
			}
		case r := <-e.Receivers["seer.downloads"]:
			m := r.(EventSeerDownload)
			e.Log.Infof("download: %s %s", m.ID, m.Event)
		}
	}
}

func (e *Events) Send(topic EventsTopic, data any) error {
	c := events.Senders[topic]
	if c == nil {
		return errors.Errorf("events: %s does not exist", topic)
	}
	e.Log.Infof("sending %s: %+v", topic, data)
	c <- data
	return nil
}

func setupEvents() (err error) {
	events, err = NewEvents()
	if err != nil {
		return err
	}
	return nil
}
