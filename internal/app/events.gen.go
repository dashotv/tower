// Code generated by github.com/dashotv/golem. DO NOT EDIT.
package app

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"

	"github.com/dashotv/fae"
	flame "github.com/dashotv/flame/client"
	"github.com/dashotv/mercury"
	"github.com/dashotv/minion"
	runic "github.com/dashotv/runic/client"
	"github.com/dashotv/tower/internal/plex"
)

func init() {
	initializers = append(initializers, setupEvents)
	healthchecks["events"] = checkEvents
	starters = append(starters, startEvents)
}

type EventsChannel string
type EventsTopic string

func setupEvents(app *Application) error {
	events, err := NewEvents(app)
	if err != nil {
		return err
	}

	app.Events = events
	return nil
}

func startEvents(ctx context.Context, app *Application) error {
	go app.Events.Start()
	return nil
}

func checkEvents(app *Application) error {
	switch app.Events.Merc.Status() {
	case nats.CONNECTED:
		return nil
	default:
		return fae.Errorf("nats status: %s", app.Events.Merc.Status())
	}
}

type Events struct {
	App           *Application
	Merc          *mercury.Mercury
	Log           *zap.SugaredLogger
	Downloading   chan *EventDownloading
	Downloads     chan *EventDownloads
	Episodes      chan *EventEpisodes
	FlameCombined chan *FlameCombined
	Logs          chan *EventLogs
	Movies        chan *EventMovies
	Notices       chan *EventNotices
	PlexSessions  chan *EventPlexSessions
	Releases      chan *Release
	Requests      chan *EventRequests
	RunicReleases chan *runic.Release
	SeerDownloads chan *EventSeerDownload
	SeerEpisodes  chan *EventSeerEpisode
	SeerLogs      chan *EventSeerLog
	SeerNotices   chan *EventSeerNotice
	Series        chan *EventSeries
	Stats         chan *minion.Stats
}

func NewEvents(app *Application) (*Events, error) {
	m, err := mercury.New("tower", app.Config.NatsURL)
	if err != nil {
		return nil, err
	}

	e := &Events{
		App:           app,
		Merc:          m,
		Log:           app.Log.Named("events"),
		Downloading:   make(chan *EventDownloading),
		Downloads:     make(chan *EventDownloads),
		Episodes:      make(chan *EventEpisodes),
		FlameCombined: make(chan *FlameCombined),
		Logs:          make(chan *EventLogs),
		Movies:        make(chan *EventMovies),
		Notices:       make(chan *EventNotices),
		PlexSessions:  make(chan *EventPlexSessions),
		Releases:      make(chan *Release),
		Requests:      make(chan *EventRequests),
		RunicReleases: make(chan *runic.Release),
		SeerDownloads: make(chan *EventSeerDownload),
		SeerEpisodes:  make(chan *EventSeerEpisode),
		SeerLogs:      make(chan *EventSeerLog),
		SeerNotices:   make(chan *EventSeerNotice),
		Series:        make(chan *EventSeries),
		Stats:         make(chan *minion.Stats),
	}

	if err := e.Merc.Sender("tower.downloading", e.Downloading); err != nil {
		return nil, err
	}

	if err := e.Merc.Sender("tower.downloads", e.Downloads); err != nil {
		return nil, err
	}

	if err := e.Merc.Sender("tower.episodes", e.Episodes); err != nil {
		return nil, err
	}

	if err := e.Merc.Receiver("flame.combined", e.FlameCombined); err != nil {
		return nil, err
	}

	if err := e.Merc.Sender("tower.logs", e.Logs); err != nil {
		return nil, err
	}

	if err := e.Merc.Sender("tower.movies", e.Movies); err != nil {
		return nil, err
	}

	if err := e.Merc.Sender("tower.notices", e.Notices); err != nil {
		return nil, err
	}

	if err := e.Merc.Sender("tower.plex_sessions", e.PlexSessions); err != nil {
		return nil, err
	}

	if err := e.Merc.Sender("tower.index.releases", e.Releases); err != nil {
		return nil, err
	}

	if err := e.Merc.Sender("tower.requests", e.Requests); err != nil {
		return nil, err
	}

	if err := e.Merc.Receiver("runic.releases", e.RunicReleases); err != nil {
		return nil, err
	}

	if err := e.Merc.Receiver("seer.downloads", e.SeerDownloads); err != nil {
		return nil, err
	}

	if err := e.Merc.Receiver("seer.episodes", e.SeerEpisodes); err != nil {
		return nil, err
	}

	if err := e.Merc.Receiver("seer.logs", e.SeerLogs); err != nil {
		return nil, err
	}

	if err := e.Merc.Receiver("seer.notices", e.SeerNotices); err != nil {
		return nil, err
	}

	if err := e.Merc.Sender("tower.series", e.Series); err != nil {
		return nil, err
	}

	if err := e.Merc.Sender("tower.stats", e.Stats); err != nil {
		return nil, err
	}
	return e, nil
}

