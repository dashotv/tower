// This file is autogenerated by Golem
// Do NOT make modifications, they will be lost
package app

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"github.com/dashotv/grimoire"
)

type Connector struct {
	log       *zap.SugaredLogger
	Download  *grimoire.Store[*Download]
	Episode   *grimoire.Store[*Episode]
	Feed      *grimoire.Store[*Feed]
	Medium    *grimoire.Store[*Medium]
	Message   *grimoire.Store[*Message]
	MinionJob *grimoire.Store[*MinionJob]
	Movie     *grimoire.Store[*Movie]
	Pin       *grimoire.Store[*Pin]
	Release   *grimoire.Store[*Release]
	Request   *grimoire.Store[*Request]
	Series    *grimoire.Store[*Series]
	User      *grimoire.Store[*User]
	Watch     *grimoire.Store[*Watch]
}

func NewConnector() (*Connector, error) {
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

	s, err = settingsFor("message")
	if err != nil {
		return nil, err
	}
	message, err := grimoire.New[*Message](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}

	s, err = settingsFor("minion_job")
	if err != nil {
		return nil, err
	}
	minion_job, err := grimoire.New[*MinionJob](s.URI, s.Database, s.Collection)
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

	s, err = settingsFor("pin")
	if err != nil {
		return nil, err
	}
	pin, err := grimoire.New[*Pin](s.URI, s.Database, s.Collection)
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

	s, err = settingsFor("request")
	if err != nil {
		return nil, err
	}
	request, err := grimoire.New[*Request](s.URI, s.Database, s.Collection)
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

	s, err = settingsFor("user")
	if err != nil {
		return nil, err
	}
	user, err := grimoire.New[*User](s.URI, s.Database, s.Collection)
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
		log:       log.Named("db"),
		Download:  download,
		Episode:   episode,
		Feed:      feed,
		Medium:    medium,
		Message:   message,
		MinionJob: minion_job,
		Movie:     movie,
		Pin:       pin,
		Release:   release,
		Request:   request,
		Series:    series,
		User:      user,
		Watch:     watch,
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
	MediumId  primitive.ObjectID `json:"medium_id" bson:"medium_id"`
	Medium    *Medium            `json:"medium" bson:"-"`
	Auto      bool               `json:"auto" bson:"auto"`
	Multi     bool               `json:"multi" bson:"multi"`
	Force     bool               `json:"force" bson:"force"`
	Url       string             `json:"url" bson:"url"`
	ReleaseId string             `json:"release_id" bson:"tdo_id"`
	Thash     string             `json:"thash" bson:"thash"`
	Selected  string             `json:"selected" bson:"selected"`
	Status    string             `json:"status" bson:"status"`
	Files     []*DownloadFile    `json:"download_files" bson:"download_files"`
}

type DownloadFile struct { // struct
	Id       primitive.ObjectID `json:"id" bson:"_id"`
	MediumId primitive.ObjectID `json:"medium_id" bson:"medium_id"`
	Medium   *Medium            `json:"medium" bson:"medium"`
	Num      int                `json:"num" bson:"num"`
}

type Episode struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	//CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	//UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	Type           string             `json:"type" bson:"_type"`
	SeriesId       primitive.ObjectID `json:"series_id" bson:"series_id"`
	SourceId       string             `json:"source_id" bson:"source_id"`
	Title          string             `json:"title" bson:"title"`
	Description    string             `json:"description" bson:"description"`
	Directory      string             `json:"directory" bson:"directory"`
	Search         string             `json:"search" bson:"search"`
	SeasonNumber   int                `json:"season_number" bson:"season_number"`
	EpisodeNumber  int                `json:"episode_number" bson:"episode_number"`
	AbsoluteNumber int                `json:"absolute_number" bson:"absolute_number"`
	Downloaded     bool               `json:"downloaded" bson:"downloaded"`
	Completed      bool               `json:"completed" bson:"completed"`
	Skipped        bool               `json:"skipped" bson:"skipped"`
	ReleaseDate    time.Time          `json:"release_date" bson:"release_date"`
	Paths          []*Path            `json:"paths" bson:"paths,omitempty"`
	Cover          string             `json:"cover" bson:"-"`
	Background     string             `json:"background" bson:"-"`
	Active         bool               `json:"active" bson:"-"`
	Unwatched      int                `json:"unwatched" bson:"-"`
	Display        string             `json:"display" bson:"-"`
}

type Feed struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	//CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	//UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	Name      string    `json:"name" bson:"name"`
	Url       string    `json:"url" bson:"url"`
	Source    string    `json:"source" bson:"source"`
	Type      string    `json:"type" bson:"type"`
	Active    bool      `json:"active" bson:"active"`
	Processed time.Time `json:"processed" bson:"processed"`
}

