package app

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/dashotv/grimoire"
)

var cfg *Config

type Connector struct {
	Download *grimoire.Store[*Download]
	Episode  *grimoire.Store[*Episode]
	Feed     *grimoire.Store[*Feed]
	Medium   *grimoire.Store[*Medium]
	Movie    *grimoire.Store[*Movie]
	Release  *grimoire.Store[*Release]
	Series   *grimoire.Store[*Series]
	Watch    *grimoire.Store[*Watch]
}

func NewConnector() (*Connector, error) {
	cfg = ConfigInstance()
	var s *Connection
	var err error

	s, err = settingsFor("download")
	if err != nil {
		return nil, err
	}

	download, err := grimoire.New[*Download](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}
	s, err = settingsFor("episode")
	if err != nil {
		return nil, err
	}

	episode, err := grimoire.New[*Episode](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}
	s, err = settingsFor("feed")
	if err != nil {
		return nil, err
	}

	feed, err := grimoire.New[*Feed](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}
	s, err = settingsFor("medium")
	if err != nil {
		return nil, err
	}

	medium, err := grimoire.New[*Medium](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}
	s, err = settingsFor("movie")
	if err != nil {
		return nil, err
	}

	movie, err := grimoire.New[*Movie](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}
	s, err = settingsFor("release")
	if err != nil {
		return nil, err
	}

	release, err := grimoire.New[*Release](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}
	s, err = settingsFor("series")
	if err != nil {
		return nil, err
	}

	series, err := grimoire.New[*Series](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}
	s, err = settingsFor("watch")
	if err != nil {
		return nil, err
	}

	watch, err := grimoire.New[*Watch](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}

	c := &Connector{
		Download: download,
		Episode:  episode,
		Feed:     feed,
		Medium:   medium,
		Movie:    movie,
		Release:  release,
		Series:   series,
		Watch:    watch,
	}

	return c, nil
}

func settingsFor(name string) (*Connection, error) {
	if cfg.Connections["default"] == nil {
		return nil, fmt.Errorf("no default config while configuring %s", name)
	}

	if _, ok := cfg.Connections[name]; !ok {
		return cfg.Connections["default"], nil
	}

	s := cfg.Connections["default"]
	a := cfg.Connections[name]

	if a.URI != "" {
		s.URI = a.URI
	}
	if a.Database != "" {
		s.Database = a.Database
	}
	if a.Collection != "" {
		s.Collection = a.Collection
	}

	return s, nil
}

type Download struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	//CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	//UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	MediumId  primitive.ObjectID `json:&#34;medium_id&#34; bson:&#34;medium_id&#34;`
	Medium    *Medium            `json:&#34;medium&#34; bson:&#34;-&#34;`
	Auto      bool               `json:&#34;auto&#34; bson:&#34;auto&#34;`
	Multi     bool               `json:&#34;multi&#34; bson:&#34;multi&#34;`
	Force     bool               `json:&#34;force&#34; bson:&#34;force&#34;`
	Url       string             `json:&#34;url&#34; bson:&#34;url&#34;`
	ReleaseId string             `json:&#34;release_id&#34; bson:&#34;tdo_id&#34;`
	Thash     string             `json:&#34;thash&#34; bson:&#34;thash&#34;`
	Selected  string             `json:&#34;selected&#34; bson:&#34;selected&#34;`
	Status    string             `json:&#34;status&#34; bson:&#34;status&#34;`
	Files     []*DownloadFile    `json:&#34;download_files&#34; bson:&#34;download_files&#34;`
}

type DownloadFile struct { // struct
	Id       primitive.ObjectID `json:&#34;id&#34; bson:&#34;_id&#34;`
	MediumId primitive.ObjectID `json:&#34;medium_id&#34; bson:&#34;medium_id&#34;`
	Medium   *Medium            `json:&#34;medium&#34; bson:&#34;medium&#34;`
	Num      int                `json:&#34;num&#34; bson:&#34;num&#34;`
}