func (e *Events) Start() error {
	e.Log.Debugf("starting events...")
	go func() {
		// wire up receivers
		for {
			select {
			case m := <-e.FlameCombined:
				v, err := onFlameCombined(e.App, m)
				if err != nil {
					e.Log.Errorf("proxy failed: onFlameCombined: %s", err)
					continue
				}
				e.Send("tower.downloading", v)
			case m := <-e.RunicReleases:
				onRunicReleases(e.App, m)

			case m := <-e.SeerDownloads:
				v, err := onSeerDownloads(e.App, m)
				if err != nil {
					e.Log.Errorf("proxy failed: onSeerDownloads: %s", err)
					continue
				}
				e.Send("tower.downloads", v)
			case m := <-e.SeerEpisodes:
				v, err := onSeerEpisodes(e.App, m)
				if err != nil {
					e.Log.Errorf("proxy failed: onSeerEpisodes: %s", err)
					continue
				}
				e.Send("tower.episodes", v)
			case m := <-e.SeerLogs:
				v, err := onSeerLogs(e.App, m)
				if err != nil {
					e.Log.Errorf("proxy failed: onSeerLogs: %s", err)
					continue
				}
				e.Send("tower.logs", v)
			case m := <-e.SeerNotices:
				v, err := onSeerNotices(e.App, m)
				if err != nil {
					e.Log.Errorf("proxy failed: onSeerNotices: %s", err)
					continue
				}
				e.Send("tower.notices", v)
			}
		}
	}()
	return nil
}

func (e *Events) Send(topic EventsTopic, data any) error {
	f := func() interface{} { return e.doSend(topic, data) }

	err, ok := WithTimeout(f, time.Second*5)
	if !ok {
		e.Log.Errorf("timeout sending: %s", topic)
		return fmt.Errorf("timeout sending: %s", topic)
	}
	if err != nil {
		e.Log.Errorf("sending: %s", err)
		return fae.Wrap(err.(error), "events.send")
	}
	return nil
}

func (e *Events) doSend(topic EventsTopic, data any) error {
	switch topic {
	case "tower.downloading":
		m, ok := data.(*EventDownloading)
		if !ok {
			return fae.Errorf("events.send: wrong data type: %t", data)
		}
		e.Downloading <- m

	case "tower.downloads":
		m, ok := data.(*EventDownloads)
		if !ok {
			return fae.Errorf("events.send: wrong data type: %t", data)
		}
		e.Downloads <- m

	case "tower.episodes":
		m, ok := data.(*EventEpisodes)
		if !ok {
			return fae.Errorf("events.send: wrong data type: %t", data)
		}
		e.Episodes <- m

	case "tower.logs":
		m, ok := data.(*EventLogs)
		if !ok {
			return fae.Errorf("events.send: wrong data type: %t", data)
		}
		e.Logs <- m

	case "tower.movies":
		m, ok := data.(*EventMovies)
		if !ok {
			return fae.Errorf("events.send: wrong data type: %t", data)
		}
		e.Movies <- m

	case "tower.notices":
		m, ok := data.(*EventNotices)
		if !ok {
			return fae.Errorf("events.send: wrong data type: %t", data)
		}
		e.Notices <- m

	case "tower.plex_sessions":
		m, ok := data.(*EventPlexSessions)
		if !ok {
			return fae.Errorf("events.send: wrong data type: %t", data)
		}
		e.PlexSessions <- m

	case "tower.index.releases":
		m, ok := data.(*Release)
		if !ok {
			return fae.Errorf("events.send: wrong data type: %t", data)
		}
		e.Releases <- m

	case "tower.requests":
		m, ok := data.(*EventRequests)
		if !ok {
			return fae.Errorf("events.send: wrong data type: %t", data)
		}
		e.Requests <- m

	case "tower.series":
		m, ok := data.(*EventSeries)
		if !ok {
			return fae.Errorf("events.send: wrong data type: %t", data)
		}
		e.Series <- m

	case "tower.stats":
		m, ok := data.(*minion.Stats)
		if !ok {
			return fae.Errorf("events.send: wrong data type: %t", data)
		}
		e.Stats <- m
	default:
		e.Log.Warnf("events.send: unknown topic: %s", topic)
	}
	return nil
}

type EventDownloading struct { // downloading
	Downloads []*Download    `bson:"downloads" json:"downloads"`
	Hashes    map[string]int `bson:"hashes" json:"hashes"`
	Metrics   *flame.Metrics `bson:"metrics" json:"metrics"`
}

type EventDownloads struct { // downloads
	Event    string    `bson:"event" json:"event"`
	ID       string    `bson:"id" json:"id"`
	Download *Download `bson:"download" json:"download"`
}

type EventEpisodes struct { // episodes
	Event   string   `bson:"event" json:"event"`
	ID      string   `bson:"id" json:"id"`
	Episode *Episode `bson:"episode" json:"episode"`
}

type EventLogs struct { // logs
	Event string   `bson:"event" json:"event"`
	ID    string   `bson:"id" json:"id"`
	Log   *Message `bson:"log" json:"log"`
}

type EventMovies struct { // movies
	Event string `bson:"event" json:"event"`
	ID    string `bson:"id" json:"id"`
	Movie *Movie `bson:"movie" json:"movie"`
}

type EventNotices struct { // notices
	Event   string `bson:"event" json:"event"`
	Time    string `bson:"time" json:"time"`
	Class   string `bson:"class" json:"class"`
	Level   string `bson:"level" json:"level"`
	Message string `bson:"message" json:"message"`
}

type EventPlexSessions struct { // plex_sessions
	Sessions []*plex.SessionMetadata `bson:"sessions" json:"sessions"`
}

type EventRequests struct { // requests
	Event   string   `bson:"event" json:"event"`
	ID      string   `bson:"id" json:"id"`
	Request *Request `bson:"request" json:"request"`
}

type EventSeries struct { // series
	Event  string  `bson:"event" json:"event"`
	ID     string  `bson:"id" json:"id"`
	Series *Series `bson:"series" json:"series"`
}
