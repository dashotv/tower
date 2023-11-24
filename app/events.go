package app

import (
	"fmt"
	"time"

	"github.com/dashotv/mercury"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var events *Events

type EventsChannel string
type EventsTopic string

type Events struct {
	Merc           *mercury.Mercury
	Log            *zap.SugaredLogger
	SeerLogs       chan *EventSeerLog
	SeerDownloads  chan *EventSeerDownload
	SeerNotices    chan *EventSeerNotice
	SeerEpisodes   chan *EventSeerEpisode
	TowerLogs      chan *EventTowerLog
	TowerEpisodes  chan *EventTowerEpisode
	TowerSeries    chan *EventTowerSeries
	TowerMovies    chan *EventTowerMovie
	TowerEvents    chan *EventTowerRequest
	TowerDownloads chan *EventTowerDownload
}

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

type EventTowerDownload struct {
	Event    string    `json:"event,omitempty"`
	ID       string    `json:"id,omitempty"`
	Download *Download `json:"download,omitempty"`
}
type EventTowerEpisode struct {
	Event   string   `json:"event,omitempty"`
	ID      string   `json:"id,omitempty"`
	Episode *Episode `json:"episode,omitempty"`
}
type EventTowerSeries struct {
	Event  string  `json:"event,omitempty"`
	ID     string  `json:"id,omitempty"`
	Series *Series `json:"series,omitempty"`
}
type EventTowerMovie struct {
	Event string `json:"event,omitempty"`
	ID    string `json:"id,omitempty"`
	Movie *Movie `json:"movie,omitempty"`
}
type EventTowerLog struct {
	Event string   `json:"event,omitempty"`
	ID    string   `json:"id,omitempty"`
	Log   *Message `json:"log,omitempty"`
}
type EventTowerRequest struct {
	Event   string   `json:"event,omitempty"`
	ID      string   `json:"id,omitempty"`
	Request *Request `json:"request,omitempty"`
}

func NewEvents() (*Events, error) {
	m, err := mercury.New("tower", cfg.Nats.URL)
	if err != nil {
		return nil, err
	}

	e := &Events{
		Merc:           m,
		Log:            log.Named("events"),
		SeerLogs:       make(chan *EventSeerLog, 5),
		SeerDownloads:  make(chan *EventSeerDownload, 5),
		SeerNotices:    make(chan *EventSeerNotice, 5),
		SeerEpisodes:   make(chan *EventSeerEpisode, 5),
		TowerLogs:      make(chan *EventTowerLog),
		TowerEpisodes:  make(chan *EventTowerEpisode),
		TowerSeries:    make(chan *EventTowerSeries),
		TowerMovies:    make(chan *EventTowerMovie),
		TowerEvents:    make(chan *EventTowerRequest),
		TowerDownloads: make(chan *EventTowerDownload),
	}

	if err := e.Merc.Receiver("seer.logs", e.SeerLogs); err != nil {
		return nil, err
	}
	if err := e.Merc.Receiver("seer.downloads", e.SeerDownloads); err != nil {
		return nil, err
	}
	if err := e.Merc.Receiver("seer.notices", e.SeerNotices); err != nil {
		return nil, err
	}
	if err := e.Merc.Receiver("seer.episodes", e.SeerEpisodes); err != nil {
		return nil, err
	}
	if err := e.Merc.Sender("tower.logs", e.TowerLogs); err != nil {
		return nil, err
	}
	if err := e.Merc.Sender("tower.episodes", e.TowerEpisodes); err != nil {
		return nil, err
	}
	if err := e.Merc.Sender("tower.series", e.TowerSeries); err != nil {
		return nil, err
	}
	if err := e.Merc.Sender("tower.movies", e.TowerMovies); err != nil {
		return nil, err
	}
	if err := e.Merc.Sender("tower.requests", e.TowerEvents); err != nil {
		return nil, err
	}
	if err := e.Merc.Sender("tower.downloads", e.TowerDownloads); err != nil {
		return nil, err
	}

	return e, nil
}

func (e *Events) Start() error {
	e.Log.Infof("starting events...")

	for {
		select {
		case m := <-e.SeerNotices:
			if m.Message == "processing downloads" {
				cache.Set("seer_downloads", time.Now().Unix())
			}
			l := &Message{
				Level:    m.Level,
				Message:  m.Message,
				Facility: m.Class,
			}
			if err := db.Message.Save(l); err != nil {
				e.Log.Errorf("error saving log: %s", err)
			}
			e.Send("tower.logs", &EventTowerLog{Event: "new", ID: l.ID.Hex(), Log: l})
		case m := <-e.SeerLogs:
			l := &Message{
				Level:    m.Level,
				Message:  m.Message,
				Facility: m.Facility,
			}
			l.CreatedAt = m.Time
			if err := db.Message.Save(l); err != nil {
				e.Log.Errorf("error saving log: %s", err)
			}
			e.Send("tower.logs", &EventTowerLog{Event: "new", ID: l.ID.Hex(), Log: l})
		case m := <-e.SeerDownloads:
			d := &Download{}
			err := db.Download.Find(m.ID, d)
			if err != nil {
				e.Log.Errorf("error finding download: %s", err)
				continue
			}
			db.processDownloads([]*Download{d})
			e.Send("tower.downloads", &EventTowerDownload{Event: m.Event, ID: d.ID.Hex(), Download: d})
		case m := <-e.SeerEpisodes:
			ep := &Episode{}
			err := db.Episode.Find(m.ID, ep)
			if err != nil {
				e.Log.Errorf("error finding episode: %s", err)
				continue
			}
			db.processEpisode(ep)
			e.Send("tower.episodes", &EventTowerEpisode{Event: m.Event, ID: ep.ID.Hex(), Episode: ep})
		}
	}
}

func (e *Events) Send(topic EventsTopic, data any) error {
	f := func() interface{} { return e.doSend(topic, data) }

	err, ok := WithTimeout(f, time.Second*2)
	if !ok {
		e.Log.Errorf("events.send: timeout sending message: %s", topic)
		return fmt.Errorf("events.send: timeout sending message: %s", topic)
	}
	if err != nil {
		e.Log.Errorf("events.send: %s", err)
		return errors.Wrap(err.(error), "events.send")
	}
	return nil
}

func (e *Events) doSend(topic EventsTopic, data any) error {
	switch topic {
	case "tower.logs":
		m, ok := data.(*EventTowerLog)
		if !ok {
			e.Log.Errorf("events.send: wrong data type: %t", data)
			return errors.New("events.send: wrong data type")
		}

		e.TowerLogs <- m
	case "tower.episodes":
		m, ok := data.(*EventTowerEpisode)
		if !ok {
			e.Log.Errorf("events.send: wrong data type: %t", data)
			return errors.New("events.send: wrong data type")
		}
		e.Log.Infof("sending episode %s", m.ID)
		e.TowerEpisodes <- m
	case "tower.series":
		m, ok := data.(*EventTowerSeries)
		if !ok {
			e.Log.Errorf("events.send: wrong data type: %t", data)
			return errors.New("events.send: wrong data type")
		}
		e.TowerSeries <- m
	case "tower.movies":
		m, ok := data.(*EventTowerMovie)
		if !ok {
			e.Log.Errorf("events.send: wrong data type: %t", data)
			return errors.New("events.send: wrong data type")
		}
		e.TowerMovies <- m
	case "tower.requests":
		m, ok := data.(*EventTowerRequest)
		if !ok {
			e.Log.Errorf("events.send: wrong data type: %t", data)
			return errors.New("events.send: wrong data type")
		}
		e.TowerEvents <- m
	case "tower.downloads":
		m, ok := data.(*EventTowerDownload)
		if !ok {
			e.Log.Errorf("events.send: wrong data type: %t", data)
			return errors.New("events.send: wrong data type")
		}
		e.TowerDownloads <- m
	default:
		e.Log.Warnf("events.send: unknown topic: %s", topic)
	}
	return nil
}

func setupEvents() (err error) {
	events, err = NewEvents()
	if err != nil {
		return err
	}
	return nil
}