type Episode struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	//CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	//UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	Type           string             `json:&#34;type&#34; bson:&#34;_type&#34;`
	SeriesId       primitive.ObjectID `json:&#34;series_id&#34; bson:&#34;series_id&#34;`
	Kind           primitive.Symbol   `json:&#34;kind&#34; bson:&#34;kind&#34;`
	Source         string             `json:&#34;source&#34; bson:&#34;source&#34;`
	SourceId       string             `json:&#34;source_id&#34; bson:&#34;source_id&#34;`
	Title          string             `json:&#34;title&#34; bson:&#34;title&#34;`
	Description    string             `json:&#34;description&#34; bson:&#34;description&#34;`
	Slug           string             `json:&#34;slug&#34; bson:&#34;slug&#34;`
	Text           []string           `json:&#34;text&#34; bson:&#34;text&#34;`
	Display        string             `json:&#34;display&#34; bson:&#34;display&#34;`
	Directory      string             `json:&#34;directory&#34; bson:&#34;directory&#34;`
	Search         string             `json:&#34;search&#34; bson:&#34;search&#34;`
	SeasonNumber   int                `json:&#34;season_number&#34; bson:&#34;season_number&#34;`
	EpisodeNumber  int                `json:&#34;episode_number&#34; bson:&#34;episode_number&#34;`
	AbsoluteNumber int                `json:&#34;absolute_number&#34; bson:&#34;absolute_number&#34;`
	SearchParams   *SearchParams      `json:&#34;search_params&#34; bson:&#34;search_params&#34;`
	Active         bool               `json:&#34;active&#34; bson:&#34;active&#34;`
	Downloaded     bool               `json:&#34;downloaded&#34; bson:&#34;downloaded&#34;`
	Completed      bool               `json:&#34;completed&#34; bson:&#34;completed&#34;`
	Skipped        bool               `json:&#34;skipped&#34; bson:&#34;skipped&#34;`
	Watched        bool               `json:&#34;watched&#34; bson:&#34;watched&#34;`
	Broken         bool               `json:&#34;broken&#34; bson:&#34;broken&#34;`
	Unwatched      int                `json:&#34;unwatched&#34; bson:&#34;-&#34;`
	ReleaseDate    time.Time          `json:&#34;release_date&#34; bson:&#34;release_date&#34;`
	Paths          []Path             `json:&#34;paths&#34; bson:&#34;paths&#34;`
	Cover          string             `json:&#34;cover&#34; bson:&#34;-&#34;`
	Background     string             `json:&#34;background&#34; bson:&#34;-&#34;`
}

type Feed struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	//CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	//UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	Name      string    `json:&#34;name&#34; bson:&#34;name&#34;`
	Url       string    `json:&#34;url&#34; bson:&#34;url&#34;`
	Source    string    `json:&#34;source&#34; bson:&#34;source&#34;`
	Type      string    `json:&#34;type&#34; bson:&#34;type&#34;`
	Active    bool      `json:&#34;active&#34; bson:&#34;active&#34;`
	Processed time.Time `json:&#34;processed&#34; bson:&#34;processed&#34;`
}

type Medium struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	//CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	//UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	Type           string             `json:&#34;type&#34; bson:&#34;_type&#34;`
	Kind           primitive.Symbol   `json:&#34;kind&#34; bson:&#34;kind&#34;`
	Source         string             `json:&#34;source&#34; bson:&#34;source&#34;`
	SourceId       string             `json:&#34;source_id&#34; bson:&#34;source_id&#34;`
	Title          string             `json:&#34;title&#34; bson:&#34;title&#34;`
	Description    string             `json:&#34;description&#34; bson:&#34;description&#34;`
	Slug           string             `json:&#34;slug&#34; bson:&#34;slug&#34;`
	Text           []string           `json:&#34;text&#34; bson:&#34;text&#34;`
	Display        string             `json:&#34;display&#34; bson:&#34;display&#34;`
	Directory      string             `json:&#34;directory&#34; bson:&#34;directory&#34;`
	Search         string             `json:&#34;search&#34; bson:&#34;search&#34;`
	SearchParams   *SearchParams      `json:&#34;search_params&#34; bson:&#34;search_params&#34;`
	Active         bool               `json:&#34;active&#34; bson:&#34;active&#34;`
	Downloaded     bool               `json:&#34;downloaded&#34; bson:&#34;downloaded&#34;`
	Completed      bool               `json:&#34;completed&#34; bson:&#34;completed&#34;`
	Skipped        bool               `json:&#34;skipped&#34; bson:&#34;skipped&#34;`
	Watched        bool               `json:&#34;watched&#34; bson:&#34;watched&#34;`
	Broken         bool               `json:&#34;broken&#34; bson:&#34;broken&#34;`
	Favorite       bool               `json:&#34;favorite&#34; bson:&#34;favorite&#34;`
	Unwatched      int                `json:&#34;unwatched&#34; bson:&#34;unwatched&#34;`
	ReleaseDate    time.Time          `json:&#34;release_date&#34; bson:&#34;release_date&#34;`
	Paths          []Path             `json:&#34;paths&#34; bson:&#34;paths&#34;`
	Cover          string             `json:&#34;cover&#34; bson:&#34;-&#34;`
	Background     string             `json:&#34;background&#34; bson:&#34;-&#34;`
	SeriesId       primitive.ObjectID `json:&#34;series_id&#34; bson:&#34;series_id&#34;`
	SeasonNumber   int                `json:&#34;season_number&#34; bson:&#34;season_number&#34;`
	EpisodeNumber  int                `json:&#34;episode_number&#34; bson:&#34;episode_number&#34;`
	AbsoluteNumber int                `json:&#34;absolute_number&#34; bson:&#34;absolute_number&#34;`
}