type Medium struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	//CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	//UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	Type           string             `json:"type" bson:"_type"`
	Kind           primitive.Symbol   `json:"kind" bson:"kind"`
	Source         string             `json:"source" bson:"source"`
	SourceId       string             `json:"source_id" bson:"source_id"`
	Title          string             `json:"title" bson:"title"`
	Description    string             `json:"description" bson:"description"`
	Slug           string             `json:"slug" bson:"slug"`
	Text           []string           `json:"text" bson:"text"`
	Display        string             `json:"display" bson:"display"`
	Directory      string             `json:"directory" bson:"directory"`
	Search         string             `json:"search" bson:"search"`
	SearchParams   *SearchParams      `json:"search_params" bson:"search_params"`
	Active         bool               `json:"active" bson:"active"`
	Downloaded     bool               `json:"downloaded" bson:"downloaded"`
	Completed      bool               `json:"completed" bson:"completed"`
	Skipped        bool               `json:"skipped" bson:"skipped"`
	Watched        bool               `json:"watched" bson:"watched"`
	Broken         bool               `json:"broken" bson:"broken"`
	Favorite       bool               `json:"favorite" bson:"favorite"`
	Unwatched      int                `json:"unwatched" bson:"unwatched"`
	ReleaseDate    time.Time          `json:"release_date" bson:"release_date"`
	Paths          []*Path            `json:"paths" bson:"paths"`
	Cover          string             `json:"cover" bson:"-"`
	Background     string             `json:"background" bson:"-"`
	SeriesId       primitive.ObjectID `json:"series_id" bson:"series_id"`
	SeasonNumber   int                `json:"season_number" bson:"season_number"`
	EpisodeNumber  int                `json:"episode_number" bson:"episode_number"`
	AbsoluteNumber int                `json:"absolute_number" bson:"absolute_number"`
}

type Message struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	//CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	//UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	Level    string `json:"level" bson:"level"`
	Facility string `json:"facility" bson:"facility"`
	Message  string `json:"message" bson:"message"`
}

type MinionJob struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	//CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	//UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	Name        string    `json:"name" bson:"name"`
	ProcessedAt time.Time `json:"processed_at" bson:"processed_at"`
	Duration    float64   `json:"duration" bson:"duration"`
	Error       string    `json:"error" bson:"error"`
}

type Movie struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	//CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	//UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	Type         string           `json:"type" bson:"_type"`
	Kind         primitive.Symbol `json:"kind" bson:"kind"`
	Source       string           `json:"source" bson:"source"`
	SourceId     string           `json:"source_id" bson:"source_id"`
	ImdbId       string           `json:"imdb_id" bson:"imdb_id"`
	Title        string           `json:"title" bson:"title"`
	Description  string           `json:"description" bson:"description"`
	Slug         string           `json:"slug" bson:"slug"`
	Text         []string         `json:"text" bson:"text"`
	Display      string           `json:"display" bson:"display"`
	Directory    string           `json:"directory" bson:"directory"`
	Search       string           `json:"search" bson:"search"`
	SearchParams *SearchParams    `json:"search_params" bson:"search_params"`
	Active       bool             `json:"active" bson:"active"`
	Downloaded   bool             `json:"downloaded" bson:"downloaded"`
	Completed    bool             `json:"completed" bson:"completed"`
	Skipped      bool             `json:"skipped" bson:"skipped"`
	Watched      bool             `json:"watched" bson:"watched"`
	Broken       bool             `json:"broken" bson:"broken"`
	Favorite     bool             `json:"favorite" bson:"favorite"`
	ReleaseDate  time.Time        `json:"release_date" bson:"release_date"`
	Paths        []*Path          `json:"paths" bson:"paths"`
	Cover        string           `json:"cover" bson:"-"`
	Background   string           `json:"background" bson:"-"`
}