type Movie struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	//CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	//UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	Type         string           `json:&#34;type&#34; bson:&#34;_type&#34;`
	Kind         primitive.Symbol `json:&#34;kind&#34; bson:&#34;kind&#34;`
	Source       string           `json:&#34;source&#34; bson:&#34;source&#34;`
	SourceId     string           `json:&#34;source_id&#34; bson:&#34;source_id&#34;`
	Title        string           `json:&#34;title&#34; bson:&#34;title&#34;`
	Description  string           `json:&#34;description&#34; bson:&#34;description&#34;`
	Slug         string           `json:&#34;slug&#34; bson:&#34;slug&#34;`
	Text         []string         `json:&#34;text&#34; bson:&#34;text&#34;`
	Display      string           `json:&#34;display&#34; bson:&#34;display&#34;`
	Directory    string           `json:&#34;directory&#34; bson:&#34;directory&#34;`
	Search       string           `json:&#34;search&#34; bson:&#34;search&#34;`
	SearchParams *SearchParams    `json:&#34;search_params&#34; bson:&#34;search_params&#34;`
	Active       bool             `json:&#34;active&#34; bson:&#34;active&#34;`
	Downloaded   bool             `json:&#34;downloaded&#34; bson:&#34;downloaded&#34;`
	Completed    bool             `json:&#34;completed&#34; bson:&#34;completed&#34;`
	Skipped      bool             `json:&#34;skipped&#34; bson:&#34;skipped&#34;`
	Watched      bool             `json:&#34;watched&#34; bson:&#34;watched&#34;`
	Broken       bool             `json:&#34;broken&#34; bson:&#34;broken&#34;`
	Favorite     bool             `json:&#34;favorite&#34; bson:&#34;favorite&#34;`
	ReleaseDate  time.Time        `json:&#34;release_date&#34; bson:&#34;release_date&#34;`
	Paths        []Path           `json:&#34;paths&#34; bson:&#34;paths&#34;`
	Cover        string           `json:&#34;cover&#34; bson:&#34;-&#34;`
	Background   string           `json:&#34;background&#34; bson:&#34;-&#34;`
}

type Release struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	//CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	//UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	Type        string    `json:&#34;type&#34; bson:&#34;type&#34;`
	Source      string    `json:&#34;source&#34; bson:&#34;source&#34;`
	Raw         string    `json:&#34;raw&#34; bson:&#34;raw&#34;`
	Title       string    `json:&#34;title&#34; bson:&#34;title&#34;`
	Description string    `json:&#34;description&#34; bson:&#34;description&#34;`
	Size        string    `json:&#34;size&#34; bson:&#34;size&#34;`
	View        string    `json:&#34;view&#34; bson:&#34;view&#34;`
	Download    string    `json:&#34;download&#34; bson:&#34;download&#34;`
	Infohash    string    `json:&#34;infohash&#34; bson:&#34;infohash&#34;`
	Name        string    `json:&#34;name&#34; bson:&#34;name&#34;`
	Season      int       `json:&#34;season&#34; bson:&#34;season&#34;`
	Episode     int       `json:&#34;episode&#34; bson:&#34;episode&#34;`
	Volume      int       `json:&#34;volume&#34; bson:&#34;volume&#34;`
	Checksum    string    `json:&#34;checksum&#34; bson:&#34;checksum&#34;`
	Group       string    `json:&#34;group&#34; bson:&#34;group&#34;`
	Author      string    `json:&#34;author&#34; bson:&#34;author&#34;`
	Verified    bool      `json:&#34;verified&#34; bson:&#34;verified&#34;`
	Widescreen  bool      `json:&#34;widescreen&#34; bson:&#34;widescreen&#34;`
	Uncensored  bool      `json:&#34;uncensored&#34; bson:&#34;uncensored&#34;`
	Bluray      bool      `json:&#34;bluray&#34; bson:&#34;bluray&#34;`
	Resolution  string    `json:&#34;resolution&#34; bson:&#34;resolution&#34;`
	Encoding    string    `json:&#34;encoding&#34; bson:&#34;encoding&#34;`
	Quality     string    `json:&#34;quality&#34; bson:&#34;quality&#34;`
	PublishedAt time.Time `json:&#34;published_at&#34; bson:&#34;published_at&#34;`
}

type SearchParams struct { // struct
	Type       string `json:&#34;type&#34; bson:&#34;type&#34;`
	Verified   bool   `json:&#34;verified&#34; bson:&#34;verified&#34;`
	Group      string `json:&#34;group&#34; bson:&#34;group&#34;`
	Author     string `json:&#34;author&#34; bson:&#34;author&#34;`
	Resolution int    `json:&#34;resolution&#34; bson:&#34;resolution&#34;`
	Source     string `json:&#34;source&#34; bson:&#34;source&#34;`
	Uncensored bool   `json:&#34;uncensored&#34; bson:&#34;uncensored&#34;`
	Bluray     bool   `json:&#34;bluray&#34; bson:&#34;bluray&#34;`
}

type Series struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	//CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	//UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	Type          string           `json:&#34;type&#34; bson:&#34;_type&#34;`
	Kind          primitive.Symbol `json:&#34;kind&#34; bson:&#34;kind&#34;`
	Source        string           `json:&#34;source&#34; bson:&#34;source&#34;`
	SourceId      string           `json:&#34;source_id&#34; bson:&#34;source_id&#34;`
	Title         string           `json:&#34;title&#34; bson:&#34;title&#34;`
	Description   string           `json:&#34;description&#34; bson:&#34;description&#34;`
	Slug          string           `json:&#34;slug&#34; bson:&#34;slug&#34;`
	Text          []string         `json:&#34;text&#34; bson:&#34;text&#34;`
	Display       string           `json:&#34;display&#34; bson:&#34;display&#34;`
	Directory     string           `json:&#34;directory&#34; bson:&#34;directory&#34;`
	Search        string           `json:&#34;search&#34; bson:&#34;search&#34;`
	SearchParams  *SearchParams    `json:&#34;search_params&#34; bson:&#34;search_params&#34;`
	Active        bool             `json:&#34;active&#34; bson:&#34;active&#34;`
	Downloaded    bool             `json:&#34;downloaded&#34; bson:&#34;downloaded&#34;`
	Completed     bool             `json:&#34;completed&#34; bson:&#34;completed&#34;`
	Skipped       bool             `json:&#34;skipped&#34; bson:&#34;skipped&#34;`
	Watched       bool             `json:&#34;watched&#34; bson:&#34;watched&#34;`
	Broken        bool             `json:&#34;broken&#34; bson:&#34;broken&#34;`
	Favorite      bool             `json:&#34;favorite&#34; bson:&#34;favorite&#34;`
	Unwatched     int              `json:&#34;unwatched&#34; bson:&#34;-&#34;`
	ReleaseDate   time.Time        `json:&#34;release_date&#34; bson:&#34;release_date&#34;`
	Paths         []Path           `json:&#34;paths&#34; bson:&#34;paths&#34;`
	Cover         string           `json:&#34;cover&#34; bson:&#34;-&#34;`
	Background    string           `json:&#34;background&#34; bson:&#34;-&#34;`
	CurrentSeason int              `json:&#34;currentSeason&#34; bson:&#34;-&#34;`
	Seasons       []int            `json:&#34;seasons&#34; bson:&#34;-&#34;`
	Episodes      []*Episode       `json:&#34;episodes&#34; bson:&#34;-&#34;`
	Watches       []*Watch         `json:&#34;watches&#34; bson:&#34;-&#34;`
}

type Watch struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	//CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	//UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	Username  string             `json:&#34;username&#34; bson:&#34;username&#34;`
	Player    string             `json:&#34;player&#34; bson:&#34;player&#34;`
	WatchedAt time.Time          `json:&#34;watched_at&#34; bson:&#34;watched_at&#34;`
	MediumId  primitive.ObjectID `json:&#34;medium_id&#34; bson:&#34;medium_id&#34;`
	Medium    *Medium            `json:&#34;medium&#34; bson:&#34;-&#34;`
}