type Path struct { // struct
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Type      primitive.Symbol   `json:"type" bson:"type"`
	Remote    string             `json:"remote" bson:"remote"`
	Local     string             `json:"local" bson:"local"`
	Size      int                `json:"size" bson:"size"`
	Extension string             `json:"extension" bson:"extension"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type Pin struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	//CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	//UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	Pin        int    `json:"id" bson:"pin"`
	Code       string `json:"code" bson:"code"`
	Token      string `json:"authToken" bson:"token"`
	Product    string `json:"product" bson:"product"`
	Identifier string `json:"clientIdentifier" bson:"identifier"`
}

type Release struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	//CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	//UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	Type        string    `json:"type" bson:"type"`
	Source      string    `json:"source" bson:"source"`
	Raw         string    `json:"raw" bson:"raw"`
	Title       string    `json:"title" bson:"title"`
	Description string    `json:"description" bson:"description"`
	Size        string    `json:"size" bson:"size"`
	View        string    `json:"view" bson:"view"`
	Download    string    `json:"download" bson:"download"`
	Infohash    string    `json:"infohash" bson:"infohash"`
	Name        string    `json:"name" bson:"name"`
	Season      int       `json:"season" bson:"season"`
	Episode     int       `json:"episode" bson:"episode"`
	Volume      int       `json:"volume" bson:"volume"`
	Checksum    string    `json:"checksum" bson:"checksum"`
	Group       string    `json:"group" bson:"group"`
	Author      string    `json:"author" bson:"author"`
	Verified    bool      `json:"verified" bson:"verified"`
	Widescreen  bool      `json:"widescreen" bson:"widescreen"`
	Uncensored  bool      `json:"uncensored" bson:"uncensored"`
	Bluray      bool      `json:"bluray" bson:"bluray"`
	Resolution  string    `json:"resolution" bson:"resolution"`
	Encoding    string    `json:"encoding" bson:"encoding"`
	Quality     string    `json:"quality" bson:"quality"`
	PublishedAt time.Time `json:"published_at" bson:"published_at"`
}

type Request struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	//CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	//UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	Title    string `json:"title" bson:"title"`
	User     string `json:"user" bson:"user"`
	Type     string `json:"type" bson:"type"`
	Source   string `json:"source" bson:"source"`
	SourceId string `json:"source_id" bson:"source_id"`
	Status   string `json:"status" bson:"status"`
}

type SearchParams struct { // struct
	Type       string `json:"type" bson:"type"`
	Verified   bool   `json:"verified" bson:"verified"`
	Group      string `json:"group" bson:"group"`
	Author     string `json:"author" bson:"author"`
	Resolution int    `json:"resolution" bson:"resolution"`
	Source     string `json:"source" bson:"source"`
	Uncensored bool   `json:"uncensored" bson:"uncensored"`
	Bluray     bool   `json:"bluray" bson:"bluray"`
}

type Series struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	//CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	//UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	Type          string           `json:"type" bson:"_type"`
	Kind          primitive.Symbol `json:"kind" bson:"kind"`
	Source        string           `json:"source" bson:"source"`
	SourceId      string           `json:"source_id" bson:"source_id"`
	ImdbId        string           `json:"imdb_id" bson:"imdb_id"`
	Title         string           `json:"title" bson:"title"`
	Description   string           `json:"description" bson:"description"`
	Slug          string           `json:"slug" bson:"slug"`
	Text          []string         `json:"text" bson:"text"`
	Display       string           `json:"display" bson:"display"`
	Directory     string           `json:"directory" bson:"directory"`
	Search        string           `json:"search" bson:"search"`
	SearchParams  *SearchParams    `json:"search_params" bson:"search_params"`
	Status        string           `json:"status" bson:"status"`
	Active        bool             `json:"active" bson:"active"`
	Downloaded    bool             `json:"downloaded" bson:"downloaded"`
	Completed     bool             `json:"completed" bson:"completed"`
	Skipped       bool             `json:"skipped" bson:"skipped"`
	Watched       bool             `json:"watched" bson:"watched"`
	Broken        bool             `json:"broken" bson:"broken"`
	Favorite      bool             `json:"favorite" bson:"favorite"`
	Unwatched     int              `json:"unwatched" bson:"-"`
	ReleaseDate   time.Time        `json:"release_date" bson:"release_date"`
	Paths         []*Path          `json:"paths" bson:"paths"`
	Cover         string           `json:"cover" bson:"-"`
	Background    string           `json:"background" bson:"-"`
	CurrentSeason int              `json:"currentSeason" bson:"-"`
	Seasons       []int            `json:"seasons" bson:"-"`
	Episodes      []*Episode       `json:"episodes" bson:"-"`
	Watches       []*Watch         `json:"watches" bson:"-"`
}

type User struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	//CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	//UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	Name  string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
	Token string `json:"token" bson:"token"`
	Thumb string `json:"thumb" bson:"thumb"`
	Home  bool   `json:"home" bson:"home"`
	Admin bool   `json:"admin" bson:"admin"`
}

type Watch struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	//CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	//UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	Username  string             `json:"username" bson:"username"`
	Player    string             `json:"player" bson:"player"`
	WatchedAt time.Time          `json:"watched_at" bson:"watched_at"`
	MediumId  primitive.ObjectID `json:"medium_id" bson:"medium_id"`
	Medium    *Medium            `json:"medium" bson:"-"`
}
